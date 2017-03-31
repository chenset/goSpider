package requestProxy

import (
	"golang.org/x/net/proxy"
)

type Socks5 struct {
	socks5Server string
}

func NewSocks5(socks5Server string) *Socks5 {
	return &Socks5{socks5Server: socks5Server}
}

func (cls *Socks5) StaticDialer() (proxy.Dialer, error) {
	return proxy.SOCKS5("tcp", cls.socks5Server, nil, proxy.Direct)
}

func (cls *Socks5) Dialer() (proxy.Dialer, error) {
	return proxy.SOCKS5("tcp", cls.socks5Server, nil, proxy.Direct)
}
