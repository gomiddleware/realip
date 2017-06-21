package realip

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type key int

const realIpKey key = 27 // doesn't matter what this value is

var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
var xRealIP = http.CanonicalHeaderKey("X-Real-IP")

// RealIp is middleware that parses either the X-Forwarded-For header or the X-Real-IP header (in that order). It stores the result
// in the request's context for later retrieval during the request lifecycle.
//
// If neither of these are available, we check the req.RemoteAddr for the original value. This may or may not contain
// what you really want, but it is the best guess we have.
//
// If none of these provide an IP address, we just store the empty string instead of an IP address.
//
// This middleware should be inserted fairly early in the middleware stack to ensure that subsequent layers (eg.,
// request loggers) which require the RealIp will see the intended value. Other GoMiddleware, such as
// gomiddleware/logger can use this information.
//
// You should only use this middleware if you can trust the headers passed to you (in particular, the two headers this
// middleware uses). If you have placed a reverse proxy like HAProxy or Nginx in front of your server and forwarded the
// correct headers then you should be fine. However, if your reverse proxy is configured to pass along arbitrary header
// values from the client, or if you use this middleware without a reverse proxy, malicious clients will be able to
// inject values that are not correct and may even cause a vulnerability if used incorrectly.
func RealIp(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// realIp can return ""
		ctx := context.WithValue(r.Context(), realIpKey, realIp(r))
		next.ServeHTTP(w, r.WithContext(ctx))

	}

	return http.HandlerFunc(fn)
}

// RealIpFromRequest can be used to obtain the RealIp from the request. This is given as a convenience method for
// RealIpFromContext if you don't have the context handy, but you do (of course) have the http.Request.
func RealIpFromRequest(r *http.Request) string {
	return r.Context().Value(realIpKey).(string)
}

// RealIpFromContext can be used to obtain the RealIp from the context (if you already have it handy).
func RealIpFromContext(ctx context.Context) string {
	return ctx.Value(realIpKey).(string)
}

// checkIp checks that the IP Address looks okay and returns it, else it returns the empty string.
func checkIP(rawIP string) string {
	// try the request connection
	ip := ""
	if strings.IndexRune(rawIP, ':') != -1 {
		var err error
		ip, _, err = net.SplitHostPort(rawIP)
		if err != nil {
			fmt.Errorf("checkIP: %q is not IP:port", rawIP)
			return ""
		}
	} else {
		ip = rawIP
	}

	// parse the ip to make sure it is okay
	userIP := net.ParseIP(ip)
	if userIP == nil {
		fmt.Errorf("checkIP: %q is not IP:port", rawIP)
		return ""
	}

	return userIP.String()
}

// realIp is an internal function to help with extracting the IP Address from the request. Hat tip to
// https://blog.golang.org/context for info regarding r.RemoteAddr and parsing it to make sure it is correct and
// reliable.
func realIp(r *http.Request) string {
	// firstly, try the "X-Forwarded-For"
	xff := r.Header.Get(xForwardedFor)
	if xff != "" {
		i := strings.Index(xff, ", ")
		if i == -1 {
			return checkIP(xff) // the complete string, since there is no comma
		}
		return checkIP(xff[:i])
	}

	// check "X-Real-IP" instead
	xrip := r.Header.Get(xRealIP)
	if xrip != "" {
		return checkIP(xrip)
	}

	// finally, try the request's RemoteAddr from the connection
	return checkIP(r.RemoteAddr)
}
