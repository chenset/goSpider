package requestor

import (
	"net/http"
	"golang.org/x/net/proxy"
	"fmt"
	"os"
	"net"
	"time"
	"context"
)

func client(ssLocalServer string) {
	httpReq, err := http.NewRequest("GET", "https://ss.flysay.com/", nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't create request:", err)
		os.Exit(2)
	}

	getHttp(httpClients(ssLocalServer), httpReq)
}

func getHttp(httpClient *http.Client, httpReq *http.Request) {
	var startTimeF float64 = float64(time.Now().UnixNano() / 1000)

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't GET page:", err)
	} else {
		defer resp.Body.Close()
	}

	fmt.Println((float64(time.Now().UnixNano()/1000)-startTimeF)/1000.0, "ms")
}

func httpClients(ssLocalServer string) *http.Client {
	dialer, err := proxy.SOCKS5("tcp", ssLocalServer, nil, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		os.Exit(1)
	}

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
