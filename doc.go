/*

Package realip provides standard http.Handler middleware that sets a value in the Request's context to the result of
parsing either the X-Forwarded-For header or the X-Real-IP header (in that order).

*/
package realip
