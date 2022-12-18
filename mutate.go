package ampache

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

// Rate rates a library item
// https://ampache.org/api/api-xml-methods/#rate
// rateType may use the Media* constants (e.g., MediaSong)
func (c *Client) Rate(rateType string, id, rating int) (string, error) {
	resp, err := c.Invoke("rate", map[string]string{
		"type":   rateType,
		"id":     strconv.Itoa(id),
		"rating": strconv.Itoa(rating),
	})
	if err != nil {
		return "", fmt.Errorf("rate failure: %w", err)
	}
	defer resp.Close()

	var v Response
	if err := xml.NewDecoder(resp).Decode(&v); err != nil {
		return "", fmt.Errorf("unexpected rate response: %w", err)
	}
	if v.Error != nil {
		return "", v.Error
	}

	return v.Success.Message, nil
}
