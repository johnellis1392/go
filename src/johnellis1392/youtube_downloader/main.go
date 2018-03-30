package main

// Look at this for how to download youtube videos:
// https://github.com/kkdai/youtube
import (
	"fmt"
)

const (
	defaultURL = "https://www.youtube.com/watch?v=-xgNc3-nBpI"
)

/* ************ */
/* *** Main *** */
/* ************ */
func main() {
	fmt.Println("Fetching Video...")
	video := NewVideo(defaultURL)

	fmt.Println("Fetching Video Data...")
	// data, err := video.FetchVideoInfo()
	_, err := video.FetchVideoInfo()
	if err != nil {
		fmt.Println("An error occurred:", err.Error())
	}

	// fmt.Println("Success:")
	// fmt.Println(data)
}
