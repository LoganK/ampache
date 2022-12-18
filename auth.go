package ampache

import (
	"crypto/sha256"
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

var (
	// Useful for testing
	Now = func() time.Time { return time.Now() }
)

func (c *Client) WithAuthPassword(user, pass string) *Client {
	// https://ampache.org/api/#sending-handshake-request

	keyHash := sha256.New()
	_, _ = keyHash.Write([]byte(pass))
	key := fmt.Sprintf("%x", keyHash.Sum(nil))

	auth := struct {
		user string
		key  string

		auth   string
		expire time.Time
	}{
		user: user,
		key:  key,
	}
	newPassphrase := func() (timestamp, passphrase string) {
		timestamp = strconv.FormatInt(Now().Unix(), 10)
		hash := sha256.New()
		_, _ = hash.Write([]byte(timestamp))
		_, _ = hash.Write([]byte(auth.key))
		passphrase = fmt.Sprintf("%x", hash.Sum(nil))
		return
	}

	c.wrapValues = func(v url.Values) error {
		remaining := -Now().Sub(auth.expire)
		if remaining > 0 && remaining < 60*time.Second {
			auth.expire, _ = extendSession(c, auth.auth)
		}

		remaining = -Now().Sub(auth.expire)
		if remaining > 60*time.Second {
			v.Set("auth", auth.auth)
			return nil
		}

		ts, pp := newPassphrase()
		newAuth, newExpire, err := createSession(c, map[string]string{
			"auth":      pp,
			"timestamp": ts,
			"user":      auth.user,
		})
		if err != nil {
			return fmt.Errorf("error creating new session: %w", err)
		}
		if -Now().Sub(newExpire) < 1*time.Second {
			return errors.New("new session is invalid: " + newExpire.String())
		}

		auth.auth = newAuth
		auth.expire = newExpire
		v.Set("auth", auth.auth)
		return nil
	}

	return c
}

// createSession returns an auth string and expiration using the given params.
func createSession(c *Client, params map[string]string) (string, time.Time, error) {
	params["action"] = "handshake"
	resp, err := c.InvokeRaw(params)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("createSession failure: %w", err)
	}
	defer resp.Close()

	var v Handshake
	if err := xml.NewDecoder(resp).Decode(&v); err != nil {
		return "", time.Time{}, fmt.Errorf("unexpected handshake response: %w", err)
	}

	return v.Auth, v.SessionExpire.Time, nil
}

// extendSession makes a request to the server to extend the given session. If
// successful, returns the new expiration time.
func extendSession(c *Client, auth string) (time.Time, error) {
	resp, err := c.InvokeRaw(map[string]string{
		"action": "ping",
		"auth":   auth,
	})
	if err != nil {
		return time.Time{}, fmt.Errorf("extendSession failure: %w", err)
	}
	defer resp.Close()

	var v Ping
	if err := xml.NewDecoder(resp).Decode(&v); err != nil {
		return time.Time{}, fmt.Errorf("unexpected ping response: %w", err)
	}

	return v.SessionExpire.Time, nil
}
