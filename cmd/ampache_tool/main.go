package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/logank/ampache"
)

const (
	progName = "ampache_tool"
)

func main() {
	host, action, input, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed parse args: %s\n", err)
		fmt.Fprintf(os.Stderr, "Use the form: %s https://example.com songs --limit=5\n", progName)
		os.Exit(-1)
	}

	user, pass := os.Getenv("AMPACHE_USER"), os.Getenv("AMPACHE_PASS")
	if action != "ping" && (user == "" || pass == "") {
		fmt.Fprint(os.Stderr, "Credentials not found.\n")
		fmt.Fprint(os.Stderr, "Set user and password via the environment\n")
		fmt.Fprintf(os.Stderr, "  AMPACHE_USER=user AMPACHE_PASS=pass %s\n", progName)
		os.Exit(-2)

	}

	c, err := ampache.New(host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed creating client: %s\n", err)
		os.Exit(-3)
	}

	if v, found := os.LookupEnv("AMPACHE_VERBOSE"); found {
		if vi, err := strconv.Atoi(v); err == nil {
			c.Verbose = vi
		}
	}

	if user != "" {
		c.WithAuthPassword(user, pass)
	}

	var v interface{}
	err = nil
	if action == "ping" {
		v, err = c.Ping()
	} else if action == "songs" {
		songs, serr := c.Songs(input)
		for _, s := range songs.Songs {
			fmt.Printf("%+v\n", s)
		}
		v = songs
		err = serr
	} else {
		resp, err := c.Invoke(action, input)
		if err == nil {
			defer resp.Close()

			var out strings.Builder
			io.Copy(&out, resp)
			v = out.String()
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed calling server: %s\n", err)
		fmt.Fprintf(os.Stderr, "Last response: '%s'\n", c.LastResponse())
		os.Exit(-4)
	}
	fmt.Printf("%+v\n", v)
}
