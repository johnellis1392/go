package main

// package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

const (
	apiURL = "https://www.youtube.com"
)

type VideoInfo struct {
}

// Video - Video
type Video struct {
	ID    string `json:"video_id"`
	URL   string
	Title string `json:"title"`

	Info VideoInfo         `json:"video_info"`
	Data map[string]string `json:"data"`

	AdaptiveFmtsRaw string `json:"adaptive_fmts"`
	AdaptiveFmts    string

	StreamsRaw string            `json:"url_encoded_fmt_stream_map"`
	Streams    map[string]string `json:"streams"`

	/*

		account_playback_token
		adaptive_fmts
		allow_embed
		allow_ratings
		apiary_host
		apiary_host_firstparty
		atc
		author
		avg_rating
		c
		cl
		cr
		csi_page_type
		csn
		cver
		enablecsi
		eventid
		external_play_video
		fexp
		fflags
		fmt_list
		gapi_hint_params
		hl
		host_language
		idpj
		innertube_api_key
		innertube_api_version
		innertube_context_client_version
		is_listed
		ismb
		itct
		keywords
		ldpj
		length_seconds
		loudness
		no_get_video_log
		of
		player_error_log_fraction
		player_response
		plid
		pltype
		ppv_remarketing_url
		probe_url
		ptk
		relative_loudness
		root_ve_type
		ssl
		status
		storyboard_spec
		swf_player_response
		t
		thumbnail_url
		timestamp
		title
		tmi
		token
		ucid
		url_encoded_fmt_stream_map
		video_id
		videostats_playback_base_url
		view_count
		vm
		vmap
		vss_host
		watermark
		xhr_apiary_host

	*/
}

func FetchVideo(url string) (*Video, error) {
	var err error
	var video Video
	video.ID = getVideoID(url)

	endpoint := fmt.Sprintf("%s/get_video_info?video_id=%s", apiURL, video.ID)
	response, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	jsonstr, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	decodedJson, err := strconv.Unquote("\"" + string(jsonstr) + "\"")
	if err != nil {
		return nil, err
	}

	decodedVideo, err := parseQueryString(string(decodedJson))
	if err != nil {
		return nil, err
	}

	video.Data = decodedVideo

	// Get Adaptive Formats
	// video.AdaptiveFmtsRaw = strconv.Unquote("\"" + video.Data[""] + "\"")
	videoJson, err := json.Marshal(decodedVideo)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(videoJson, &video)
	if err != nil {
		return nil, err
	}

	streamMap, err := parseQueryString(video.StreamsRaw)
	if err != nil {
		return nil, err
	}

	video.Streams = streamMap

	// Get Stream Map

	// _, err = video.FetchVideoInfo()
	// if err != nil {
	// 	return nil, err
	// }

	return &video, nil
}

// FetchVideoInfo - FetchVideoInfo
func (v Video) FetchVideoInfo() (*VideoInfo, error) {
	var videoInfo VideoInfo

	endpoint := fmt.Sprintf("%s/get_video_info?video_id=%s", apiURL, v.ID)
	response, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	videoInfoRaw := string(body)
	// fmt.Println(videoInfoRaw)
	data, err := parseQueryString(videoInfoRaw)
	if err != nil {
		return nil, err
	}

	datastr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(datastr, &videoInfo)
	if err != nil {
		return nil, err
	}

	v.Info = videoInfo
	return &videoInfo, nil
}

func parseQueryString(querystr string) (map[string]string, error) {
	query, err := url.ParseQuery(querystr)
	if err != nil {
		return nil, err
	}

	params := make(map[string]string)
	for key, value := range query {
		s, err := strconv.Unquote("\"" + value[0] + "\"")
		if err != nil {
			continue
		}

		params[key] = s
	}

	return params, nil
}

func getVideoID(url string) string {
	re := regexp.MustCompile("^https://www\\.youtube\\.com/watch\\?v=(.*?)(?:&.*)?$")
	id := re.ReplaceAllString(url, "$1")
	return id
}
