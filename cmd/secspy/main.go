package main

/* This is just a test app to demonstrate basic usage of the securityspy library. */

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	securityspy "github.com/davidnewhall/go-securityspy"
	flg "github.com/spf13/pflag"
)

// Version of the app
var Version = "1.0.0"

// Config represents the CLI args + securityspy.Server.
type Config struct {
	UseSSL bool
	User   string
	Pass   string
	URL    string
	Cmd    string
	Arg    string
	Server securityspy.Server
}

func main() {
	config := parseFlags()
	securityspy.Encoder = "/usr/local/bin/ffmpeg"
	switch config.Cmd {

	// Demonstrates event channels. Events always happen in order.
	// Do not block the channel or things stop working.
	case "events", "event", "e":
		config.getServer()
		fmt.Println("Watching Event Stream (specific events, until disconnect)")
		channel := make(chan securityspy.Event)
		config.Server.Events().BindChan(securityspy.EventStreamDisconnect, channel)
		config.Server.Events().BindChan(securityspy.EventStreamConnect, channel)
		config.Server.Events().BindChan(securityspy.EventMotionDetected, channel)
		config.Server.Events().BindChan(securityspy.EventOnline, channel)
		config.Server.Events().BindChan(securityspy.EventOffline, channel)
		go config.Server.Events().Watch(10*time.Second, true)
		for event := range channel {
			config.showEvent(event)
			if event.Event == securityspy.EventStreamDisconnect {
				config.Server.Events().UnbindAll()
				config.Server.Events().Stop()
				fmt.Println("Got disconnect, bailing out.")
				os.Exit(1)
			}
		}

		// Demonstrates event callbacks. Sometimes they fire out of order.
		// They happen in a go routine, so they can be blocking operations.
	case "callbacks", "callback", "call", "l":
		config.getServer()
		fmt.Println("Watching Event Stream (all events, forever)")
		config.Server.Events().BindFunc(securityspy.EventAllEvents, config.showEvent)
		config.Server.Events().Watch(10*time.Second, true)

	case "cameras", "cams", "cam", "c":
		config.printCamData()
	case "video", "vid", "v":
		config.saveVideo()
	case "picture", "pic", "p":
		config.savePicture()
	case "trigger", "t":
		config.triggerMotion()
	case "files", "file", "f":
		config.showFiles()
	case "download", "d":
		config.downloadFile()
	default:
		_, _ = fmt.Fprintln(os.Stderr, "invalid command:", config.Cmd)
		flg.Usage()
		os.Exit(1)
	}
}

// Turn CLI flags into a config struct.
func parseFlags() *Config {
	config := &Config{}
	flg.Usage = func() {
		fmt.Println("Usage: secspy [--user <user>] [--pass <pass>] [--url <url>] [-c <cmd>] [-a <arg>]")
		flg.PrintDefaults()
	}
	flg.StringVarP(&config.User, "user", "u", os.Getenv("SECSPY_USERNAME"), "Username to authenticate with")
	flg.StringVarP(&config.Pass, "pass", "p", os.Getenv("SECSPY_PASSWORD"), "Password to authenticate with")
	flg.StringVarP(&config.URL, "url", "U", "http://127.0.0.1:8000", "SecuritySpy URL")
	flg.BoolVarP(&config.UseSSL, "verify-ssl", "s", false, "Validate SSL certificate if using https")
	flg.StringVarP(&config.Cmd, "command", "c", "", "Command to run. Currently supports: events/callback, cams, pic, vid, trigger, files, download")
	flg.StringVarP(&config.Arg, "arg", "a", "", "if cmd supports an argument, pass it here. ie. -c pic -a Porch:/tmp/filename.jpg")
	version := flg.BoolP("version", "v", false, "Print the version and exit")
	if flg.Parse(); *version {
		fmt.Println("secspy version:", Version)
		os.Exit(0) // don't run anything else.
	}
	return config
}

// getServer makes, saves and returns a securitypy handle.
func (c *Config) getServer() securityspy.Server {
	var err error
	if c.Server, err = securityspy.GetServer(c.User, c.Pass, c.URL, c.UseSSL); err != nil {
		fmt.Println("SecuritySpy Error:", err)
		os.Exit(1)
	}
	fmt.Printf("%v %v (http://%v:%v/) %d cameras, %d scripts, %d sounds\n",
		c.Server.Info().Name, c.Server.Info().Version, c.Server.Info().IP1,
		c.Server.Info().HTTPPort, len(c.Server.GetCameras()),
		len(c.Server.Info().Scripts.Names), len(c.Server.Info().Sounds.Names))
	return c.Server
}

func (c *Config) triggerMotion() {
	if c.Arg == "" {
		fmt.Println("Triggers motion on a camera.")
		fmt.Println("Supply a camera name with -a <cam>[,<cam>][,<cam>]")
		fmt.Println("Example: secspy -c trigger -a Door,Gate")
		fmt.Println("See camera names with -c cams")
		os.Exit(1)
	}
	srv := c.getServer()
	for _, arg := range strings.Split(c.Arg, ",") {
		if cam := srv.GetCameraByName(arg); cam == nil {
			fmt.Println("Camera does not exist:", arg)
			continue
		} else if err := cam.TriggerMotion(); err != nil {
			fmt.Printf("Error Triggering Motion for camera '%v': %v", arg, err)
			continue
		}
		fmt.Println("Triggered Motion for Camera:", arg)
	}
}

// showEvent is a callback function fired by the event watcher in securityspy library.
func (c *Config) showEvent(e securityspy.Event) {
	camString := "No Camera"
	// Always check Camera interface for nil.
	if e.Camera != nil {
		camString = "Camera " + e.Camera.Num() + ": " + e.Camera.Name()
	} else if e.ID < 0 {
		camString = "SecuritySpy Server"
	}
	fmt.Printf("[%v] Event %d: %v, %v, Msg: %v\n",
		e.When, e.ID, e.Event.Event(), camString, e.Msg)
}

// printCamData formats camera data onto a screen for an operator.
func (c *Config) printCamData() {
	for _, camera := range c.getServer().GetCameras() {
		cam := camera.Device()
		fmt.Printf("%2v: %-14v (%-9v %5v/%-7v %v) connected: %3v, down %v, modes: C:%-8v M:%-8v A:%-8v "+
			"%2vFPS, Audio:%3v, MD: %3v/pre:%v/post:%3v idle %-10v Script: %v (reset %v)\n",
			camera.Num(), cam.Name, camera.Size(), cam.DeviceName, cam.DeviceType, cam.Address,
			cam.Connected.Val, cam.TimeSinceLastFrame.Dur.String(), cam.ModeC.Txt, cam.ModeM.Txt,
			cam.ModeA.Txt+",", int(cam.CurrentFPS), cam.HasAudio.Txt, cam.MDenabled.Txt,
			cam.MDpreCapture.Dur.String(), cam.MDpostCapture.Dur.String(),
			cam.TimeSinceLastMotion.Dur.String(), cam.ActionScriptName, cam.ActionResetTime.Dur.String())
	}
}

func (c *Config) savePicture() {
	if c.Arg == "" || !strings.Contains(c.Arg, ":") {
		fmt.Println("Saves a single still JPEG image from a camera.")
		fmt.Println("Supply a camera name and file path with -a <cam>:<path>")
		fmt.Println("Example: secspy -c pic -a Porch:/tmp/Porch.jpg")
		fmt.Println("See camera names with -c cams")
		os.Exit(1)
	}
	split := strings.Split(c.Arg, ":")
	cam := c.getServer().GetCameraByName(split[0])
	if cam == nil {
		fmt.Println("Camera does not exist:", split[0])
		os.Exit(1)
	} else if err := cam.SaveJPEG(&securityspy.VidOps{}, split[1]); err != nil {
		fmt.Printf("Error Saving Image for camera '%v' to file '%v': %v\n", cam.Name(), split[1], err)
		os.Exit(1)
	}
	fmt.Printf("Image for camera '%v' saved to: %v\n", cam.Name(), split[1])
}

func (c *Config) saveVideo() {
	if c.Arg == "" || !strings.Contains(c.Arg, ":") {
		fmt.Println("Saves a 10 second video from a camera.")
		fmt.Println("Supply a camera name and file path with -a <cam>:<path>")
		fmt.Println("Example: secspy -c pic -a Gate:/tmp/Gate.mov")
		fmt.Println("See camera names with -c cams")
		os.Exit(1)
	}
	split := strings.Split(c.Arg, ":")
	cam := c.getServer().GetCameraByName(split[0])
	if cam == nil {
		fmt.Println("Camera does not exist:", split[0])
		os.Exit(1)
	} else if err := cam.SaveVideo(&securityspy.VidOps{}, 10*time.Second, 9999999999, split[1]); err != nil {
		fmt.Printf("Error Saving Video for camera '%v' to file '%v': %v\n", cam.Name(), split[1], err)
		os.Exit(1)
	}
	fmt.Printf("10 Second video for camera '%v' saved to: %v\n", cam.Name(), split[1])
}

func (c *Config) showFiles() {
	if c.Arg == "" {
		fmt.Println("Shows last files captured by securityspy")
		fmt.Println("Supply camera names and file age with -a <cam>,<cam>:<days old>")
		fmt.Println("Example: secspy -c files -a Porch,Gate:10")
		fmt.Println("See camera names with -c cams")
		os.Exit(1)
	}
	split := strings.Split(c.Arg, ":")
	daysOld := 14
	if len(split) > 1 {
		daysOld, _ = strconv.Atoi(split[1])
		if daysOld < 1 {
			daysOld = 14
		}
	}
	srv := c.getServer()
	var cameraNums []int
	// Loop the provided camera names and find their numbers.
	for _, name := range strings.Split(split[0], ",") {
		cam := srv.GetCameraByName(name)
		if cam == nil {
			fmt.Println("Camera does not exist:", name)
			continue
		}
		cameraNums = append(cameraNums, cam.Number())
	}
	age := time.Now().Add(-time.Duration(daysOld*24) * time.Hour)
	files, err := srv.Files().GetAll(cameraNums, age, time.Now())
	if err != nil {
		fmt.Println("Received error from Files.GetAll() method:", err)
	}
	fmt.Printf("Found %d files. From %v to %v:\n", len(files), age.Format("01/02/2006"), time.Now().Format("01/02/2006"))
	for _, file := range files {
		camName := "<no camera>"
		if file.Camera() != nil {
			camName = file.Camera().Name()
		}
		fmt.Printf("[%v] %v %v: '%v' (%vMB)\n",
			file.Date(), camName, file.Type(), file.Name(), file.Size()/1024/1024)
	}
}

func (c *Config) downloadFile() {
	if c.Arg == "" || !strings.Contains(c.Arg, ":") {
		fmt.Println("Downloads a saved media file from SecuritySpy.")
		fmt.Println("Supply file name and save-path with -a 'filename:path'")
		fmt.Println("Example: secspy -c download -a '01-19-2019 00-01-23 M Porch.m4v:/tmp/file.m4v'")
		fmt.Println("See file names with -c files")
		os.Exit(1)
	}

	srv := c.getServer()
	fileName := strings.Split(c.Arg, ":")[0]
	savePath := strings.Split(c.Arg, ":")[1]
	if _, err := os.Stat(savePath); !os.IsNotExist(err) {
		fmt.Println("File already exists:", savePath)
		os.Exit(1)
	}
	file, err := srv.Files().GetFile(fileName)
	if err != nil {
		fmt.Println("Error getting file:", err)
		os.Exit(1)
	}
	size, err := file.Save(savePath)
	if err != nil {
		fmt.Println("Error writing file:", err)
		os.Exit(1)
	}
	fmt.Println("File saved to:", savePath, "->", size/1024/1024, "MB")
}
