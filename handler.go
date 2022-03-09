package delegation

import (
	"context"

	"github.com/coredns/coredns/plugin"

	"github.com/miekg/dns"
)

type Delegation struct {
	Next plugin.Handler
}

// ServeDNS implements the plugin.Handler interface.
func (d Delegation) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	drr := &DelegationResponseWriter{w}
	return plugin.NextOrFailure(d.Name(), d.Next, ctx, drr, r)
}

// Name implements the Handler interface.
func (d Delegation) Name() string { return "delegation" }
