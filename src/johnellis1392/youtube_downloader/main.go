package main

// Look at this for how to download youtube videos:
// https://github.com/kkdai/youtube
import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	DEBUG      = true
	defaultURL = "https://www.youtube.com/watch?v=-xgNc3-nBpI"
	outfile    = "video_data.txt"
)

func log(message string) {
	if DEBUG {
		fmt.Println(message)
	}
}

func saveVideo(video Video, dest string) {
	f, err := os.Create(dest)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	video_json, err := json.MarshalIndent(video, "", "  ")
	if err != nil {
		panic(err)
	}

	log("Saving json...")
	_, err = f.Write(video_json)
	if err != nil {
		panic(err)
	}
}

/* ************ */
/* *** Main *** */
/* ************ */
func main() {
	var urlstr string
	if len(os.Args) == 2 {
		urlstr = os.Args[1]
	} else {
		urlstr = defaultURL
	}

	video, err := FetchVideo(urlstr)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(outfile)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	encodedData, err := json.MarshalIndent(video, "", "  ")
	if err != nil {
		panic(err)
	}

	_, err = f.Write(encodedData)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success")
}
