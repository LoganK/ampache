package ampache

import (
	"encoding/xml"
	"fmt"
)

type Ping struct {
	Handshake
	Server     string `xml:"server"`
	Version    string `xml:"version"`
	Compatible int    `xml:"compatible"`
}

func (c *Client) InvokePing() (*Ping, error) {
	resp, err := c.Invoke("ping", nil)
	if err != nil {
		return nil, fmt.Errorf("ping failure: %w", err)
	}
	defer resp.Close()

	var v Ping
	if err := xml.NewDecoder(resp).Decode(&v); err != nil {
		return nil, fmt.Errorf("unexpected ping response: %w", err)
	}

	return &v, nil
}
