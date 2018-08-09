package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"path"
	"strconv"
	"strings"
	"time"
)

func isTransportOver(data string) (over bool) {
	over = strings.HasSuffix(data, StopCharacter)
	return
}

func handler(socket net.Conn, config Config, imageDir string) {
	defer socket.Close()

	buf := make([]byte, 1024)
	r := bufio.NewReader(socket)
	w := bufio.NewWriter(socket)

	// Read from client
	var completeData string
	for {
		n, err := r.Read(buf)
		data := string(buf[:n])

		if err == nil {
			log.Printf("Receive: %s", data)
			completeData += data
			if isTransportOver(data) {
				break
			}
		} else if err == io.EOF {
			break
		} else {
			log.Fatalf("Receive data failed: %s", err)
			return
		}
	}
	completeData = strings.Trim(completeData, StopCharacter)

	var result string
	switch completeData {
	case "set-random-wallpaper":
		wallpaperErr := SetRandomWallpaper(imageDir)
		if wallpaperErr != nil {
			result = wallpaperErr.Error()
			break
		}

		result = "Done"
		break
	case "fetch-wallpaper":
		imagePath := path.Join(imageDir, time.Now().Format("2006-01-02 15:04:05")+".jpg")
		requestUrl := config.getUrl()

		saveErr := SaveImage(requestUrl, imagePath)
		if saveErr != nil {
			result = saveErr.Error()
			break
		}

		wallpaperErr := SetWallpaper(imagePath)
		if wallpaperErr != nil {
			result = wallpaperErr.Error()
			break
		}

		result = "Done"
		break
	default:
		result = "Unexpected command"
	}

	// Reply to client
	w.Write([]byte(result))
	w.Flush()
	log.Printf("Sent: %s", result)
}

func slideshowHandler(ticker *time.Ticker, quitTicker chan struct{}, imageDir string) {
    for {
        select {
            case <- ticker.C:
                wallpaperErr := SetRandomWallpaper(imageDir)
                if wallpaperErr != nil {
                    log.Printf("Error: %s", wallpaperErr.Error())
                }
            case <- quitTicker:
                ticker.Stop();
                return;
        }
    }
}

func RunServer(verbose bool, port int, configFile string, imageDir string, slideshow time.Duration, setOnStartup bool) int {
	if verbose {
		log.Printf("Config file: %s", configFile)
		log.Printf("Image directory: %s", imageDir)
        if slideshow.Seconds() > 0 {
            log.Printf("Slideshow is enabled on every %f", slideshow.Seconds())
        } else {
            log.Printf("Slideshow is disabled")
        }
	}

    if setOnStartup {
        wallpaperErr := SetRandomWallpaper(imageDir)
        if wallpaperErr != nil {
            log.Printf("Error: %s", wallpaperErr.Error())
        }
    }

    // set wallpaper at given interval
    if (slideshow.Seconds() > 0) {
        ticker := time.NewTicker(slideshow)
        defer ticker.Stop()

        quitTicker := make(chan struct{})

        // FIXME: restar ticker on any other action by user
        go slideshowHandler(ticker, quitTicker, imageDir);
    }


	config, configErr := LoadConfig(configFile)
	if configErr != nil {
		log.Printf("Error: %s", configErr.Error())
		return 1
	}

	// start listening
	listen, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Socket listen on port %d failed.\n%s", port, err)
		return 1
	}
	defer listen.Close()
	log.Printf("Listening on port: %d", port)

	for {
		// wait for a connection
		socket, err := listen.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		// open up a handler
		go handler(socket, config, imageDir)
	}
}
