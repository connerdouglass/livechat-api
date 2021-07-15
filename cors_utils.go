package main

import (
	"net/http"
	"strings"
)

func simplifyOrigin(origin string) string {
	origin = strings.TrimRight(origin, "/")
	origin = strings.ToLower(origin)
	return origin
}

func checkOrigin(allowedOrigins []string) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		origin = simplifyOrigin(origin)
		for _, allowed := range allowedOrigins {
			if simplifyOrigin(allowed) == origin {
				return true
			}
		}
		return false
	}
}
