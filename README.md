<p align="center">
  <img alt="GoReleaser Logo" src="https://storage.googleapis.com/trufflehog-static-sources/pixel_pig.png" height="140" />
  <h2 align="center">TruffleHog</h2>
  <p align="center">Find leaked credentials.</p>
</p>

---


[![CI Status](https://github.com/trufflesecurity/trufflehog/actions/workflows/release.yml/badge.svg)](https://github.com/trufflesecurity/trufflehog/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/trufflesecurity/trufflehog)](https://goreportcard.com/report/github.com/trufflesecurity/trufflehog)
[![Docker Hub Build Status](https://img.shields.io/docker/cloud/build/trufflesecurity/trufflehog.svg)](https://hub.docker.com/r/trufflesecurity/trufflehog/)
![GitHub](https://img.shields.io/github/license/trufflesecurity/trufflehog)

---

## Join The Slack
Have questions? Feedback? Jump in slack and hang out with us

https://join.slack.com/t/trufflehog-community/shared_invite/zt-pw2qbi43-Aa86hkiimstfdKH9UCpPzQ


## Demo

![GitHub scanning demo](https://storage.googleapis.com/truffle-demos/non-interactive.svg)

```bash
docker run -it -v "$PWD:/pwd" trufflesecurity/trufflehog:latest github --org=trufflesecurity
```

# What's new in v3?

TruffleHog v3 is a complete rewrite in Go with many new powerful features.

- We've **added over 600 credential detectors that support active verification against their respective APIs**.
- We've also added native **support for scanning GitHub, GitLab, filesystems, and S3**.


## What is credential verification?
For every potential credential that is detected, we've painstakingly implemented programatic verification against the API that we think it belongs to. Verification eliminates false positives. For example, the [AWS credential detector](pkg/detectors/aws/aws.go) performs a `GetCallerIdentity` API call against the AWS API to verify if an AWS credential is active.

## Installation

Several options:

### 1. Go
```
git clone https://github.com/trufflesecurity/trufflehog.git

cd trufflehog; go install
```

### 2. [Release binaries](https://github.com/trufflesecurity/trufflehog/releases)

### 3. Docker


> Note: Apple M1 hardware users should run with `docker run --platform linux/arm64` for better performance.

#### **Most users**

```bash
docker run -it -v "$PWD:/pwd" trufflesecurity/trufflehog:latest github --repo https://github.com/trufflesecurity/test_keys
```

#### **Apple M1 users**

The `linux/arm64` image is better to run on the M1 than the amd64 image.
Even better is running the native darwin binary avilable, but there is not container image for that.

```bash
docker run --platform linux/arm64 -it -v "$PWD:/pwd" trufflesecurity/trufflehog:latest github --repo https://github.com/trufflesecurity/test_keys 
```

### 4. Pip (help wanted)

It's possible to distribute binaries in pip wheels.

Here is an example of a [project that does it](https://github.com/Yelp/dumb-init).

Help with setting up this packaging would be appreciated!

### 5. Brew (help wanted)

We'd love to distribute via brew and could use your help.

## Usage

TruffleHog has a sub-command for each source of data that you may want to scan:

- git
- github
- gitlab
- S3
- filesystem
- file and stdin

Each subcommand can have options that you can see with the `-h` flag provided to the sub command:

```
$ trufflehog git --help
usage: TruffleHog git [<flags>] <uri>

Find credentials in git repositories.

Flags:
      --help           Show context-sensitive help (also try --help-long and --help-man).
      --debug          Run in debug mode
      --json           Output in JSON format.
      --concurrency=8  Number of concurrent workers.
      --verification   Verify the results.
  -i, --include_paths=INCLUDE_PATHS  
                       Path to file with newline separated regexes for files to include in scan.
  -x, --exclude_paths=EXCLUDE_PATHS  
                       Path to file with newline separated regexes for files to exclude in scan.
      --branch=BRANCH  Branch to scan.
      --allow          No-op flag for backwards compat.
      --entropy        No-op flag for backwards compat.
      --regex          No-op flag for backwards compat.

Args:
  <uri>  Git repository URL. https:// or file:// schema expected.
```

For example, to scan a  `git` repository, start with

```
$ trufflehog git https://github.com/trufflesecurity/trufflehog.git
```


#### Scanning an orginization

Try scanning an entire GitHub orginization with the following:

```bash
docker run -it -v "$PWD:/pwd" trufflesecurity/trufflehog:latest github --org=trufflesecurity
```


## Contributors

This project exists thanks to all the people who contribute. [[Contribute](CONTRIBUTING.md)].


<a href="https://github.com/trufflesecurity/trufflehog/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=trufflesecurity/trufflehog" />
</a>


## Contributing

Contributions are very welcome! Please see our [contribution guidelines first](CONTRIBUTING.md).

We no longer accept contributions to TruffleHog v2, but that code is available in the `v2` branch.

### Adding new secret detectors

We have published some [documentation and tooling to get started on adding new secret detectors](hack/docs/Adding_Detectors_external.md). Let's improve detection together!

## License Change

Since v3.0, TruffleHog is released under a AGPL 3 license, included in [`LICENSE`](LICENSE). TruffleHog v3.0 uses none of the previous codebase, but care was taken to preserve backwards compatibility on the command line interface. The work previous to this release is still available licensed under GPL 2.0 in the history of this repository and the previous package releases and tags. A completed CLA is required for us to accept contributions going forward.
