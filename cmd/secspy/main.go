package main

/* This is just a test app to demonstrate basic usage of the securityspy library. */

import (
	"errors"
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
	Server *securityspy.Server
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
		config.Server.Events.BindChan(securityspy.EventStreamDisconnect, channel)
		config.Server.Events.BindChan(securityspy.EventStreamConnect, channel)
		config.Server.Events.BindChan(securityspy.EventMotionDetected, channel)
		config.Server.Events.BindChan(securityspy.EventOnline, channel)
		config.Server.Events.BindChan(securityspy.EventOffline, channel)
		go config.Server.Events.Watch(10*time.Second, true)
		for event := range channel {
			config.showEvent(event)
			if event.Event == securityspy.EventStreamDisconnect {
				config.Server.Events.UnbindAll()
				config.Server.Events.Stop()
				fmt.Println("Got disconnect, bailing out.")
				os.Exit(1)
			}
		}

		// Demonstrates event callbacks. Sometimes they fire out of order.
		// They happen in a go routine, so they can be blocking operations.
	case "callbacks", "callback", "call", "l":
		config.getServer()
		fmt.Println("Watching Event Stream (all events, forever)")
		config.Server.Events.BindFunc(securityspy.EventAllEvents, config.showEvent)
		config.Server.Events.Watch(10*time.Second, true)

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
	case "ptz", "z":
		config.controlPTZ()
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
	flg.StringVarP(&config.Cmd, "command", "c", "", "Command to run. Currently supports: events/callback, cams, pic, vid, trigger, files, download, ptz")
	flg.StringVarP(&config.Arg, "arg", "a", "", "if cmd supports an argument, pass it here. ie. -c pic -a Porch:/tmp/filename.jpg")
	version := flg.BoolP("version", "v", false, "Print the version and exit")
	if flg.Parse(); *version {
		fmt.Println("secspy version:", Version)
		os.Exit(0) // don't run anything else.
	}
	return config
}

// getServer makes, saves and returns a securitypy handle.
func (c *Config) getServer() *securityspy.Server {
	var err error
	if c.Server, err = securityspy.GetServer(c.User, c.Pass, c.URL, c.UseSSL); err != nil {
		fmt.Println("SecuritySpy Error:", err)
		os.Exit(1)
	}
	fmt.Printf("%v %v (http://%v:%v/) %d cameras, %d scripts, %d sounds\n",
		c.Server.Info.Name, c.Server.Info.Version, c.Server.Info.IP1,
		c.Server.Info.HTTPPort, c.Server.Cameras.Count,
		len(c.Server.Info.Scripts.Names), len(c.Server.Info.Sounds.Names))
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
		if cam := srv.Cameras.ByName(arg); cam == nil {
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
		camString = "Camera " + strconv.Itoa(e.Camera.Number) + ": " + e.Camera.Name
	} else if e.ID < 0 {
		camString = "SecuritySpy Server"
	}
	fmt.Printf("[%v] Event %d: %v, %v, Msg: %v\n",
		e.When, e.ID, e.Event.Event(), camString, e.Msg)
}

// printCamData formats camera data onto a screen for an operator.
func (c *Config) printCamData() {
	for _, camera := range c.getServer().Cameras.All() {
		fmt.Printf("%2v: %-14v (%-4vx%-4v %5v/%-7v %v) connected: %3v, down %v, modes: C:%-8v M:%-8v A:%-8v "+
			"%2vFPS, Audio:%3v, MD: %3v/pre:%v/post:%3v idle %-10v Script: %v (reset %v)\n",
			camera.Number, camera.Name, camera.Width, camera.Height, camera.DeviceName, camera.DeviceType, camera.Address,
			camera.Connected.Val, camera.TimeSinceLastFrame.Dur.String(), camera.ModeC.Txt, camera.ModeM.Txt,
			camera.ModeA.Txt+",", int(camera.CurrentFPS), camera.HasAudio.Txt, camera.MDenabled.Txt,
			camera.MDpreCapture.Dur.String(), camera.MDpostCapture.Dur.String(),
			camera.TimeSinceLastMotion.Dur.String(), camera.ActionScriptName, camera.ActionResetTime.Dur.String())
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
	cam := c.getServer().Cameras.ByName(split[0])
	if cam == nil {
		fmt.Println("Camera does not exist:", split[0])
		os.Exit(1)
	} else if err := cam.SaveJPEG(&securityspy.VidOps{}, split[1]); err != nil {
		fmt.Printf("Error Saving Image for camera '%v' to file '%v': %v\n", cam.Name, split[1], err)
		os.Exit(1)
	}
	fmt.Printf("Image for camera '%v' saved to: %v\n", cam.Name, split[1])
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
	cam := c.getServer().Cameras.ByName(split[0])
	if cam == nil {
		fmt.Println("Camera does not exist:", split[0])
		os.Exit(1)
	} else if err := cam.SaveVideo(&securityspy.VidOps{}, 10*time.Second, 9999999999, split[1]); err != nil {
		fmt.Printf("Error Saving Video for camera '%v' to file '%v': %v\n", cam.Name, split[1], err)
		os.Exit(1)
	}
	fmt.Printf("10 Second video for camera '%v' saved to: %v\n", cam.Name, split[1])
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
		cam := srv.Cameras.ByName(name)
		if cam == nil {
			fmt.Println("Camera does not exist:", name)
			continue
		}
		cameraNums = append(cameraNums, cam.Number)
	}
	age := time.Now().Add(-time.Duration(daysOld*24) * time.Hour)
	files, err := srv.Files.GetAll(cameraNums, age, time.Now())
	if err != nil {
		fmt.Println("Received error from Files.All() method:", err)
	}
	fmt.Printf("Found %d files. From %v to %v:\n", len(files), age.Format("01/02/2006"), time.Now().Format("01/02/2006"))
	for _, file := range files {
		camName := "<no camera>"
		if file.Camera != nil {
			camName = file.Camera.Name
		}
		fmt.Printf("[%v] %v %v: '%v' (%vMB)\n",
			file.Updated, camName, file.Link.Type, file.Title, file.Link.Length/1024/1024)
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
	file, err := srv.Files.GetFile(fileName)
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

func (c *Config) controlPTZ() {
	if c.Arg == "" || !strings.Contains(c.Arg, ":") {
		fmt.Println("Controls Camera PTZ.")
		fmt.Println("Supply the Camera and action with -a 'Camera:action'")
		fmt.Println("Example: secspy -c z -a 'Door Cam:Home'")
		fmt.Println("Actions: Home, Up, Down, Left, Right, In, Out, Preset1 .. Preset8")
		os.Exit(1)
	}
	srv := c.getServer()
	splitStr := strings.Split(c.Arg, ":")
	command := strings.ToLower(splitStr[1])
	camera := srv.Cameras.ByName(splitStr[0])
	if camera == nil {
		fmt.Println("camera not found:", splitStr[0])
		os.Exit(1)
	}
	var err error
	switch command {
	case "home":
		err = camera.PTZ.Home()
	case "up":
		err = camera.PTZ.Up()
	case "down":
		err = camera.PTZ.Down()
	case "left":
		err = camera.PTZ.Left()
	case "right":
		err = camera.PTZ.Right()
	case "in":
		err = camera.PTZ.Zoom(true)
	case "out":
		err = camera.PTZ.Zoom(false)
	case "preset1":
		err = camera.PTZ.Preset(securityspy.PTZpreset1)
	case "preset2":
		err = camera.PTZ.Preset(securityspy.PTZpreset2)
	case "preset3":
		err = camera.PTZ.Preset(securityspy.PTZpreset3)
	case "preset4":
		err = camera.PTZ.Preset(securityspy.PTZpreset4)
	case "preset5":
		err = camera.PTZ.Preset(securityspy.PTZpreset5)
	case "preset6":
		err = camera.PTZ.Preset(securityspy.PTZpreset6)
	case "preset7":
		err = camera.PTZ.Preset(securityspy.PTZpreset7)
	case "preset8":
		err = camera.PTZ.Preset(securityspy.PTZpreset8)
	default:
		err = errors.New("invalid command: " + command)
	}
	if err != nil {
		fmt.Println("ptz error:", err)
		os.Exit(1)
	}
	fmt.Println(command, "command sent to", camera.Name)
}
