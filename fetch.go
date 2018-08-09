package main

import (
	"io"
	"net/http"
	"os"
)

func SaveImage(requestUrl string, filePath string) error {
	// Pull background and save it
	response, requestErr := http.Get(requestUrl)
	if requestErr != nil {
		return requestErr
		// fmt.Println("ERROR:", requestErr.Error())
	}
	defer response.Body.Close()
	// NOTE: check for other errors

	imageFile, fileErr := os.Create(filePath)
	if fileErr != nil {
		// fmt.Println(fileErr.Error())
		return fileErr
		// return 1
	}
	defer imageFile.Close()

	// save response body into a file
	io.Copy(imageFile, response.Body)
	return nil
}
