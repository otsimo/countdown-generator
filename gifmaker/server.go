package gifmaker

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

type GifServerConf struct {
	Port     int
	Dpi      float64
	FontFile string
}

func CreateGifServer(conf GifServerConf) {
	err := SetContext(conf)
	if err != nil {
		log.Fatal("Context could not be set err:", err)
	}
	http.HandleFunc("/countdown", CountdownRequest)             // set router
	err = http.ListenAndServe(":"+strconv.Itoa(conf.Port), nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// CountdownRequest Method that takes in a date and serves a countdown
func CountdownRequest(w http.ResponseWriter, r *http.Request) {
	str := r.URL.Query().Get("expires")
	fg := r.URL.Query().Get("fg")
	bg := r.URL.Query().Get("bg")
	wd := r.URL.Query().Get("width")
	if wd == "" {
		wd = "500"
	}
	width, _ := strconv.Atoi(wd)
	if str == "" {
		str = "2016-01-01T01:01:01"
	}
	expires, err := time.Parse(
		"2006-01-02T15:04:05",
		str)
	if err != nil {
		w.Write([]byte("Error Parsing the Time"))
	}
	w.Header().Set("Content-Type", "image/gif")
	w.Header().Set("Cache-Control", "no-cache")
	// expires = expires.AddDate(0, -3, 0)
	expires = expires.Add(time.Duration(-3) * time.Hour)

	gifBuffer, err := MakeGif(expires, width, bg, fg)
	if err != nil {
		w.Write([]byte("Error Creating the GIF"))
	}
	w.Write(gifBuffer.Bytes())

}
