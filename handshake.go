package ampache

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
