package gifmaker

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

var palette = []color.Color{color.White, color.Black}
var context *freetype.Context

func MakeGif(expires time.Time, maxWidth int, bg string, fg string) (bytes.Buffer, error) {
	palette = []color.Color{color.White, color.Black}
	out := &gif.GIF{}
	now := time.Now()
	if fg == "" {
		fg = "000000"
	}
	fontColor, _ := hexToRGBA(fg)
	if fontColor != nil {
		palette = append(palette, fontColor)
		context.SetSrc(image.NewUniform(fontColor))
	}
	if bg == "" {
		bg = "ffffff"
	}

	bCol, _ := hexToRGBA(bg)
	if bCol != nil {
		palette = append(palette, bCol)
	}

	for n := 0; n < 60; n++ {
		dif := expires.Sub(now)
		offset := maxWidth / 50
		timeString := GetTimeFragments(dif)
		img := image.NewPaletted(image.Rect(0, 0, maxWidth, maxWidth/5), palette)
		draw.Draw(img, img.Bounds(), &image.Uniform{bCol}, image.ZP, draw.Src)
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
	var days, hours, minutes, seconds float64
	if dif > 0 {
		days = math.Floor(dif.Seconds() / (60 * 60 * 24))
		hours = math.Floor((dif.Seconds()/(60*60) - (days * 24)))
		minutes = math.Floor((dif.Seconds()/(60) - (days * 24 * 60) - (hours * 60)))
		seconds = math.Floor(dif.Seconds() - (days * 60 * 60 * 24) - (hours * 60 * 60) - (minutes * 60))
	} else {
		days = 0
		hours = 0
		minutes = 0
		seconds = 0
	}
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

func SetContext(conf GifServerConf) error {
	fontBytes, err := ioutil.ReadFile(conf.FontFile)
	if err != nil {
		log.Println(err)
		return err
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return err
	}

	fg := image.NewUniform(color.NRGBA{255, 0, 0, 255})

	c := freetype.NewContext()
	c.SetDPI(conf.Dpi)

	c.SetFont(f)
	c.SetHinting(font.HintingNone)
	c.SetSrc(fg)

	context = c
	return nil
}

func hexToRGBA(hex string) (color.Color, error) {
	log.Printf("%s", hex)
	if len(hex) != 6 {
		return nil, errors.New("Color Provided not Hex")
	}
	r, err := strconv.ParseInt(hex[0:2], 16, 64)
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}
	g, err := strconv.ParseInt(hex[2:4], 16, 64)
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}
	b, err := strconv.ParseInt(hex[4:6], 16, 64)
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	a := 255
	return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}, nil
}
