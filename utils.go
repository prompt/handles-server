package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func HostnameToHandle(hostname string) (Handle, error) {
	domainParser := regexp.MustCompile(`^(.+?)\.(.+)$`)

	parts := domainParser.FindStringSubmatch(strings.ToLower(hostname))

	if len(parts) != 3 {
		return Handle{}, fmt.Errorf("Handle could not be parsed from hostname %s", hostname)
	}

	return Handle{
		Domain:   Domain(parts[2]),
		Username: Username(parts[1]),
	}, nil
}

type URLTemplate string

func URLFromTemplate(
	template URLTemplate,
	request *http.Request,
	handle Handle,
	did DecentralizedID,
) string {
	replacements := map[string]string{
		"{handle}":          handle.String(),
		"{did}":             string(did),
		"{handle.domain}":   string(handle.Domain),
		"{handle.username}": string(handle.Username),
		"{request.scheme}":  string(request.URL.Scheme),
		"{request.host}":    string(request.Host),
		"{request.path}":    string(request.URL.Path),
		"{request.query}":   string(request.URL.RawQuery),
	}

	url := string(template)

	for token, value := range replacements {
		url = strings.ReplaceAll(url, token, value)
	}

	return url
}
