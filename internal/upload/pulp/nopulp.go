//go:build nopulp

package pulp

type Client struct {
}

type Credentials struct {
	Username string
	Password string
}

func NewClient(url string, creds *Credentials) *Client {
	panic("pulp integration is not supported on this platform")
}

func NewClientFromFile(url, path string) (*Client, error) {
	panic("pulp integration is not supported on this platform")
}

func (cl *Client) UploadAndDistributeCommit(archivePath, repoName, basePath string) (string, error) {
	panic("pulp integration is not supported on this platform")
}
