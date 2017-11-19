package main

import (
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

type ProxyRunner interface {
	Run(w http.ResponseWriter, req *http.Request)
}

type DirectProxyRunner struct {
}

func (r *DirectProxyRunner) Run(w http.ResponseWriter, req *http.Request) {
}

func NewDirectProxyRunner() ProxyRunner {
	return &DirectProxyRunner{}
}

type PeerProxyRunner struct {
	u *url.URL
}

func NewPeerProxyRunner(urlStr string) (ProxyRunner, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse url %s", urlStr)
	}
	return &PeerProxyRunner{
		u: u,
	}, nil
}

func (r *PeerProxyRunner) Run(w http.ResponseWriter, req *http.Request) {
}
