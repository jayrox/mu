package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s URL\n", os.Args[0])
		os.Exit(1)
	}
	response, err := http.Get(os.Args[1])

	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		buf := bytes.NewBuffer(make([]byte, 0, response.ContentLength+1))
		_, readErr := buf.ReadFrom(response.Body)
		if readErr != nil {
			panic(err)
		}
		rb := buf.Bytes()
		fmt.Println(string(rb))

		var r radarrMovie
		err = json.Unmarshal(rb, &r)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Title: %s\n", r.Title)
		if r.Title == "" {
			fmt.Println("Invalid title.")
			return
		}

		fmt.Printf("Path %s\n", r.Path)
		firstChar := strings.ToUpper(string(r.SortTitle[0]))

		commitChanges := false
		if strings.Contains(r.Path, "_Unsorted") {
			fmt.Printf("First Char: %s\n", firstChar)
			newPath := firstChar
			r.Path = strings.Replace(r.Path, "_Unsorted", newPath, -1)
			fmt.Printf("New Path: %s\n", r.Path)
			commitChanges = true
		}

		if r.Status == "" {
			fmt.Printf("Status - Replacing %s => %s\n", r.Status, "released")
			r.Status = "released"
			commitChanges = true
		}

		if r.Monitored == false {
			fmt.Printf("Monitored - Replacing %s => %s\n", "false", "true")
			r.Monitored = true
			commitChanges = true
		}

		if r.MinimumAvailability != "released" {
			fmt.Printf("MinAvail - Replacing %s => %s\n", r.MinimumAvailability, "released")
			r.MinimumAvailability = "released"
			commitChanges = true
		}

		if r.PathState != "static" {
			fmt.Printf("PathState - Replacing %s => %s\n", r.PathState, "static")
			r.PathState = "static"
			commitChanges = true
		}

		if r.ProfileID != 6 {
			fmt.Printf("ProfileID - Replacing %d => %d\n", r.ProfileID, 6)
			r.ProfileID = 6
			commitChanges = true
		}

		// If no changes need to be commit, go ahead and exit.
		if commitChanges == false {
			fmt.Println("Radarr settings are correct.")
			return
		}

		// Only update this if there are other changes to be committed.
		if r.MovieFile.MediaInfo.RunTime == "" {
			r.MovieFile.MediaInfo.RunTime = "0"
		}

		fmt.Println("Radarr settings are incorrect.")
		fmt.Println("Updating Radarr settings.")

		b, err := json.Marshal(r)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))

		putRequest(os.Args[1], bytes.NewBuffer(b))
	}
}

func putRequest(url string, data io.Reader) {
	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, data)
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	//fmt.Printf("%+v", res)
	if res.Status == "202 Accepted" {
		fmt.Printf("Res: %s\n", res.Status)
		return
	}
	fmt.Printf("Res: %s\n", res.Status)
}

type radarrMovie struct {
	Title             string `json:"title"`
	AlternativeTitles []struct {
		SourceType string `json:"sourceType"`
		MovieID    int    `json:"movieId"`
		Title      string `json:"title"`
		SourceID   int    `json:"sourceId"`
		Votes      int    `json:"votes"`
		VoteCount  int    `json:"voteCount"`
		Language   string `json:"language"`
		ID         int    `json:"id"`
	} `json:"alternativeTitles"`
	SecondaryYearSourceID int       `json:"secondaryYearSourceId"`
	SortTitle             string    `json:"sortTitle"`
	SizeOnDisk            int64     `json:"sizeOnDisk"`
	Status                string    `json:"status"`
	Overview              string    `json:"overview"`
	InCinemas             time.Time `json:"inCinemas"`
	Images                []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
	} `json:"images"`
	Website             string        `json:"website"`
	Downloaded          bool          `json:"downloaded"`
	Year                int           `json:"year"`
	HasFile             bool          `json:"hasFile"`
	YouTubeTrailerID    string        `json:"youTubeTrailerId"`
	Studio              string        `json:"studio,omitempty"`
	Path                string        `json:"path"`
	ProfileID           int           `json:"profileId"`
	PathState           string        `json:"pathState"`
	Monitored           bool          `json:"monitored"`
	MinimumAvailability string        `json:"minimumAvailability"`
	IsAvailable         bool          `json:"isAvailable"`
	FolderName          string        `json:"folderName"`
	Runtime             int           `json:"runtime"`
	LastInfoSync        time.Time     `json:"lastInfoSync"`
	CleanTitle          string        `json:"cleanTitle"`
	ImdbID              string        `json:"imdbId"`
	TmdbID              int           `json:"tmdbId"`
	TitleSlug           string        `json:"titleSlug"`
	Genres              []string      `json:"genres"`
	Tags                []interface{} `json:"tags"`
	Added               time.Time     `json:"added"`
	Ratings             struct {
		Votes int     `json:"votes"`
		Value float64 `json:"value"`
	} `json:"ratings"`
	MovieFile struct {
		MovieID      int       `json:"movieId"`
		RelativePath string    `json:"relativePath"`
		Size         int64     `json:"size"`
		DateAdded    time.Time `json:"dateAdded"`
		Quality      struct {
			Quality struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"quality"`
			Revision struct {
				Version int `json:"version"`
				Real    int `json:"real"`
			} `json:"revision"`
		} `json:"quality"`
		Edition   string `json:"edition"`
		MediaInfo struct {
			VideoCodec                string  `json:"videoCodec"`
			VideoBitrate              int     `json:"videoBitrate"`
			VideoBitDepth             int     `json:"videoBitDepth"`
			Width                     int     `json:"width"`
			Height                    int     `json:"height"`
			AudioFormat               string  `json:"audioFormat"`
			AudioBitrate              int     `json:"audioBitrate"`
			RunTime                   string  `json:"runTime"`
			AudioStreamCount          int     `json:"audioStreamCount"`
			AudioChannels             int     `json:"audioChannels"`
			AudioChannelPositions     string  `json:"audioChannelPositions"`
			AudioChannelPositionsText string  `json:"audioChannelPositionsText"`
			AudioProfile              string  `json:"audioProfile"`
			VideoFps                  float64 `json:"videoFps"`
			AudioLanguages            string  `json:"audioLanguages"`
			Subtitles                 string  `json:"subtitles"`
			ScanType                  string  `json:"scanType"`
			SchemaRevision            int     `json:"schemaRevision"`
		} `json:"mediaInfo"`
		ID int `json:"id"`
	} `json:"movieFile"`
	QualityProfileID int `json:"qualityProfileId"`
	ID               int `json:"id"`
}
