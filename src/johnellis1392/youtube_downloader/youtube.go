package main

// package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	apiURL = "https://www.youtube.com"
)

// Video - Video
type Video struct {
	ID           string
	URL          string
	VideoInfoRaw string
	VideoInfo    string
}

// FetchVideoInfo - FetchVideoInfo
func (v Video) FetchVideoInfo() (string, error) {
	endpoint := fmt.Sprintf("%s/get_video_info?video_id=%s", apiURL, v.ID)
	response, err := http.Get(endpoint)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	videoInfoRaw := string(body)
	videoInfo, err := parseVideoInfo(videoInfoRaw)
	if err != nil {
		return "", err
	}

	v.VideoInfo = videoInfo
	v.VideoInfoRaw = videoInfoRaw
	return v.VideoInfoRaw, nil
}

func parseVideoInfo(videoInfo string) (string, error) {
	query, err := url.ParseQuery(videoInfo)
	if err != nil {
		return "", err
	}

	dumpVideoInfo(query)
	return "", nil
}

// type Values map[string][]string
func dumpVideoInfo(values url.Values) {
	const LENGTH int = 50
	padding := getPadding(values)

	for key, v := range values {
		value := v[0]
		var truncatedValue string

		if len(value) > LENGTH {
			truncatedValue = value[:LENGTH]
		} else {
			truncatedValue = value
		}

		strpad := strings.Repeat(" ", padding-len(key))
		fmt.Printf("%s%s=> %v ...\n", key, strpad, truncatedValue)
	}
}

func getPadding(values url.Values) int {
	const paddingBuffer int = 2
	var maxLength int

	for key := range values {
		if len(key) > maxLength {
			maxLength = len(key)
		}
	}

	return maxLength + paddingBuffer
}

// NewVideo - NewVideo
func NewVideo(url string) Video {
	videoID := getVideoID(url)
	return Video{
		ID:  videoID,
		URL: url,
	}
}

func getVideoID(url string) string {
	re := regexp.MustCompile("^https://www\\.youtube\\.com/watch\\?v=(.*?)(?:&.*)?$")
	id := re.ReplaceAllString(url, "$1")
	return id
}
