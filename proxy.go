package main

import (
	"github.com/pkg/errors"
	"log"
	"net"
	"net/url"
	"regexp"
)

type Proxy struct {
	matcher     Matcher
	proxyRunner ProxyRunner
}

type Matcher interface {
	Match(host string) bool
}

type RegexpMatcher struct {
	reg *regexp.Regexp
}

func NewRegexpMatcher(str string) (*RegexpMatcher, error) {
	reg, err := regexp.Compile(str)
	if err != nil {
		return nil, err
	}

	return &RegexpMatcher{
		reg: reg,
	}, nil
}

func (r *RegexpMatcher) Match(host string) bool {
	return r.reg.Match([]byte(host))
}

type IPMatcher struct {
	ip string
}

func NewIPMatcher(ip string) *IPMatcher {
	return &IPMatcher{
		ip: ip,
	}
}

func (r *IPMatcher) Match(host string) bool {
	host, port, err := net.SplitHostPort(host)
	if err != nil {
		log.Printf("warning: net.SplitHostPort failed %s", err)
		return false
	}
	return host == r.ip
}

type ProxyHandler struct {
}

func NewProxyHandler(c *Config) (*ProxyHandler, error) {
	proxies := make([]Proxy, 0, len(c.Matches))
	for _, match := range c.Matches {
		var matcher Matcher
		if len(match.Regexp) > 0 {
			matcher, err := NewRegexpMatcher(match.Regexp)
			if err != nil {
				return nil, errors.Wrap(err, "failed to create regexp matcher")
			}
		} else if len(match.IP) > 0 {
			matcher = NewIPMatcher(match.IP)
		} else {
			return nil, errors.Errorf("failed to detect matcher type %v", match)
		}

		var runner ProxyRunner
		if match.Direct {
			runner = NewDirectProxyRunner()
		} else if len(match.Proxy) > 0 {
			runner, err := NewPeerProxyRunner(match.Proxy)
			if err != nil {
				return nil, errors.Wrap(err, "failed to create peer proxy runner")
			}
		} else {
			return nil, errors.Errorf("failed to detect proxy type %v", match)
		}

		proxies = append(proxies, Proxy{
			matcher:     matcher,
			proxyRunner: runner,
		})
	}
	return proxies, nil
}
