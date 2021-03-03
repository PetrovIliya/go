package urlresolver

import (
	"fmt"
	"net/http"
)

type HttpHandler struct {
	ShortUrlPaths ShortUrlsPaths
}

type ShortUrlsPaths struct {
	Paths map[string]string `json:"paths"`
}

func (handler HttpHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	requestUri := request.RequestURI
	if requestUri == "/favicon.ico" {
		return
	}

	if len(handler.ShortUrlPaths.Paths) == 0 {
		fmt.Println("No short paths")
		return
	}
	realUrl, shortUrlExist := handler.ShortUrlPaths.Paths[requestUri]
	if shortUrlExist {
		http.Redirect(response, request, realUrl, 301)
	} else if requestUri != "/" {
		response.WriteHeader(http.StatusNotFound)
		_, _ = response.Write([]byte("<h1>Unknown short url " + requestUri + " </h1>"))
	}
}
