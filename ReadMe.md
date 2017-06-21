# GoMiddleware : IP #

Middleware that sets a RealIP in your request's context (r.Context()).

* [Project](https://github.com/gomiddleware/ip)
* [GoDoc](https://godoc.org/github.com/gomiddleware/ip)

## Synopsis ##

```go
package main

import (
	"net/http"

	"github.com/gomiddleware/realip"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.URL.Path))
}

func main() {
    handle := realip(http.HandlerFunc(handler))
	http.Handle("/", handle)
	http.ListenAndServe(":8080", nil)
}
```

## Author ##

By [Andrew Chilton](https://chilts.org/), [@twitter](https://twitter.com/andychilton).

For [AppsAttic](https://appsattic.com/), [@AppsAttic](https://twitter.com/AppsAttic).

## License ##

MIT.

(Ends)
