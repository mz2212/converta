package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/floostack/transcoder/ffmpeg"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type converter struct {
	format      string
	overwrite   bool
	codec       string
	ffmpegPath  string
	ffprobePath string
}

func main() {
	viper.SetDefault("inFolder", filepath.Join(os.Getenv("USERPROFILE"), "/AppData/LocalLow/Nolla_Games_Noita/save_rec/screenshots_animated"))
	viper.SetDefault("outFolder", filepath.Join(os.Getenv("USERPROFILE"), "/Videos"))
	viper.SetDefault("remove", true)
	viper.SetDefault("format", "webm")
	viper.SetDefault("codec", "librav1e")
	viper.SetDefault("ffmpegPath", "C:/bin/ffmpeg.exe")
	viper.SetDefault("ffprobePath", "C:/bin/ffprobe.exe")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		if err == err.(viper.ConfigFileNotFoundError) {
			fmt.Println("Couldn't find config, using defaults!")
		} else {
			panic(err)
		}
	}

	c := converter{
		format:      viper.GetString("format"),
		overwrite:   viper.GetBool("remove"),
		codec:       viper.GetString("codec"),
		ffmpegPath:  viper.GetString("ffmpegPath"),
		ffprobePath: viper.GetString("ffprobePath"),
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	err = watcher.Add(viper.GetString("inFolder")) // Tell the watcher to observe the gif folder
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
				webm := strings.Split(filepath.Base(event.Name), ".")[0] + viper.GetString("format") // Construct the video filename
				outFile := filepath.Join(viper.GetString("outFolder"), webm)
				c.convert(event.Name, outFile)
				fmt.Println("Encoded: ", outFile)
				if viper.GetBool("remove") {
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

func (c converter) convert(in string, out string) {
	opts := ffmpeg.Options{
		OutputFormat: &c.format,
		Overwrite:    &c.overwrite,
		VideoCodec:   &c.codec,
	}

	ffmpegConf := ffmpeg.Config{
		FfmpegBinPath:   c.ffmpegPath,
		FfprobeBinPath:  c.ffprobePath,
		ProgressEnabled: false,
	}

	_, err := ffmpeg.New(&ffmpegConf).Input(in).Output(out).WithOptions(opts).Start(opts)
	if err != nil {
		panic(err)
	}
}
