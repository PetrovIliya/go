package util

import (
	"fmt"
	"io"
	"os"
)

func ReadFile(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	data := make([]byte, 64)

	result := ""
	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		result += string(data[:n])
	}

	_ = file.Close()
	return result
}
