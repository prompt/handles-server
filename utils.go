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
		Hostname: Hostname(parts[0]),
		Domain:   Domain(parts[2]),
		Username: Username(parts[1]),
	}, nil
}

func FormatTemplateUrl(
	template string,
	request *http.Request,
	handle Handle,
	did DecentralizedID,
) string {
	replacements := map[string]string{
		"{handle}":          handle.String(),
		"{did}":             string(did),
		"{handle.hostname}": string(handle.Hostname),
		"{handle.domain}":   string(handle.Domain),
		"{handle.username}": string(handle.Username),
		"{request.scheme}":  string(request.URL.Scheme),
		"{request.host}":    string(request.Host),
		"{request.path}":    string(request.URL.Path),
		"{request.query}":   string(request.URL.RawQuery),
	}

	for token, value := range replacements {
		template = strings.ReplaceAll(template, token, value)
	}

	return template
}
