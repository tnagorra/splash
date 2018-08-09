package main

import (
	"io/ioutil"
	"math/rand"
	"os/exec"
	"path"
	"time"
)

func SetWallpaper(filePath string) error {
	return exec.Command("feh", "--bg-fill", filePath).Run()
}

func SetRandomWallpaper(imageDirPath string) error {
	rand.Seed(time.Now().Unix())

	fileInfo, err := ioutil.ReadDir(imageDirPath)
	if err != nil {
		return err
	}

	filteredFileInfo := fileInfo[:0]
	for _, x := range fileInfo {
		if !x.IsDir() {
			filteredFileInfo = append(filteredFileInfo, x)
		}
	}

	totalItemsCount := len(filteredFileInfo)

	if totalItemsCount <= 0 {
		// NOTE: No need to do anything
		return nil
	}

	randInt := rand.Intn(totalItemsCount)
	fileName := filteredFileInfo[randInt].Name()
	absolutePath := path.Join(imageDirPath, fileName)

	return exec.Command("feh", "--bg-fill", absolutePath).Run()
}
