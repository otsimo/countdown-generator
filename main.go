package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang/freetype/truetype"

	"github.com/golang/freetype"
	"github.com/otsimo/countdown-generator/gifmaker"
)

var PTM *truetype.Font
var Arial *truetype.Font
var OpenSans *truetype.Font

func cacheFonts() error {
	openSansFontBytes, err := ioutil.ReadFile("./gifmaker/fonts/OpenSans.ttf")
	if err != nil {
		log.Println(err)
		return err
	}
	openSans, err := freetype.ParseFont(openSansFontBytes)
	if err != nil {
		log.Println(err)
		return err
	}
	OpenSans = openSans

	ptmFontBytes, err := ioutil.ReadFile("./gifmaker/fonts/PTM55FT.ttf")
	if err != nil {
		log.Println(err)
		return err
	}
	ptm, err := freetype.ParseFont(ptmFontBytes)
	if err != nil {
		log.Println(err)
		return err
	}
	PTM = ptm

	arialFontBytes, err := ioutil.ReadFile("./gifmaker/fonts/arial.ttf")
	if err != nil {
		log.Println(err)
		return err
	}
	arial, err := freetype.ParseFont(arialFontBytes)
	if err != nil {
		log.Println(err)
		return err
	}
	Arial = arial

	return nil
}

func main() {
	err := cacheFonts()
	if err != nil {
		log.Fatal("Fonts could not be read")
	}
	http.HandleFunc("/countdownArial", CountdownRequestArial)       // set router
	http.HandleFunc("/countdownOpenSans", CountdownRequestOpenSans) // set router
	http.HandleFunc("/countdownPTM", CountdownRequestPTM)           // set router
	err = http.ListenAndServe(":8090", nil)                         // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func CountdownRequestArial(w http.ResponseWriter, r *http.Request) {
	str := r.URL.Query().Get("expires")
	fg := r.URL.Query().Get("fg")
	bg := r.URL.Query().Get("bg")
	fontSize, _ := strconv.ParseFloat(r.URL.Query().Get("fontSize"), 64)

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

	conf := gifmaker.Config{
		FontSize: fontSize,
		Dpi:      72,
		Font:     Arial,
		Fg:       fg,
		Bg:       bg,
	}
	gm, err := gifmaker.NewGifMaker(conf)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	gifBuffer, err := gm.MakeGif(expires)
	if err != nil {
		w.Write([]byte("Error Creating the GIF"))
	}
	w.Write(gifBuffer.Bytes())

}
func CountdownRequestPTM(w http.ResponseWriter, r *http.Request) {
	str := r.URL.Query().Get("expires")
	fg := r.URL.Query().Get("fg")
	bg := r.URL.Query().Get("bg")
	fontSize, _ := strconv.ParseFloat(r.URL.Query().Get("fontSize"), 64)
	tm := r.URL.Query().Get("marker")
	timeMarker := false

	if tm == "true" {
		timeMarker = true
	}

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

	conf := gifmaker.Config{
		FontSize:   fontSize,
		Dpi:        72,
		Font:       PTM,
		Fg:         fg,
		Bg:         bg,
		TimeMarker: timeMarker,
	}
	gm, err := gifmaker.NewGifMaker(conf)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	gifBuffer, err := gm.MakeGif(expires)
	if err != nil {
		w.Write([]byte("Error Creating the GIF"))
	}
	w.Write(gifBuffer.Bytes())

}

func CountdownRequestOpenSans(w http.ResponseWriter, r *http.Request) {
	str := r.URL.Query().Get("expires")
	fg := r.URL.Query().Get("fg")
	bg := r.URL.Query().Get("bg")
	fontSize, _ := strconv.ParseFloat(r.URL.Query().Get("fontSize"), 64)
	tm := r.URL.Query().Get("marker")
	timeMarker := false

	if tm == "true" {
		timeMarker = true
	}
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

	conf := gifmaker.Config{
		FontSize:         fontSize,
		Dpi:              72,
		Font:             OpenSans,
		Fg:               fg,
		Bg:               bg,
		TimeMarker:       timeMarker,
		MarkerFontOffset: 0.85,
	}
	gm, err := gifmaker.NewGifMaker(conf)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	gifBuffer, err := gm.MakeGif(expires)
	if err != nil {
		w.Write([]byte("Error Creating the GIF"))
	}
	w.Write(gifBuffer.Bytes())

}
