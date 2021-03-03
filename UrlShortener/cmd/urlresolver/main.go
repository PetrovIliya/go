package main

import (
	"UrlShortener/pkg/urlresolver"
	"UrlShortener/pkg/util"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	args := os.Args
	var fileName string
	if len(args) < 2 {
		fmt.Println("Required argument 'fileName' missed")
		return
	}
	fileName = args[1]
	fileContentJson := util.ReadFile(fileName)
	var shortUrlPaths urlresolver.ShortUrlsPaths
	_ = json.Unmarshal([]byte(fileContentJson), &shortUrlPaths)

	handler := urlresolver.HttpHandler{ShortUrlPaths: shortUrlPaths}
	_ = http.ListenAndServe(":9000", handler)
}
