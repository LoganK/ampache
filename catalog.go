package ampache

import (
	"encoding/xml"
	"fmt"
)

func (c *Client) Songs(input map[string]string) (Songs, error) {
	resp, err := c.Invoke("songs", input)
	if err != nil {
		return Songs{}, fmt.Errorf("songs failure: %w", err)
	}
	defer resp.Close()

	var v Songs
	if err := xml.NewDecoder(resp).Decode(&v); err != nil {
		return Songs{}, fmt.Errorf("unexpected songs response: %w", err)
	}

	return v, nil
}
