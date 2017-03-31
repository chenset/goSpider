package requestor

import (
	"net/http"
	"fmt"
	"os"
	"time"
)

func request(client *http.Client) {
	httpReq, err := http.NewRequest("GET", "https://ss.flysay.com/", nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't create request:", err)
		os.Exit(2)
	}

	getHttp(client, httpReq)
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
