package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path"
)

const StopCharacter string = "\r\n"

var HomeDir string = os.Getenv("HOME")
var DefaultConfigFilePath = path.Join(HomeDir, ".splash.json")
var DefaultImageDirPath = path.Join(HomeDir, "Pictures")

var (
	// APP
	app = kingpin.New("splash", "Splash your screen with beautiful wallpapers.")

	verbose = app.Flag("verbose", "Enable verbose mode.").
		Short('v').
		Bool()
	port = app.Flag("port", "Port to connect.").
		Short('p').
		Default("8086").
		Int()

	// APP SERVER
	server = app.Command("server", "Open splash server.")

	configFile = server.Flag("config", "Configuration file").
			Short('c').
			PlaceHolder("path").
			Default(DefaultConfigFilePath).
			ExistingFile()
	imageDir = server.Flag("image", "Image directory").
			Short('i').
			Default(DefaultImageDirPath).
			PlaceHolder("path").
			ExistingDir()
    slideshow = server.Flag("slideshow", "Change wallpaper automatically.").
            Short('S').
            PlaceHolder("duration").
            Duration()
    setOnStartup = server.Flag("set-on-startup", "Set random wallpaper at startup").
            Short('u').
            Bool()

	// APP CLIENT
	client = app.Command("client", "Open splash client.")

	ip = client.Flag("server", "Server address.").
		Short('s').
		Default("127.0.0.1").
		IP()
	setRandomWallpaper = client.Flag("change", "Set random wallpaper").
				Short('r').
				Bool()
	fetchWallpaper = client.Flag("fetch", "Fetch new wallpapers").
			Short('f').
			Bool()
)

func main() {
	var exitCode int

	kingpin.Version("0.0.1")
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case client.FullCommand():
		var message string
		if *setRandomWallpaper {
			message = "set-random-wallpaper"
		} else if *fetchWallpaper {
			message = "fetch-wallpaper"
		}
		exitCode = RunClient(*verbose, ip.String(), *port, message)
	case server.FullCommand():
		exitCode = RunServer(*verbose, *port, *configFile, *imageDir, *slideshow, *setOnStartup)
	}
	os.Exit(exitCode)
}
