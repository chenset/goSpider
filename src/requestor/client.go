package requestor

import (
	"net/http"
	"fmt"
	"os"
	"net"
	"time"
	"context"
	"golang.org/x/net/proxy"
)



func client(dialer proxy.Dialer) {
	httpReq, err := http.NewRequest("GET", "https://ss.flysay.com/", nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't create request:", err)
		os.Exit(2)
	}

	getHttp(httpClients(dialer), httpReq)
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

func httpClients(dialer proxy.Dialer) *http.Client {

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
