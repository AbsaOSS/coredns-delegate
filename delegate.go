package delegation

import (
	"github.com/miekg/dns"
	"math/rand"
	"strings"
)

const dnsPort = ":53"

type DelegationResponseWriter struct {
	dns.ResponseWriter
}

func (dw *DelegationResponseWriter) WriteMsg(m *dns.Msg) error {
	if m.Authoritative {
		return dw.ResponseWriter.WriteMsg(m)
	}

	resolver := new(dns.Client)
	msg := new(dns.Msg)

	ns := extractNs(m.Ns)

	msg.SetQuestion(m.Question[0].Name, m.Question[0].Qtype)
	in, _, err := resolver.Exchange(msg, ns[rand.Intn(len(ns))])
	if err != nil {
		log.Info("err: ", err.Error())
	}
	m.Answer = in.Answer

	return dw.ResponseWriter.WriteMsg(m)
}

func (dw *DelegationResponseWriter) Write(buf []byte) (int, error) {
	n, err := dw.ResponseWriter.Write(buf)
	return n, err
}

func extractNs(rrset []dns.RR) []string {
	var servers []string
	for _, rr := range rrset {
		ns := dns.Field(rr, 1)
		nameserver := strings.TrimSuffix(ns, ".") + dnsPort
		servers = append(servers, nameserver)
	}
	return servers
}
