package ampache_test

import (
	"encoding/xml"
	"net/url"
	"testing"
	"time"

	"github.com/logank/ampache"
)

func TestHandshake(t *testing.T) {
	// https://raw.githubusercontent.com/ampache/python3-ampache/api6/docs/xml-responses/handshake.xml
	raw := `<?xml version="1.0" encoding="UTF-8"?>
<root>
  <auth><![CDATA[cfj3f237d563f479f5223k23189dbb34]]></auth>
  <api><![CDATA[6.0.0]]></api>
  <session_expire><![CDATA[2022-08-17T04:34:55+00:00]]></session_expire>
  <update><![CDATA[2021-07-21T02:51:36+00:00]]></update>
  <add><![CDATA[2021-08-03T00:04:14+00:00]]></add>
  <clean><![CDATA[2021-08-03T00:05:54+00:00]]></clean>
  <songs><![CDATA[75]]></songs>
  <albums><![CDATA[9]]></albums>
  <artists><![CDATA[17]]></artists>
  <genres><![CDATA[7]]></genres>
  <playlists><![CDATA[2]]></playlists>
  <searches><![CDATA[17]]></searches>
  <playlists_searches><![CDATA[19]]></playlists_searches>
  <users><![CDATA[4]]></users>
  <catalogs><![CDATA[4]]></catalogs>
  <videos><![CDATA[2]]></videos>
  <podcasts><![CDATA[2]]></podcasts>
  <podcast_episodes><![CDATA[13]]></podcast_episodes>
  <shares><![CDATA[2]]></shares>
  <licenses><![CDATA[14]]></licenses>
  <live_streams><![CDATA[3]]></live_streams>
  <labels><![CDATA[3]]></labels>
</root>`

	var v ampache.Handshake
	if err := xml.Unmarshal([]byte(raw), &v); err != nil {
		t.Fatalf("Unmarshal failed for Handshake: %s input =\n%s", err, raw)
	}

	wantAuth := "cfj3f237d563f479f5223k23189dbb34"
	if v.Auth != wantAuth {
		t.Errorf("Unexpected <auth> = '%s'; want '%s'", v.Auth, wantAuth)
	}

	wantSe := time.Date(2022, 8, 17, 4, 34, 55, 0, time.UTC)
	if !v.SessionExpire.Equal(wantSe) {
		t.Errorf("Unexpected <session_expire> = '%s'; want '%s'", v.SessionExpire, wantSe)
	}

	wantSongs := 75
	if v.Songs != wantSongs {
		t.Errorf("Unexpected <songs> = %d; want %d", v.Songs, wantSongs)
	}
}

func TestSong(t *testing.T) {
	raw := `<?xml version="1.0" encoding="UTF-8"?>
<song id="8195">
  <title><![CDATA[Duality]]></title>
  <name><![CDATA[Duality]]></name>
  <artist id="400"><![CDATA[Bayside]]></artist>
  <album id="691"><![CDATA[The Walking Wounded]]></album>
  <albumartist id="0"><![CDATA[]]></albumartist>
  <disk><![CDATA[1]]></disk>
  <track>3</track>
  <genre id="36"><![CDATA[Alternative & Punk]]></genre>
  <filename><![CDATA[/media/iTunes/iTunes Music/Bayside/The Walking Wounded/03 Duality.m4a]]></filename>
  <playlisttrack>2</playlisttrack>
  <time>180</time>
  <year>2007</year>
  <bitrate>162837</bitrate>
  <rate>44100</rate>
  <mode><![CDATA[vbr]]></mode>
  <mime><![CDATA[audio/mp4]]></mime>
  <url><![CDATA[https://example.com/play/index.php?ssid=35e1bb81dc07bc36862a9a042b470b67&type=song&oid=8195&uid=2&player=api&name=Bayside%20-%20Duality.m4a]]></url>
  <size>3707354</size>
  <mbid><![CDATA[]]></mbid>
  <album_mbid><![CDATA[]]></album_mbid>
  <artist_mbid><![CDATA[]]></artist_mbid>
  <albumartist_mbid><![CDATA[]]></albumartist_mbid>
  <art><![CDATA[https://example.com/image.php?object_id=691&object_type=album&auth=35e1bb81dc07bc36862a9a042b470b67]]></art>
  <flag>0</flag>
  <preciserating/>
  <rating>5</rating>
  <averagerating/>
  <playcount>0</playcount>
  <catalog>1</catalog>
  <composer><![CDATA[Bayside]]></composer>
  <channels/>
  <comment><![CDATA[]]></comment>
  <license><![CDATA[]]></license>
  <publisher><![CDATA[]]></publisher>
  <language/>
  <lyrics><![CDATA[]]></lyrics>
  <replaygain_album_gain/>
  <replaygain_album_peak/>
  <replaygain_track_gain/>
  <replaygain_track_peak/>
  <r128_album_gain/>
  <r128_track_gain/>
</song>`

	var v ampache.Song
	if err := xml.Unmarshal([]byte(raw), &v); err != nil {
		t.Fatalf("Unmarshal failed for Song: %s input =\n%s", err, raw)
	}

	wantId := 8195
	if v.Id != wantId {
		t.Errorf("Unexpected <song id> = %d; want %d", v.Id, wantId)
	}

	wantTime := time.Duration(180 * time.Second)
	if v.Time.Duration != wantTime {
		t.Errorf("Unexpected <time> = %d; want %d", v.Time, wantTime)
	}

	wantArt, _ := url.Parse("https://example.com/image.php?object_id=691&object_type=album&auth=35e1bb81dc07bc36862a9a042b470b67")
	if v.Art.URL.String() != wantArt.String() {
		t.Errorf("Unexpected <art> = '%s'; want '%s'", v.Art.URL.String(), wantArt)
	}
}
