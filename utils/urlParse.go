package utils

import (
	"net/url"
)

func ParseUrl(u string) (host string, ssl bool, username string, err error) {
	//u := "fc+ssl://srv1.local:2283/vrypan"
	p, err := url.Parse(u)
	if err != nil {
		return "", false, "", err
	}

	host = p.Host
	ssl = false
	if p.Scheme == "fc+ssl" {
		ssl = true
	}

	if len(p.Path) > 0 {
		username = p.Path
	}
	return host, ssl, username, nil
}
