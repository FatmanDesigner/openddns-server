package main

import (
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// Serve OpenDDNS server at a given port
func Serve(port int) {
	http.HandleFunc("/ping", pingHandler)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func pingHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		res.WriteHeader(http.StatusBadRequest)
		io.WriteString(res, "Bad Request")
		return
	}

	appid := req.URL.Query().Get("appid")
	if appid == "" {
		res.WriteHeader(http.StatusBadRequest)
		io.WriteString(res, "Bad Request")
		return
	}

	if req.ContentLength == 0 {
		res.WriteHeader(http.StatusBadRequest)
		io.WriteString(res, "Bad Request")
		return
	}

	body := make([]byte, req.ContentLength)
	req.Body.Read(body)

	_splat := strings.SplitAfterN(string(body), "\n", 2)
	secret, domainName := strings.TrimSpace(_splat[0]), strings.TrimSpace(_splat[1])

	if len(secret) == 0 {
		res.WriteHeader(http.StatusForbidden)
		io.WriteString(res, "Request not authorized")
		return
	}

	// TODO: Honor appid and secret as security measure

	if len(domainName) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		io.WriteString(res, "Invalid domain name")
		return
	}

	// TODO: Validate domainName as per RFC 1034 - IETF (see https://www.ietf.org/rfc/rfc1034.txt)
	// TODO: Take into account req.Headers().Get("X-Forwarded-For")
	// .. when router is behind proxies
	host, _, _ := net.SplitHostPort(req.RemoteAddr)
	ip := net.ParseIP(host)
	Register(domainName, ip.String())

	res.WriteHeader(http.StatusOK)
	io.WriteString(res, "OK")
}
