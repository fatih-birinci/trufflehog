package github

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/go-errors/errors"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-logr/logr"
	"github.com/gobwas/glob"
	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/trufflesecurity/trufflehog/v3/pkg/cache"
	"github.com/trufflesecurity/trufflehog/v3/pkg/cache/memory"
	"github.com/trufflesecurity/trufflehog/v3/pkg/common"
	"github.com/trufflesecurity/trufflehog/v3/pkg/context"
	"github.com/trufflesecurity/trufflehog/v3/pkg/pb/credentialspb"
	"github.com/trufflesecurity/trufflehog/v3/pkg/pb/source_metadatapb"
	"github.com/trufflesecurity/trufflehog/v3/pkg/pb/sourcespb"
	"github.com/trufflesecurity/trufflehog/v3/pkg/sanitizer"
	"github.com/trufflesecurity/trufflehog/v3/pkg/sources"
	"github.com/trufflesecurity/trufflehog/v3/pkg/sources/git"
)

const (
	unauthGithubOrgRateLimt = 30
	defaultPagination       = 100
	membersAppPagination    = 500
)

type Source struct {
	name string
	// Protects the user and token.
	userMu      sync.Mutex
	githubUser  string
	githubToken string

	sourceID          int64
	jobID             int64
	verify            bool
	repos             []string
	members           []string
	orgsCache         cache.Cache
	filteredRepoCache *filteredRepoCache
	memberCache       map[string]struct{}
	repoSizes         repoSize
	totalRepoSize     int // total size in bytes of all repos
	git               *git.Git

	scanOptMu   sync.Mutex // protects the scanOptions
	scanOptions *git.ScanOptions

	httpClient      *http.Client
	log             logr.Logger
	conn            *sourcespb.GitHub
	jobPool         *errgroup.Group
	resumeInfoMutex sync.Mutex
	resumeInfoSlice []string
	apiClient       *github.Client

	mu        sync.Mutex // protects the visibility maps
	publicMap map[string]source_metadatapb.Visibility

	includePRComments    bool
	includeIssueComments bool
	sources.Progress
	sources.CommonSourceUnitUnmarshaller
}

func (s *Source) WithScanOptions(scanOptions *git.ScanOptions) {
	s.scanOptions = scanOptions
}

func (s *Source) setScanOptions(base, head string) {
	s.scanOptMu.Lock()
	defer s.scanOptMu.Unlock()
	s.scanOptions.BaseHash = base
	s.scanOptions.BaseHash = head
}

// Ensure the Source satisfies the interfaces at compile time
var _ sources.Source = (*Source)(nil)
var _ sources.SourceUnitUnmarshaller = (*Source)(nil)

var endsWithGithub = regexp.MustCompile(`github\.com/?$`)

// Type returns the type of source.
// It is used for matching source types in configuration and job input.
func (s *Source) Type() sourcespb.SourceType {
	return sourcespb.SourceType_SOURCE_TYPE_GITHUB
}

func (s *Source) SourceID() int64 {
	return s.sourceID
}

func (s *Source) JobID() int64 {
	return s.jobID
}

type repoSize struct {
	mu        sync.RWMutex
	repoSizes map[string]int // size in bytes of each repo
}

func (r *repoSize) addRepo(repo string, size int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.repoSizes[repo] = size
}

func (r *repoSize) getRepo(repo string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.repoSizes[repo]
}

func newRepoSize() repoSize {
	return repoSize{repoSizes: make(map[string]int)}
}

// filteredRepoCache is a wrapper around cache.Cache that filters out repos
// based on include and exclude globs.
type filteredRepoCache struct {
	cache.Cache
	include, exclude []glob.Glob
}

func (s *Source) newFilteredRepoCache(c cache.Cache, include, exclude []string) *filteredRepoCache {
	includeGlobs := make([]glob.Glob, 0, len(include))
	excludeGlobs := make([]glob.Glob, 0, len(exclude))
	for _, ig := range include {
		g, err := glob.Compile(ig)
		if err != nil {
			s.log.V(1).Info("invalid include glob", "glob", g, "err", err)
		}
		includeGlobs = append(includeGlobs, g)
	}
	for _, eg := range exclude {
		g, err := glob.Compile(eg)
		if err != nil {
			s.log.V(1).Info("invalid exclude glob", "glob", g, "err", err)
		}
		excludeGlobs = append(excludeGlobs, g)
	}
	return &filteredRepoCache{Cache: c, include: includeGlobs, exclude: excludeGlobs}
}

// Set overrides the cache.Cache Set method to filter out repos based on
// include and exclude globs.
func (c *filteredRepoCache) Set(key, val string) {
	if c.ignoreRepo(key) {
		return
	}
	if !c.includeRepo(key) {
		return
	}
	c.Cache.Set(key, val)
}

func (c *filteredRepoCache) ignoreRepo(s string) bool {
	for _, g := range c.exclude {
		if g.Match(s) {
			return true
		}
	}
	return false
}

func (c *filteredRepoCache) includeRepo(s string) bool {
	if len(c.include) == 0 {
		return true
	}

	for _, g := range c.include {
		if g.Match(s) {
			return true
		}
	}
	return false
}

// Init returns an initialized GitHub source.
func (s *Source) Init(aCtx context.Context, name string, jobID, sourceID int64, verify bool, connection *anypb.Any, concurrency int) error {
	s.log = aCtx.Logger()

	s.name = name
	s.sourceID = sourceID
	s.jobID = jobID
	s.verify = verify
	s.jobPool = &errgroup.Group{}
	s.jobPool.SetLimit(concurrency)

	s.httpClient = common.RetryableHttpClientTimeout(60)
	s.apiClient = github.NewClient(s.httpClient)

	var conn sourcespb.GitHub
	err := anypb.UnmarshalTo(connection, &conn, proto.UnmarshalOptions{})
	if err != nil {
		return errors.WrapPrefix(err, "error unmarshalling connection", 0)
	}
	s.conn = &conn

	s.filteredRepoCache = s.newFilteredRepoCache(memory.New(), s.conn.IncludeRepos, s.conn.IgnoreRepos)
	s.memberCache = make(map[string]struct{})

	s.repoSizes = newRepoSize()
	s.repos = s.conn.Repositories
	for _, repo := range s.repos {
		r, err := s.normalizeRepo(repo)
		if err != nil {
			aCtx.Logger().Error(err, "invalid repository", "repo", repo)
			continue
		}
		s.filteredRepoCache.Set(r, r)
	}

	s.includeIssueComments = s.conn.IncludeIssueComments
	s.includePRComments = s.conn.IncludePullRequestComments

	s.orgsCache = memory.New()
	for _, org := range s.conn.Organizations {
		s.orgsCache.Set(org, org)
	}

	// Head or base should only be used with incoming webhooks
	if (len(s.conn.Head) > 0 || len(s.conn.Base) > 0) && len(s.repos) != 1 {
		return fmt.Errorf("cannot specify head or base with multiple repositories")
	}

	err = git.GitCmdCheck()
	if err != nil {
		return err
	}

	s.publicMap = map[string]source_metadatapb.Visibility{}

	s.git = git.NewGit(s.Type(), s.JobID(), s.SourceID(), s.name, s.verify, runtime.NumCPU(),
		func(file, email, commit, timestamp, repository string, line int64) *source_metadatapb.MetaData {
			return &source_metadatapb.MetaData{
				Data: &source_metadatapb.MetaData_Github{
					Github: &source_metadatapb.Github{
						Commit:     sanitizer.UTF8(commit),
						File:       sanitizer.UTF8(file),
						Email:      sanitizer.UTF8(email),
						Repository: sanitizer.UTF8(repository),
						Link:       git.GenerateLink(repository, commit, file, line),
						Timestamp:  sanitizer.UTF8(timestamp),
						Line:       line,
						Visibility: s.visibilityOf(aCtx, repository),
					},
				},
			}
		})

	return nil
}

func (s *Source) visibilityOf(ctx context.Context, repoURL string) (visibility source_metadatapb.Visibility) {
	s.mu.Lock()
	visibility, ok := s.publicMap[repoURL]
	s.mu.Unlock()
	if ok {
		return visibility
	}

	visibility = source_metadatapb.Visibility_public
	defer func() {
		s.mu.Lock()
		s.publicMap[repoURL] = visibility
		s.mu.Unlock()
	}()
	logger := s.log.WithValues("repo", repoURL)
	logger.V(2).Info("Checking public status")
	u, err := url.Parse(repoURL)
	if err != nil {
		logger.Error(err, "Could not parse repository URL.")
		return
	}

	var resp *github.Response
	urlPathParts := strings.Split(u.Path, "/")
	switch len(urlPathParts) {
	case 2:
		// Check if repoURL is a gist.
		var gist *github.Gist
		repoName := urlPathParts[1]
		repoName = strings.TrimSuffix(repoName, ".git")
		for {
			gist, resp, err = s.apiClient.Gists.Get(ctx, repoName)
			if !s.handleRateLimit(err, resp) {
				break
			}
		}
		if err != nil || gist == nil {
			if _, unauthenticated := s.conn.GetCredential().(*sourcespb.GitHub_Unauthenticated); unauthenticated {
				logger.Info("Unauthenticated scans cannot determine if a repository is private.")
				visibility = source_metadatapb.Visibility_private
			}
			logger.Error(err, "Could not get Github repository")
			return
		}
		if !(*gist.Public) {
			visibility = source_metadatapb.Visibility_private
		}
	case 3:
		var repo *github.Repository
		owner := urlPathParts[1]
		repoName := urlPathParts[2]
		repoName = strings.TrimSuffix(repoName, ".git")
		for {
			repo, resp, err = s.apiClient.Repositories.Get(ctx, owner, repoName)
			if !s.handleRateLimit(err, resp) {
				break
			}
		}
		if err != nil || repo == nil {
			logger.Error(err, "Could not get Github repository")
			if _, unauthenticated := s.conn.GetCredential().(*sourcespb.GitHub_Unauthenticated); unauthenticated {
				logger.Info("Unauthenticated scans cannot determine if a repository is private.")
				visibility = source_metadatapb.Visibility_private
			}
			return
		}
		if *repo.Private {
			visibility = source_metadatapb.Visibility_private
		}
	default:
		logger.Error(fmt.Errorf("unexpected number of parts"), "RepoURL should split into 2 or 3 parts",
			"got", len(urlPathParts),
		)
	}
	return
}

// Chunks emits chunks of bytes over a channel.
func (s *Source) Chunks(ctx context.Context, chunksChan chan *sources.Chunk) error {
	apiEndpoint := s.conn.Endpoint
	if len(apiEndpoint) == 0 || endsWithGithub.MatchString(apiEndpoint) {
		apiEndpoint = "https://api.github.com"
	}

	installationClient, err := s.enumerate(ctx, apiEndpoint)
	if err != nil {
		return err
	}

	return s.scan(ctx, installationClient, chunksChan)
}

func (s *Source) enumerate(ctx context.Context, apiEndpoint string) (*github.Client, error) {
	var (
		installationClient *github.Client
		err                error
	)

	switch cred := s.conn.GetCredential().(type) {
	case *sourcespb.GitHub_BasicAuth:
		if err = s.enumerateBasicAuth(ctx, apiEndpoint, cred.BasicAuth); err != nil {
			return nil, err
		}
	case *sourcespb.GitHub_Unauthenticated:
		s.enumerateUnauthenticated(ctx, apiEndpoint)
	case *sourcespb.GitHub_Token:
		if err = s.enumerateWithToken(ctx, apiEndpoint, cred.Token); err != nil {
			return nil, err
		}
	case *sourcespb.GitHub_GithubApp:
		if installationClient, err = s.enumerateWithApp(ctx, apiEndpoint, cred.GithubApp); err != nil {
			return nil, err
		}
	default:
		// TODO: move this error to Init
		return nil, errors.Errorf("Invalid configuration given for source. Name: %s, Type: %s", s.name, s.Type())
	}

	s.repos = s.filteredRepoCache.Values()
	s.log.Info("Completed enumeration", "num_repos", len(s.repos), "num_orgs", s.orgsCache.Count(), "num_members", len(s.memberCache))

	// We must sort the repos so we can resume later if necessary.
	sort.Strings(s.repos)
	return installationClient, nil
}

func (s *Source) enumerateBasicAuth(ctx context.Context, apiEndpoint string, basicAuth *credentialspb.BasicAuth) error {
	s.httpClient.Transport = &github.BasicAuthTransport{
		Username: basicAuth.Username,
		Password: basicAuth.Password,
	}
	ghClient, err := createGitHubClient(s.httpClient, apiEndpoint)
	if err != nil {
		s.log.Error(err, "error creating GitHub client")
	}
	s.apiClient = ghClient

	for _, org := range s.orgsCache.Keys() {
		if err := s.getReposByOrg(ctx, org); err != nil {
			s.log.Error(err, "error fetching repos for org or user")
		}
	}

	return nil
}

func (s *Source) enumerateUnauthenticated(ctx context.Context, apiEndpoint string) {
	ghClient, err := createGitHubClient(s.httpClient, apiEndpoint)
	if err != nil {
		s.log.Error(err, "error creating GitHub client")
	}
	s.apiClient = ghClient
	if s.orgsCache.Count() > unauthGithubOrgRateLimt {
		s.log.Info("You may experience rate limiting when using the unauthenticated GitHub api. Consider using an authenticated scan instead.")
	}

	for _, org := range s.orgsCache.Keys() {
		if err := s.getReposByOrg(ctx, org); err != nil {
			s.log.Error(err, "error fetching repos for org or user")
		}

		// We probably don't need to do this, since getting repos by org makes more sense?
		if err := s.getReposByUser(ctx, org); err != nil {
			s.log.Error(err, "error fetching repos for org or user")
		}

		if s.conn.ScanUsers {
			s.log.Info("Enumerating unauthenticated does not support scanning organization members")
		}
	}
}

func (s *Source) enumerateWithToken(ctx context.Context, apiEndpoint, token string) error {
	// Needed for clones.
	s.githubToken = token

	// Needed to list repos.
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	s.httpClient.Transport = &oauth2.Transport{
		Base:   s.httpClient.Transport,
		Source: oauth2.ReuseTokenSource(nil, ts),
	}

	// If we're using public GitHub, make a regular client.
	// Otherwise, make an enterprise client.
	var isGHE bool = apiEndpoint != "https://api.github.com"
	ghClient, err := createGitHubClient(s.httpClient, apiEndpoint)
	if err != nil {
		s.log.Error(err, "error creating GitHub client")
	}
	s.apiClient = ghClient

	// TODO: this should support scanning users too

	specificScope := false

	if len(s.repos) > 0 {
		specificScope = true
	}

	var (
		ghUser *github.User
		resp   *github.Response
	)

	ctx.Logger().V(1).Info("Enumerating with token", "endpoint", apiEndpoint)
	for {
		ghUser, resp, err = s.apiClient.Users.Get(ctx, "")
		if handled := s.handleRateLimit(err, resp); handled {
			continue
		}
		if err != nil {
			return errors.New(err)
		}
		break
	}

	if s.orgsCache.Count() > 0 {
		specificScope = true
		for _, org := range s.orgsCache.Keys() {
			logger := s.log.WithValues("org", org)
			if err := s.getReposByOrg(ctx, org); err != nil {
				logger.Error(err, "error fetching repos for org")
			}

			if s.conn.ScanUsers {
				err := s.addMembersByOrg(ctx, org)
				if err != nil {
					logger.Error(err, "Unable to add members by org")
					continue
				}
			}
		}
	}

	// If no scope was provided, enumerate them.
	if !specificScope {
		if err := s.getReposByUser(ctx, ghUser.GetLogin()); err != nil {
			s.log.Error(err, "error fetching repos by user")
		}

		if isGHE {
			s.addAllVisibleOrgs(ctx)
		} else {
			// Scan for orgs is default with a token. GitHub App enumerates the repositories
			// that were assigned to it in GitHub App settings.
			s.addOrgsByUser(ctx, ghUser.GetLogin())
		}

		for _, org := range s.orgsCache.Keys() {
			logger := s.log.WithValues("org", org)
			if err := s.getReposByOrg(ctx, org); err != nil {
				logger.Error(err, "error fetching repos by org")
			}

			if err := s.getReposByUser(ctx, ghUser.GetLogin()); err != nil {
				logger.Error(err, "error fetching repos by user")
			}

			if s.conn.ScanUsers {
				err := s.addMembersByOrg(ctx, org)
				if err != nil {
					logger.Error(err, "Unable to add members by org for org")
				}
			}
		}

		// If we enabled ScanUsers above, we've already added the gists for the current user and users from the orgs.
		// So if we don't have ScanUsers enabled, add the user gists as normal.
		if err := s.addUserGistsToCache(ctx, ghUser.GetLogin()); err != nil {
			s.log.Error(err, "error fetching gists", "user", ghUser.GetLogin())
		}

		return nil
	}

	if s.conn.ScanUsers {
		s.log.Info("Adding repos", "members", len(s.members), "orgs", s.orgsCache.Count())
		s.addReposForMembers(ctx)
		return nil
	}

	return nil
}

func (s *Source) enumerateWithApp(ctx context.Context, apiEndpoint string, app *credentialspb.GitHubApp) (installationClient *github.Client, err error) {
	installationID, err := strconv.ParseInt(app.InstallationId, 10, 64)
	if err != nil {
		return nil, errors.New(err)
	}

	appID, err := strconv.ParseInt(app.AppId, 10, 64)
	if err != nil {
		return nil, errors.New(err)
	}

	// This client is used for most APIs.
	itr, err := ghinstallation.New(
		s.httpClient.Transport,
		appID,
		installationID,
		[]byte(app.PrivateKey))
	if err != nil {
		return nil, errors.New(err)
	}
	itr.BaseURL = apiEndpoint

	s.httpClient.Transport = itr
	s.apiClient, err = github.NewEnterpriseClient(apiEndpoint, apiEndpoint, s.httpClient)
	if err != nil {
		return nil, errors.New(err)
	}

	// This client is required to create installation tokens for cloning.
	// Otherwise, the required JWT is not in the request for the token :/
	appItr, err := ghinstallation.NewAppsTransport(
		s.httpClient.Transport,
		appID,
		[]byte(app.PrivateKey))
	if err != nil {
		return nil, errors.New(err)
	}
	appItr.BaseURL = apiEndpoint

	// Does this need to be separate from |s.httpClient|?
	instHttpClient := common.RetryableHttpClientTimeout(60)
	instHttpClient.Transport = appItr
	installationClient, err = github.NewEnterpriseClient(apiEndpoint, apiEndpoint, instHttpClient)
	if err != nil {
		return nil, errors.New(err)
	}

	// If no repos were provided, enumerate them.
	if len(s.repos) == 0 {
		if err = s.getReposByApp(ctx); err != nil {
			return nil, err
		}

		// Check if we need to find user repos.
		if s.conn.ScanUsers {
			err := s.addMembersByApp(ctx, installationClient)
			if err != nil {
				return nil, err
			}
			s.log.Info("Scanning repos", "org_members", len(s.members))
			for _, member := range s.members {
				logger := s.log.WithValues("member", member)
				if err := s.getReposByUser(ctx, member); err != nil {
					logger.Error(err, "error fetching gists by user")
				}
				if err := s.getReposByUser(ctx, member); err != nil {
					logger.Error(err, "error fetching repos by user")
				}
			}
		}
	}

	return installationClient, nil
}

func createGitHubClient(httpClient *http.Client, apiEndpoint string) (ghClient *github.Client, err error) {
	// If we're using public GitHub, make a regular client.
	// Otherwise, make an enterprise client.
	if apiEndpoint == "https://api.github.com" {
		ghClient = github.NewClient(httpClient)
	} else {
		ghClient, err = github.NewEnterpriseClient(apiEndpoint, apiEndpoint, httpClient)
		if err != nil {
			return nil, errors.New(err)
		}
	}

	return ghClient, err
}

func (s *Source) scan(ctx context.Context, installationClient *github.Client, chunksChan chan *sources.Chunk) error {
	var scanned uint64

	s.log.V(2).Info("Found repos to scan", "count", len(s.repos))

	// If there is resume information available, limit this scan to only the repos that still need scanning.
	reposToScan, progressIndexOffset := sources.FilterReposToResume(s.repos, s.GetProgress().EncodedResumeInfo)
	s.repos = reposToScan

	scanErrs := sources.NewScanErrors()
	// Setup scan options if it wasn't provided.
	if s.scanOptions == nil {
		s.scanOptions = &git.ScanOptions{}
	}

	for i, repoURL := range s.repos {
		i, repoURL := i, repoURL
		s.jobPool.Go(func() error {
			if common.IsDone(ctx) {
				return nil
			}

			// TODO: set progress complete is being called concurrently with i
			s.setProgressCompleteWithRepo(i, progressIndexOffset, repoURL)
			// Ensure the repo is removed from the resume info after being scanned.
			defer func(s *Source, repoURL string) {
				s.resumeInfoMutex.Lock()
				defer s.resumeInfoMutex.Unlock()
				s.resumeInfoSlice = sources.RemoveRepoFromResumeInfo(s.resumeInfoSlice, repoURL)
			}(s, repoURL)

			if !strings.HasSuffix(repoURL, ".git") {
				scanErrs.Add(fmt.Errorf("repo %s does not end in .git", repoURL))
				return nil
			}

			logger := s.log.WithValues("repo", repoURL)
			logger.V(2).Info(fmt.Sprintf("attempting to clone repo %d/%d", i+1, len(s.repos)))
			var path string
			var repo *gogit.Repository
			var err error

			path, repo, err = s.cloneRepo(ctx, repoURL, installationClient)
			if err != nil {
				scanErrs.Add(err)
				return nil
			}

			defer os.RemoveAll(path)
			if err != nil {
				scanErrs.Add(fmt.Errorf("error cloning repo %s: %w", repoURL, err))
				return nil
			}

			s.setScanOptions(s.conn.Base, s.conn.Head)

			repoSize := s.repoSizes.getRepo(repoURL)
			logger.V(2).Info(fmt.Sprintf("scanning repo %d/%d", i, len(s.repos)), "repo_size", repoSize)

			now := time.Now()
			defer func(start time.Time) {
				logger.V(2).Info(fmt.Sprintf("scanned %d/%d repos", scanned, len(s.repos)), "repo_size", repoSize, "duration_seconds", time.Since(start).Seconds())
			}(now)

			if err = s.scanComments(ctx, repoURL, chunksChan); err != nil {
				scanErrs.Add(fmt.Errorf("error scanning comments in repo %s: %w", repoURL, err))
				return nil
			}

			if err = s.git.ScanRepo(ctx, repo, path, s.scanOptions, chunksChan); err != nil {
				scanErrs.Add(fmt.Errorf("error scanning repo %s: %w", repoURL, err))
				return nil
			}
			atomic.AddUint64(&scanned, 1)

			return nil
		})
	}

	_ = s.jobPool.Wait()
	if scanErrs.Count() > 0 {
		s.log.V(2).Info("Errors encountered while scanning", "error-count", scanErrs.Count(), "errors", scanErrs)
	}
	s.SetProgressComplete(len(s.repos), len(s.repos), "Completed Github scan", "")

	return nil
}

// handleRateLimit returns true if a rate limit was handled
// Unauthenticated access to most github endpoints has a rate limit of 60 requests per hour.
// This will likely only be exhausted if many users/orgs are scanned without auth
func (s *Source) handleRateLimit(errIn error, res *github.Response) bool {
	limit, ok := errIn.(*github.RateLimitError)
	if !ok {
		return false
	}

	if res != nil {
		knownWait := true
		remaining, err := strconv.Atoi(res.Header.Get("x-ratelimit-remaining"))
		if err != nil {
			knownWait = false
		}
		resetTime, err := strconv.Atoi(res.Header.Get("x-ratelimit-reset"))
		if err != nil || resetTime == 0 {
			knownWait = false
		}

		if knownWait && remaining == 0 {
			waitTime := int64(resetTime) - time.Now().Unix()
			if waitTime > 0 {
				duration := time.Duration(waitTime+1) * time.Second
				s.log.V(2).Info("rate limited", "resumeTime", time.Now().Add(duration).String())
				time.Sleep(duration)
				return true
			}
		}
	}

	s.log.V(2).Info("handling rate limit (5 minutes retry)", "retry-after", limit.Message)
	time.Sleep(time.Minute * 5)
	return true
}

func (s *Source) addReposForMembers(ctx context.Context) {
	s.log.Info("Fetching repos from members", "members", len(s.members))
	for member := range s.memberCache {
		if err := s.addUserGistsToCache(ctx, member); err != nil {
			s.log.Info("Unable to fetch gists by user", "user", member, "error", err)
		}
		if err := s.getReposByUser(ctx, member); err != nil {
			s.log.Info("Unable to fetch repos by user", "user", member, "error", err)
		}
	}
}

// addUserGistsToCache collects all the gist urls for a given user,
// and adds them to the filteredRepoCache.
func (s *Source) addUserGistsToCache(ctx context.Context, user string) error {
	gistOpts := &github.GistListOptions{}
	logger := s.log.WithValues("user", user)
	for {
		gists, res, err := s.apiClient.Gists.List(ctx, user, gistOpts)
		if err == nil {
			res.Body.Close()
		}
		if handled := s.handleRateLimit(err, res); handled {
			continue
		}
		if err != nil {
			return fmt.Errorf("could not list gists for user %s: %w", user, err)
		}
		for _, gist := range gists {
			s.filteredRepoCache.Set(gist.GetID(), gist.GetGitPullURL())
		}
		if res == nil || res.NextPage == 0 {
			break
		}
		logger.V(2).Info("Listed gists", "page", gistOpts.Page, "last_page", res.LastPage)
		gistOpts.Page = res.NextPage
	}
	return nil
}

func (s *Source) addMembersByApp(ctx context.Context, installationClient *github.Client) error {
	opts := &github.ListOptions{
		PerPage: membersAppPagination,
	}

	// TODO: Check rate limit for this call.
	installs, _, err := installationClient.Apps.ListInstallations(ctx, opts)
	if err != nil {
		return fmt.Errorf("could not enumerate installed orgs: %w", err)
	}

	for _, org := range installs {
		if org.Account.GetType() != "Organization" {
			continue
		}
		if err := s.addMembersByOrg(ctx, *org.Account.Login); err != nil {
			return err
		}
	}

	return nil
}

func (s *Source) addAllVisibleOrgs(ctx context.Context) {
	s.log.V(2).Info("enumerating all visible organizations on GHE")
	// Enumeration on this endpoint does not use pages it uses a since ID.
	// The endpoint will return organizations with an ID greater than the given since ID.
	// Empty org response is our cue to break the enumeration loop.
	orgOpts := &github.OrganizationsListOptions{
		Since: 0,
		ListOptions: github.ListOptions{
			PerPage: defaultPagination,
		},
	}
	for {
		orgs, resp, err := s.apiClient.Organizations.ListAll(ctx, orgOpts)
		if err == nil {
			resp.Body.Close()
		}
		if handled := s.handleRateLimit(err, resp); handled {
			continue
		}
		if err != nil {
			s.log.Error(err, "Could not list all organizations")
			return
		}
		if len(orgs) == 0 {
			break
		}
		lastOrgID := *orgs[len(orgs)-1].ID
		s.log.V(2).Info(fmt.Sprintf("listed organization IDs %d through %d", orgOpts.Since, lastOrgID))
		orgOpts.Since = lastOrgID

		for _, org := range orgs {
			var name string
			if org.Name != nil {
				name = *org.Name
			} else if org.Login != nil {
				name = *org.Login
			} else {
				continue
			}
			s.orgsCache.Set(name, name)
			s.log.V(2).Info("adding organization for repository enumeration", "id", org.ID, "name", name)
		}
	}
}

func (s *Source) addOrgsByUser(ctx context.Context, user string) {
	orgOpts := &github.ListOptions{
		PerPage: defaultPagination,
	}
	logger := s.log.WithValues("user", user)
	for {
		orgs, resp, err := s.apiClient.Organizations.List(ctx, "", orgOpts)
		if err == nil {
			resp.Body.Close()
		}
		if handled := s.handleRateLimit(err, resp); handled {
			continue
		}
		if err != nil {
			logger.Error(err, "Could not list organizations")
			return
		}
		if resp == nil {
			break
		}
		logger.V(2).Info("Listed orgs", "page", orgOpts.Page, "last_page", resp.LastPage)
		for _, org := range orgs {
			if org.Login == nil {
				continue
			}
			s.orgsCache.Set(*org.Login, *org.Login)
		}
		if resp.NextPage == 0 {
			break
		}
		orgOpts.Page = resp.NextPage
	}
}

func (s *Source) addMembersByOrg(ctx context.Context, org string) error {
	opts := &github.ListMembersOptions{
		PublicOnly: false,
		ListOptions: github.ListOptions{
			PerPage: membersAppPagination,
		},
	}

	logger := s.log.WithValues("org", org)
	for {
		members, res, err := s.apiClient.Organizations.ListMembers(ctx, org, opts)
		if err == nil {
			defer res.Body.Close()
		}
		if handled := s.handleRateLimit(err, res); handled {
			continue
		}
		if err != nil || len(members) == 0 {
			return errors.New("Could not list organization members: account may not have access to list organization members")
		}
		if res == nil {
			break
		}
		logger.V(2).Info("Listed members", "page", opts.Page, "last_page", res.LastPage)
		for _, m := range members {
			usr := m.Login
			if usr == nil || *usr == "" {
				continue
			}
			if _, ok := s.memberCache[*usr]; !ok {
				s.memberCache[*usr] = struct{}{}
			}
		}
		if res.NextPage == 0 {
			break
		}
		opts.Page = res.NextPage
	}

	return nil
}

// setProgressCompleteWithRepo calls the s.SetProgressComplete after safely setting up the encoded resume info string.
func (s *Source) setProgressCompleteWithRepo(index int, offset int, repoURL string) {
	s.resumeInfoMutex.Lock()
	defer s.resumeInfoMutex.Unlock()

	// Add the repoURL to the resume info slice.
	s.resumeInfoSlice = append(s.resumeInfoSlice, repoURL)
	sort.Strings(s.resumeInfoSlice)

	// Make the resume info string from the slice.
	encodedResumeInfo := sources.EncodeResumeInfo(s.resumeInfoSlice)

	s.SetProgressComplete(index+offset, len(s.repos)+offset, fmt.Sprintf("Repo: %s", repoURL), encodedResumeInfo)
}

func (s *Source) scanComments(ctx context.Context, repoPath string, chunksChan chan *sources.Chunk) error {
	s.log.Info("scanning comments", "repository", repoPath)

	// Support ssh and https URLs
	repoURL, err := git.GitURLParse(repoPath)
	if err != nil {
		return err
	}

	trimmedURL := removeURLAndSplit(repoURL.String())
	owner := trimmedURL[1]
	repo := trimmedURL[2]

	var (
		sortType      = "created"
		directionType = "desc"
		allComments   = 0
	)

	if s.includeIssueComments {

		issueOpts := &github.IssueListCommentsOptions{
			Sort:      &sortType,
			Direction: &directionType,
			ListOptions: github.ListOptions{
				PerPage: defaultPagination,
				Page:    1,
			},
		}

		for {
			issueComments, resp, err := s.apiClient.Issues.ListComments(ctx, owner, repo, allComments, issueOpts)
			if s.handleRateLimit(err, resp) {
				break
			}

			if err != nil {
				return err
			}

			err = s.chunkIssueComments(ctx, repo, issueComments, chunksChan, repoPath)
			if err != nil {
				return err
			}

			issueOpts.ListOptions.Page++

			if len(issueComments) < defaultPagination {
				break
			}
		}

	}

	if s.includePRComments {
		prOpts := &github.PullRequestListCommentsOptions{
			Sort:      sortType,
			Direction: directionType,
			ListOptions: github.ListOptions{
				PerPage: defaultPagination,
				Page:    1,
			},
		}

		for {
			prComments, resp, err := s.apiClient.PullRequests.ListComments(ctx, owner, repo, allComments, prOpts)
			if s.handleRateLimit(err, resp) {
				break
			}

			if err != nil {
				return err
			}

			err = s.chunkPullRequestComments(ctx, repo, prComments, chunksChan, repoPath)
			if err != nil {
				return err
			}

			prOpts.ListOptions.Page++

			if len(prComments) < defaultPagination {
				break
			}
		}
	}

	return nil
}

func (s *Source) chunkIssueComments(ctx context.Context, repo string, comments []*github.IssueComment, chunksChan chan *sources.Chunk, repoPath string) error {
	for _, comment := range comments {
		// Create chunk and send it to the channel.
		chunk := &sources.Chunk{
			SourceName: s.name,
			SourceID:   s.SourceID(),
			SourceType: s.Type(),
			SourceMetadata: &source_metadatapb.MetaData{
				Data: &source_metadatapb.MetaData_Github{
					Github: &source_metadatapb.Github{
						Link:       sanitizer.UTF8(comment.GetHTMLURL()),
						Username:   sanitizer.UTF8(comment.GetUser().GetLogin()),
						Email:      sanitizer.UTF8(comment.GetUser().GetEmail()),
						Repository: sanitizer.UTF8(repo),
						Timestamp:  sanitizer.UTF8(comment.GetCreatedAt().String()),
						Visibility: s.visibilityOf(ctx, repoPath),
					},
				},
			},
			Data:   []byte(sanitizer.UTF8(comment.GetBody())),
			Verify: s.verify,
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case chunksChan <- chunk:
		}
	}
	return nil
}

func (s *Source) chunkPullRequestComments(ctx context.Context, repo string, comments []*github.PullRequestComment, chunksChan chan *sources.Chunk, repoPath string) error {
	for _, comment := range comments {
		// Create chunk and send it to the channel.
		chunk := &sources.Chunk{
			SourceName: s.name,
			SourceID:   s.SourceID(),
			SourceType: s.Type(),
			SourceMetadata: &source_metadatapb.MetaData{
				Data: &source_metadatapb.MetaData_Github{
					Github: &source_metadatapb.Github{
						Link:       sanitizer.UTF8(comment.GetHTMLURL()),
						Username:   sanitizer.UTF8(comment.GetUser().GetLogin()),
						Email:      sanitizer.UTF8(comment.GetUser().GetEmail()),
						Repository: sanitizer.UTF8(repo),
						Timestamp:  sanitizer.UTF8(comment.GetCreatedAt().String()),
					},
				},
			},
			Data:   []byte(sanitizer.UTF8(comment.GetBody())),
			Verify: s.verify,
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case chunksChan <- chunk:
		}
	}
	return nil
}

func removeURLAndSplit(url string) []string {
	trimmedURL := strings.TrimPrefix(url, "https://")
	trimmedURL = strings.TrimSuffix(trimmedURL, ".git")
	splitURL := strings.Split(trimmedURL, "/")

	return splitURL
}
