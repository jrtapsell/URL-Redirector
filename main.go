package main

import (
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
)

func securityHeaders(header http.Header) {
	header.Add("Content-Security-Policy", "default-src 'none'")
	header.Add("X-Content-Type-Options", "nosniff")
	header.Add("X-Frame-Options", "DENY")
	header.Add("X-Xss-Protection", "1; mode=block")
	header.Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
}

var REGEX, _ = regexp.Compile(`https?://\S+`)
type ReplaceArgs struct {
	tag string
	value string
}

func replace(start string, actions ...ReplaceArgs) string {
	tmp := start
	for _, action := range actions {
		tmp = strings.ReplaceAll(tmp, action.tag, action.value)
	}
	return tmp
}
func onRequest(w http.ResponseWriter, r *http.Request) {
	var header = w.Header()
	securityHeaders(header)

	requestedHost := r.Header.Get("Host")
	name := strings.Join([]string{"_redirect_", requestedHost}, ".")
	req, err := net.LookupTXT(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var url string
	for _, element := range req {
		if REGEX.MatchString(element) {
			url = element
			break
		}
	}

	path := r.URL.Path

	qs := r.URL.RawQuery
	if len(qs) > 0 {
		qs = "?" + qs
	}

	replaced := replace(url,
		ReplaceArgs{"|{p}", path},
		ReplaceArgs{"|{q}", qs},
	)

	header.Add("Location", replaced)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	http.HandleFunc("/", onRequest)
	println("Starting server on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
