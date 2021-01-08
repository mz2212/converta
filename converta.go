package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/floostack/transcoder/ffmpeg"
	"github.com/fsnotify/fsnotify"
)

func main() {
	inFolder := filepath.Join(os.Getenv("USERPROFILE"), "/AppData/LocalLow/Nolla_Games_Noita/save_rec/screenshots_animated")
	outFolder := filepath.Join(os.Getenv("USERPROFILE"), "/Videos") // Construct the paths for the Videos folder and Noita Gif folder
	remove := true                                                  // Set to false if you don't want to delete the gifs after conversion

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	err = watcher.Add(inFolder) // Tell the watcher to observe the gif folder
	if err != nil {
		panic(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write { // Check if the event contains a write event. Note that that's a binary "and"
				fmt.Println("Encoding: ", event.Name)
				webm := strings.Split(filepath.Base(event.Name), ".")[0] + ".webm" // Construct the video filename
				outFile := filepath.Join(outFolder, webm)
				convert(event.Name, outFile)
				fmt.Println(outFile)
				if remove {
					os.Remove(event.Name)
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("Error: ", err)
		}
	}

}

func convert(in string, out string) {
	format := "webm" // Twiddle these to change the ffmpeg settings. I'll probably do something different in the end.
	overwrite := true
	codec := "librav1e"

	opts := ffmpeg.Options{
		OutputFormat: &format,
		Overwrite:    &overwrite,
		VideoCodec:   &codec,
	}

	ffmpegConf := ffmpeg.Config{ // :thinking: Not quite sure what to do about these. The library _needs_ to know where they are, and it might be different on between PCs.
		FfmpegBinPath:   "C:/bin/ffmpeg.exe",
		FfprobeBinPath:  "C:/bin/ffprobe.exe",
		ProgressEnabled: false,
	}

	_, err := ffmpeg.New(&ffmpegConf).Input(in).Output(out).WithOptions(opts).Start(opts)
	if err != nil {
		panic(err)
	}
}
