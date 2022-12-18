# Ampache API Library

This is a small client library to support https://ampache.org/api along with a
CLI demo.

## Getting Started

```sh
$ go get github.com/logank/ampache/cmd/ampache_tool
$ ampache_tool https://example.com songs --limit=5
```

### Library

```golang
import (
	"github.com/logank/ampache"
)

func main() {
	c, _ := ampache.New("https://example.com")
	c.WithAuthPassword("user", "pass")
	p, _ := c.Ping()
}
```
