package web

import (
	"net/http"
	"net/url"
)

type HttpHandle struct {
	method string
	url    *url.URL
}

func (self HttpHandle) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	self.method = request.Method
	self.url = request.URL

	response.Write([]byte(self.get()))
}

func (self *HttpHandle) get() string {
	return "response"
}

func Listen() {
	//http.ListenAndServeTLS(":443", "./10.0.0.110.crt", "./10.0.0.110.key", myHandle{})
	http.ListenAndServe(":80", HttpHandle{})
}
