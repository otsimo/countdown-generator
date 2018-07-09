package gifmaker

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io/ioutil"
	"log"
	"math"
	"time"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

var (
	dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "../../testdata/luxisr.ttf", "filename of the ttf font")
	hinting  = flag.String("hinting", "none", "none | full")
	spacing  = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	wonb     = flag.Bool("whiteonblack", false, "white text on a black background")
)
var palette = []color.Color{color.White, color.Black}
var context *freetype.Context

func MakeGif(expires time.Time, maxWidth int) (bytes.Buffer, error) {
	out := &gif.GIF{}
	now := time.Now()
	for n := 0; n < 60; n++ {
		dif := expires.Sub(now)
		offset := maxWidth / 14
		timeString := GetTimeFragments(dif)
		img := image.NewPaletted(image.Rect(0, 0, maxWidth, maxWidth/4), palette)
		AddLabel(img, offset, 0, timeString, float64(maxWidth)*0.175)
		out.Image = append(out.Image, img)
		out.Delay = append(out.Delay, 100)
		expires = expires.Add(time.Duration(-1) * time.Second)
	}
	buf := new(bytes.Buffer)
	gif.EncodeAll(buf, out)

	return *buf, nil
}

func GetTimeFragments(dif time.Duration) (timeString string) {
	days := math.Floor(dif.Seconds() / (60 * 60 * 24))
	hours := math.Floor((dif.Seconds()/(60*60) - (days * 24)))
	minutes := math.Floor((dif.Seconds()/(60) - (days * 24 * 60) - (hours * 60)))
	seconds := math.Floor(dif.Seconds() - (days * 60 * 60 * 24) - (hours * 60 * 60) - (minutes * 60))
	return fmt.Sprintf("%02.f:%02.f:%02.f:%02.f", days, hours, minutes, seconds)
}

//AddLabel function that takes in maxWidth, labels and their locations in x and y coordinates
func AddLabel(img *image.Paletted, x, y int, label string, fontSize float64) error {
	context.SetClip(img.Bounds())
	context.SetDst(img)
	context.SetFontSize(fontSize)
	pt := freetype.Pt(x, y+int(context.PointToFixed(fontSize)>>6))
	_, err := context.DrawString(label, pt)
	if err != nil {
		return err
	}
	return nil
}

func SetContext() error {
	fontBytes, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		log.Println(err)
		return err
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return err
	}
	fg := image.Black

	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(f)
	c.SetHinting(font.HintingNone)
	c.SetSrc(fg)

	context = c
	return nil
}