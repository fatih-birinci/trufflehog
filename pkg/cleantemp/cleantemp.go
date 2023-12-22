package cleantemp

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/go-ps"

	logContext "github.com/trufflesecurity/trufflehog/v3/pkg/context"
)

const (
	defaultExecPath             = "trufflehog"
	defaultArtifactPrefixFormat = "%s-%d-"
)

// MkdirTemp returns a temporary directory path formatted as:
// trufflehog-<pid>-<randint>
func MkdirTemp() (string, error) {
	pid := os.Getpid()
	tmpdir := fmt.Sprintf(defaultArtifactPrefixFormat, defaultExecPath, pid)
	dir, err := os.MkdirTemp(os.TempDir(), tmpdir)
	if err != nil {
		return "", err
	}
	return dir, nil
}

// Unlike MkdirTemp, we only want to generate the filename string.
// The tempfile creation in trufflehog we're interested in
// is generally handled by "github.com/trufflesecurity/disk-buffer-reader"
func MkFilename() string {
	pid := os.Getpid()
	filename := fmt.Sprintf(defaultArtifactPrefixFormat, defaultExecPath, pid)
	return filename
}

// Only compile during startup.
var trufflehogRE = regexp.MustCompile(`^trufflehog-\d+-\d+$`)

// CleanTempArtifacts deletes orphaned temp directories and files that do not contain running PID values.
func CleanTempArtifacts(ctx logContext.Context) error {
	executablePath, err := os.Executable()
	if err != nil {
		executablePath = defaultExecPath
	}
	execName := filepath.Base(executablePath)

	var pids []string
	procs, err := ps.Processes()
	if err != nil {
		return fmt.Errorf("error getting jobs PIDs: %w", err)
	}

	for _, proc := range procs {
		if proc.Executable() == execName {
			pids = append(pids, strconv.Itoa(proc.Pid()))
		}
	}

	if len(pids) == 0 {
		ctx.Logger().V(5).Info("No trufflehog processes were found")
		return nil
	}

	tempDir := os.TempDir()
	err = filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking temp dir: %w", err)
		}
		if trufflehogRE.MatchString(d.Name()) {
			// Mark these artifacts initially as ones that should be deleted.
			shouldDelete := true
			// Check if the name matches any live PIDs.
			// Potential race condition here if a PID is started and creates tmp data after the initial check.
			for _, pidval := range pids {
				if strings.Contains(d.Name(), fmt.Sprintf("-%s-", pidval)) {
					shouldDelete = false
					break
				}
			}

			if shouldDelete {
				var err error
				if d.IsDir() {
					err = os.RemoveAll(path)
				} else {
					err = os.Remove(path)
				}
				if err != nil {
					return fmt.Errorf("error deleting temp artifact: %s", path)
				}

				ctx.Logger().V(4).Info("Deleted orphaned temp artifact", "artifact", path)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking temp dir: %w", err)
	}

	return nil
}

// RunCleanupLoop runs a loop that cleans up orphaned directories every 15 seconds.
func RunCleanupLoop(ctx logContext.Context) {
	err := CleanTempArtifacts(ctx)
	if err != nil {
		ctx.Logger().Error(err, "Error cleaning up orphaned directories ")
	}

	const cleanupLoopInterval = 15 * time.Second
	ticker := time.NewTicker(cleanupLoopInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := CleanTempArtifacts(ctx); err != nil {
				ctx.Logger().Error(err, "error cleaning up orphaned directories")
			}
		case <-ctx.Done():
			ctx.Logger().Info("Cleanup loop exiting due to context cancellation")
			return
		}
	}
}
