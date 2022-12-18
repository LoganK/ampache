package ampache

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"time"
)

const (
	MediaSong           string = "song"
	MediaAlbum                 = "album"
	MediaArtist                = "artist"
	MediaPlaylist              = "playlist"
	MediaPodcast               = "podcast"
	MediaPodcastEpisode        = "podcast_episode"
	MediaVideo                 = "video"
	MediaTvShow                = "tvshow"
	MediaTvShowSeason          = "tvshow_season"
)

type xmlTime struct {
	time.Time
}

func (c *xmlTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse("2006-01-02T15:04:05-07:00", v)
	if err != nil {
		return err
	}
	*c = xmlTime{parse}
	return nil
}

type xmlDuration struct {
	time.Duration
}

func (c *xmlDuration) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v int
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	c.Duration = time.Duration(v) * time.Second
	return nil
}

type xmlURL struct {
	url.URL
}

func (c *xmlURL) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	u, err := url.Parse(v)
	if err != nil {
		return err
	}
	c.URL = *u
	return nil
}

type Success struct {
	Code    int    `xml:"code,attr"`
	Message string `xml:",chardata"`
}

type Error struct {
	Action  string `xml:"errorAction"`
	Code    int    `xml:"errorCode,attr"`
	Type    string `xml:"errorType"`
	Message string `xml:"errorMessage"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s [%d]: '%s' for '%s'", e.Action, e.Code, e.Message, e.Type)
}

type Response struct {
	Success *Success `xml:"success"`
	Error   *Error   `xml:"error"`
}

type Artist struct {
	Id   int    `xml:"id,attr"`
	Name string `xml:",chardata"`
}

type Album struct {
	Id   int    `xml:"id,attr"`
	Name string `xml:",chardata"`
}

type Song struct {
	Id              int         `xml:"id,attr"`
	Title           string      `xml:"title"`
	Name            string      `xml:"name"`
	Artist          Artist      `xml:"artist"`
	Album           Album       `xml:"album"`
	AlbumArtist     Artist      `xml:"albumartist"`
	Disk            int         `xml:"disk"`
	Track           int         `xml:"track"`
	Genre           string      `xml:"genre"`
	Filename        string      `xml:"filename"`
	PlaylistTrack   int         `xml:"playlisttrack"`
	Time            xmlDuration `xml:"time"`
	Year            int         `xml:"year"`
	BitRate         int         `xml:"bitrate"`
	Rate            int         `xml:"rate"`
	Mode            string      `xml:"mode"`
	Mime            string      `xml:"mime"`
	Url             xmlURL      `xml:"url"`
	Size            int         `xml:"size"`
	Mbid            string      `xml:"mbid"`
	AlbumMbid       string      `xml:"album_mbid"`
	ArtistMbid      string      `xml:"artist_mbid"`
	AlbumArtistMbid string      `xml:"albumartist_mbid"`
	Art             xmlURL      `xml:"art"`
	Flag            int         `xml:"flag"`
	PreciseRating   int         `xml:"preciserating"`
	Rating          int         `xml:"rating"`
	AverageRating   string      `xml:"averagerating"`
	PlayCount       int         `xml:"playcount"`
	Catalog         int         `xml:"catalog"`
	Composer        string      `xml:"composer"`
	Channels        int         `xml:"channels"`
	Comment         string      `xml:"comment"`
	License         string      `xml:"license"`
	Publisher       string      `xml:"publisher"`
	Language        string      `xml:"language"`
	Lyrics          string      `xml:"lyrics"`
	AlbumGain       float64     `xml:"replaygain_album_gain"`
	AlbumPeak       float64     `xml:"replaygain_album_peak"`
	TrackGain       float64     `xml:"replaygain_track_gain"`
	TrackPeak       float64     `xml:"replaygain_track_peak"`
	R128AlbumGain   int         `xml:"r128_album_gain"`
	R128TrackGain   int         `xml:"r128_track_gain"`
}

type Songs struct {
	TotalCount int    `xml:"total_count"`
	Songs      []Song `xml:"song"`
}
