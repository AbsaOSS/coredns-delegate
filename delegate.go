package delegation

import (
	"github.com/miekg/dns"
	"strings"
)

const dnsPort = ":53"

type DelegationResponseWriter struct {
	dns.ResponseWriter
}

func (dw *DelegationResponseWriter) WriteMsg(m *dns.Msg) error {
	resolver := new(dns.Client)
	msg := new(dns.Msg)

	if m.Authoritative {
		return dw.ResponseWriter.WriteMsg(m)
	}

	ns := strings.Split(m.Ns[0].String(), "\t")
	nameserver := strings.TrimSuffix(ns[4], ".") + dnsPort
	msg.SetQuestion(m.Question[0].Name, m.Question[0].Qtype)
	in, _, err := resolver.Exchange(msg, nameserver)
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
