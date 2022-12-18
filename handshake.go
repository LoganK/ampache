package ampache

import (
	"encoding/xml"
	"fmt"
)

type Handshake struct {
	Auth            string  `xml:"auth"`
	Api             string  `xml:"api"`
	SessionExpire   xmlTime `xml:"session_expire"`
	Update          xmlTime `xml:"update"`
	Add             xmlTime `xml:"add"`
	Clean           xmlTime `xml:"clean"`
	Songs           int     `xml:"songs"`
	Albums          int     `xml:"albums"`
	Artists         int     `xml:"artists"`
	Genres          int     `xml:"genres"`
	Playlists       int     `xml:"playlists"`
	Users           int     `xml:"users"`
	Catalogs        int     `xml:"catalogs"`
	Videos          int     `xml:"videos"`
	Podcast         int     `xml:"podcasts"`
	PodcastEpisodes int     `xml:"podcast_episodes"`
	Shares          int     `xml:"shares"`
	Licenses        int     `xml:"licenses"`
	LiveStreams     int     `xml:"live_streams"`
	Labels          int     `xml:"labels"`
}

func (c *Client) Handshake() (*Handshake, error) {
	resp, err := c.Invoke("handshake", nil)
	if err != nil {
		return nil, fmt.Errorf("handshake failure: %w", err)
	}
	defer resp.Close()

	var v Handshake
	if err := xml.NewDecoder(resp).Decode(&v); err != nil {
		return nil, fmt.Errorf("unexpected handshake response: %w", err)
	}

	return &v, nil
}
