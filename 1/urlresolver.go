package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type HttpHandler struct {
	shortUrlPaths ShortUrlsPaths
}

type ShortUrlsPaths struct {
	Paths map[string]string `json: "paths"`
}

func (handler HttpHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	requestUri := request.RequestURI
	if requestUri == "/favicon.ico" {
		return
	}

	if len(handler.shortUrlPaths.Paths) == 0 {
		fmt.Println("No short paths")
		return
	}
	realUrl, shortUrlExist := handler.shortUrlPaths.Paths[requestUri]
	if shortUrlExist {
		http.Redirect(response, request, realUrl, 301)
	} else if requestUri != "/" {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte("<h1>Unknown short url " + requestUri + " </h1>"))
	}
}

func readFile(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil{
		fmt.Println(err)
		return ""
	}
	data := make([]byte, 64)

	result := ""
	for{
		n, err := file.Read(data)
		if err == io.EOF{
			break
		}
		result += string(data[:n])
	}

	_ = file.Close()
	return result
}

func main() {
	args := os.Args
	var fileName string
	if len(args) < 2 {
		fmt.Println("Required argument 'fileName' missed")
		return
	}
	fileName = args[1]
	fileContentJson := readFile(fileName)
	var shortUrlPaths ShortUrlsPaths
	json.Unmarshal([]byte(fileContentJson), &shortUrlPaths)

	handler := HttpHandler{shortUrlPaths}
	http.ListenAndServe(":9000", handler)
}