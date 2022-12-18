package ampache

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type Client struct {
	Verbose      bool
	lastResponse *strings.Builder

	host   *url.URL
	key    string
	client http.Client

	// Used to modify the request (i.e., to add authentication)
	wrapValues func(url.Values) error
}

// Create a new client connection. The given host should be the base URL to
// reach Ampache. The caller is expected to call the appropriate WithAuth*
// method before the client will be fully functional.
func New(host string) (*Client, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("unable to parse host: %w", err)
	}

	u.Path = path.Join(u.Path, "server/xml.server.php")

	return &Client{
		host:       u,
		wrapValues: func(url.Values) error { return nil },
	}, nil
}

// Invoke calls the API for the given action. Parameters may be passed in input.
func (c *Client) Invoke(action string, input map[string]string) (io.ReadCloser, error) {
	params := make(url.Values)
	for k, v := range input {
		params.Set(k, v)
	}
	params.Set("action", action)

	if err := c.wrapValues(params); err != nil {
		return nil, fmt.Errorf("failed creating request: %w", err)
	}

	return c.invokeInternal(params)
}

// InvokeRaw calls the API but skips any additional processing such as authentication.
func (c *Client) InvokeRaw(input map[string]string) (io.ReadCloser, error) {
	params := make(url.Values)
	for k, v := range input {
		params.Set(k, v)
	}

	return c.invokeInternal(params)
}

// LastResponse gets the last response body. Only works if Verbose is true.
func (c *Client) LastResponse() string {
	return c.lastResponse.String()
}

func (c *Client) invokeInternal(params url.Values) (io.ReadCloser, error) {
	// Deep copy. This can't fail.
	req, err := url.Parse(c.host.String())
	if err != nil {
		log.Fatalf("Failed on internal URL copy: %s", err)
	}
	req.RawQuery = params.Encode()

	resp, err := c.client.Get(req.String())
	if err != nil {
		return nil, fmt.Errorf("failed during API call: %w", err)
	}

	if c.Verbose {
		log.Printf("%s [%d]", req.String(), resp.StatusCode)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected response during API call: %w", err)
	}

	if c.Verbose {
		c.lastResponse = &strings.Builder{}
		io.TeeReader(resp.Body, c.lastResponse)
	}

	return resp.Body, nil
}
