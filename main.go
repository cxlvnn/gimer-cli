package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	// "strconv"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

const DEF_VALUE_FOR_ALARM = 1
const DEF_PATH_FOR_ALARM_SOUND = "its-time-to-stop.mp3"

var help = flag.Bool("help", false, "show help information")

func main() {
	time_amount_ptr := flag.Int("time", DEF_VALUE_FOR_ALARM, "time - determine how long the timer runs for")
	path_to_an_alarm := flag.String("path", DEF_PATH_FOR_ALARM_SOUND, "path - pick an alarm sound to play after the timer ends")
	flag.Parse()

	if *help {
		fmt.Println("Usage of gimer-cli:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	path := expand_path(*path_to_an_alarm)
	check_path(path)
	alarm_path := filepath.Clean(*path_to_an_alarm)

	if *time_amount_ptr < 1 {
		fmt.Println("You entered a value that is less than one, using the default value: ", DEF_VALUE_FOR_ALARM)
		*time_amount_ptr = DEF_VALUE_FOR_ALARM
	}

	fmt.Println()
	set_timer(*time_amount_ptr)
	play_alarm(alarm_path)
	fmt.Println()

}

func set_timer(minutes int) {
	seconds := minutes*60 - 1
	for remaining_seconds := seconds; remaining_seconds >= 0; remaining_seconds-- {
		minute := remaining_seconds / 60
		seconds := remaining_seconds % 60
		fmt.Printf("\r  %02d:%02d  ", minute, seconds)
		time.Sleep(time.Second)
	}

	fmt.Println("The time is up!")

}

func expand_path(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err.Error())
			strings.Replace(path, "~", home, 1)
		}
	}
	return path
}

func check_path(path string) {
	_, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("The alarm path is likely to be wrong, the file does not exist")
		} else {
			fmt.Println("An unexpected error occured while opening file in that path")
			fmt.Println(err.Error())
		}
	}
}

func play_alarm(path string) {
	alarm_sound, err := os.Open(path)

	decoded, err := mp3.NewDecoder(alarm_sound)
	if err != nil {
		fmt.Println(err.Error())
	}

	context, ready, err := oto.NewContext(decoded.SampleRate(), 2, 2)
	if err != nil {
		fmt.Println(err.Error())
	}
	<-ready

	player := context.NewPlayer(decoded)
	player.Play()

	for player.IsPlaying() {
		time.Sleep(time.Second)
	}
}
