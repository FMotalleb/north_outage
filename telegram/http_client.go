package telegram

import (
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

func httpClient(proxyUrl *url.URL) *http.Client {
	dialer := proxy.FromEnvironmentUsing(proxy.Direct)

	if proxyUrl != nil {
		if d, err := proxy.FromURL(proxyUrl, dialer); err != nil {
			panic("failed to use given url as proxy")
		} else {
			dialer = d
		}
	}

	transport := &http.Transport{
		Dial: dialer.Dial,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
	}
	return client
}
