package requestClient

import (
	"net/http"
	"net"
	"time"
	"golang.org/x/net/proxy"
	"context"
)

func Client(dialer proxy.Dialer) *http.Client {
	netTransport := &http.Transport{
		// Go version < 1.6
		//Dial:dialer.Dial,

		// Go version > 1.6
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		},
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &http.Client{Transport: netTransport}
}
