package github

import (
	"net/http"

	gogit "github.com/go-git/go-git/v5"
	"github.com/google/go-github/v62/github"
	"github.com/trufflesecurity/trufflehog/v3/pkg/common"
	"github.com/trufflesecurity/trufflehog/v3/pkg/context"
	"github.com/trufflesecurity/trufflehog/v3/pkg/pb/credentialspb"
	"github.com/trufflesecurity/trufflehog/v3/pkg/sources/git"
)

type basicAuthConnector struct {
	httpClient *http.Client
	apiClient  *github.Client
	username   string
	password   string
}

var _ connector = (*basicAuthConnector)(nil)

func newBasicAuthConnector(apiEndpoint string, cred *credentialspb.BasicAuth) (*basicAuthConnector, error) {
	httpClient := common.RetryableHTTPClientTimeout(60)
	httpClient.Transport = &github.BasicAuthTransport{
		Username: cred.Username,
		Password: cred.Password,
	}

	apiClient, err := createGitHubClient(httpClient, apiEndpoint)
	if err != nil {
		return nil, err
	}

	return &basicAuthConnector{
		httpClient: httpClient,
		apiClient:  apiClient,
		username:   cred.Username,
		password:   cred.Password,
	}, nil
}

func (c basicAuthConnector) ApiClient() *github.Client {
	return c.apiClient
}

func (c basicAuthConnector) Clone(ctx context.Context, repoURL string) (string, *gogit.Repository, error) {
	return git.CloneRepoUsingToken(ctx, c.password, repoURL, c.username)
}

func (c basicAuthConnector) IsGithubEnterprise() bool {
	return false
}

func (c basicAuthConnector) HttpClient() *http.Client {
	return c.httpClient
}

func (c basicAuthConnector) ListAppInstallations(ctx context.Context) ([]*github.Installation, error) {
	return nil, nil
}
