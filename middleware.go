package shortly

import (
	"net"
	"net/http"
)

func RateLimiterMiddleware(store RateLimiterAllow) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}
			allow, err := store.Allow(ip)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			if !allow {
				WriteErrorResponse(w, http.StatusTooManyRequests, "too many requests", nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
