package ampache_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/logank/ampache"
)

func TestPingNoAuth(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/server/xml.server.php") && r.URL.Query().Get("action") == "ping" {
			io.WriteString(w, `<?xml version="1.0" encoding="UTF-8"?>
<root>
  <server><![CDATA[5.5.5-release]]></server>
  <version><![CDATA[5.5.5]]></version>
  <compatible><![CDATA[350001]]></compatible>
</root>`)
			return
		}
		t.Errorf("Unexpected request: %s", r.URL)
	}))
	defer svr.Close()

	c, err := ampache.New(svr.URL)
	if err != nil {
		t.Fatalf("Unexpected error New: %s", err)
	}

	v, err := c.InvokePing()
	if err != nil {
		t.Fatalf("Unexpected error InvokePing: %s", err)
	}

	want := "5.5.5-release"
	if v.Server != want {
		t.Errorf("Unexpected <server> = '%s'; want '%s'", v.Server, want)
	}
}

func TestPing(t *testing.T) {
	user := "test_user"
	password := "test_password"
	var timestamp int64 = 1671325631
	ampache.Now = func() time.Time { return time.Unix(timestamp, 0) }

	// TIME=1671325631
	// KEY=$(echo -n "test_password" | openssl sha256 - | sed -e 's|.* ||')
	// echo -n ${TIME}${KEY} | openssl sha256 - | sed -e 's|.* ||'
	passphrase := "78dddc4fafaea6653b54c02e9bb285ad0f4f8560eef4d0088e7fa672c5293b73"

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("action") == "handshake" &&
			q.Get("timestamp") == strconv.FormatInt(timestamp, 10) &&
			q.Get("user") == user &&
			q.Get("auth") == passphrase {
			io.WriteString(w, `<?xml version="1.0" encoding="UTF-8"?>
<root>
  <auth><![CDATA[cfj3f237d563f479f5223k23189dbb34]]></auth>
  <api><![CDATA[6.0.0]]></api>
  <session_expire><![CDATA[2022-12-18T01:12:11+00:00]]></session_expire>
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
</root>`)
			return
		} else if q.Get("action") == "ping" &&
			q.Get("auth") == "cfj3f237d563f479f5223k23189dbb34" {
			io.WriteString(w, `<?xml version="1.0" encoding="UTF-8"?>
<root>
  <session_expire><![CDATA[2022-12-18T01:12:37+00:00]]></session_expire>
  <server><![CDATA[develop]]></server>
  <version><![CDATA[6.0.0]]></version>
  <compatible><![CDATA[350001]]></compatible>
  <auth><![CDATA[cfj3f237d563f479f5223k23189dbb34]]></auth>
  <api><![CDATA[6.0.0]]></api>
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
</root>`)
			return
		}
		t.Errorf("Unexpected request: %s", r.URL)
	}))
	defer svr.Close()

	c, err := ampache.New(svr.URL)
	if err != nil {
		t.Fatalf("Unexpected error New: %s", err)
	}
	c.WithAuthPassword(user, password)

	v, err := c.InvokePing()
	if err != nil {
		t.Fatalf("Unexpected error InvokePing: %s", err)
	}

	want := time.Date(2022, 12, 18, 1, 12, 37, 0, time.UTC)
	if !v.SessionExpire.Equal(want) {
		t.Errorf("Unexpected <session_expire> = '%s'; want '%s'", v.SessionExpire, want)
	}
}
