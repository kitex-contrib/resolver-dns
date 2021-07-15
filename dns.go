// Package dns resolver
package dns

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

const (
	defaultPort   = "443"
	defaultWeight = 10
)

// dnsResolver is a resolver using consul with asynccache.
type dnsResolver struct{}

// NewDNSResolver create a dns based resolver
func NewDNSResolver() discovery.Resolver {
	return new(dnsResolver)
}

var _ discovery.Resolver = (*dnsResolver)(nil)

// Target implements the Resolver interface.
func (r *dnsResolver) Target(ctx context.Context, target rpcinfo.EndpointInfo) string {
	return target.ServiceName()
}

func (r *dnsResolver) lookupHost(host, port string) ([]discovery.Instance, error) {
	addrs, err := net.LookupHost(host)
	if err != nil {
		return nil, err
	}
	ins := make([]discovery.Instance, 0, len(addrs))
	for _, a := range addrs {
		ins = append(ins, discovery.NewInstance("tcp", net.JoinHostPort(a, port), defaultWeight, nil))
	}
	return ins, nil
}

// Resolve implements the Resolver interface.
func (r *dnsResolver) Resolve(ctx context.Context, target string) (discovery.Result, error) {
	host, port, err := parseTarget(target, defaultPort)
	if err != nil {
		return discovery.Result{}, err
	}

	eps, err := r.lookupHost(host, port)
	if err != nil {
		return discovery.Result{}, fmt.Errorf("failed to resolve '%s': %w", target, err)
	}

	// Empty slice from dns server should be treated as a normal result.
	if len(eps) == 0 {
		return discovery.Result{}, fmt.Errorf("no instance remains for %s", target)
	}
	res := discovery.Result{
		Cacheable: true,
		CacheKey:  target,
		Instances: eps,
	}
	return res, nil
}

func parseTarget(target, defaultPort string) (host, port string, err error) {
	if target == "" {
		return "", "", errors.New("missing target address")
	}
	if ip := net.ParseIP(target); ip != nil {
		// target is an IPv4 or IPv6(without brackets) address
		return target, defaultPort, nil
	}
	if host, port, err = net.SplitHostPort(target); err == nil {
		if port == "" {
			// If the port field is empty (target ends with colon), e.g. "[::1]:", this is an error.
			return "", "", errors.New("missing port after port-separator colon")
		}
		// target has port, i.e ipv4-host:port, [ipv6-host]:port, host-name:port
		if host == "" {
			// Keep consistent with net.Dial(): If the host is empty, as in ":80", the local system is assumed.
			host = "localhost"
		}
		return host, port, nil
	}
	if host, port, err = net.SplitHostPort(target + ":" + defaultPort); err == nil {
		// target doesn't have port
		return host, port, nil
	}
	return "", "", fmt.Errorf("invalid target address %v, error info: %v", target, err)
}

// Diff implements the Resolver interface.
func (r *dnsResolver) Diff(key string, prev, next discovery.Result) (discovery.Change, bool) {
	return discovery.DefaultDiff(key, prev, next)
}

// Name implements the Resolver interface.
func (r *dnsResolver) Name() string {
	return "dns"
}
