package ampache_test

import (
	"encoding/xml"
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
