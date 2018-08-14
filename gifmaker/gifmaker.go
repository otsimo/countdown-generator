package gifmaker

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/golang/freetype/truetype"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

type GifMaker interface {
	MakeGif(expires time.Time) (*bytes.Buffer, error)
}

func NewGifMaker(g Config) (GifMaker, error) {
	c := freetype.NewContext()

	if g.Fg == "" {
		g.Fg = "000000"
	}
	fontColor, err := hexToRGBA(g.Fg)
	if err != nil {
		fontColor = color.NRGBA{0, 0, 0, 255}
	}
	fmt.Println("Font Color", g.Fg)

	if g.Bg == "" {
		g.Bg = "ffffff"
	}
	bCol, err := hexToRGBA(g.Bg)
	if err != nil {
		bCol = color.NRGBA{255, 255, 255, 255}
	}
	fmt.Println("Backfground Color", g.Bg)

	if g.FontSize == 0 {
		g.FontSize = 16
	}

	c.SetDPI(g.Dpi)
	c.SetFont(g.Font)
	c.SetHinting(font.HintingNone)
	c.SetSrc(image.NewUniform(fontColor))
	c.SetFontSize(g.FontSize)

	return &gifMaker{
		context:  *c,
		bg:       bCol,
		fg:       fontColor,
		font:     g.Font,
		fontSize: g.FontSize,
		dpi:      g.Dpi,
	}, nil
}

type Config struct {
	FontSize float64
	Dpi      float64
	Font     *truetype.Font
	Fg       string
	Bg       string
}

type gifMaker struct {
	bg       color.Color
	fg       color.Color
	font     *truetype.Font
	fontSize float64
	dpi      float64
	context  freetype.Context
}

func (gm *gifMaker) MakeGif(expires time.Time) (*bytes.Buffer, error) {
	out := &gif.GIF{}
	now := time.Now()

	for n := 0; n < 60; n++ {
		dif := expires.Sub(now)
		timeString := getTimeFragments(dif)
		img, err := gm.createFrame(timeString)
		if err != nil {
			return nil, err
		}
		out.Image = append(out.Image, img)
		out.Delay = append(out.Delay, 100)
		expires = expires.Add(time.Duration(-1) * time.Second)
	}
	buf := new(bytes.Buffer)
	gif.EncodeAll(buf, out)

	return buf, nil
}

func getTimeFragments(dif time.Duration) (timeString string) {
	var days, hours, minutes, seconds float64
	if dif > 0 {
		days = math.Floor(dif.Seconds() / (60 * 60 * 24))
		hours = math.Floor((dif.Seconds()/(60*60) - (days * 24)))
		minutes = math.Floor((dif.Seconds()/(60) - (days * 24 * 60) - (hours * 60)))
		seconds = math.Floor(dif.Seconds() - (days * 60 * 60 * 24) - (hours * 60 * 60) - (minutes * 60))
		if days == 0 {
			return fmt.Sprintf("%02.f:%02.f:%02.f", hours, minutes, seconds)
		}
		if days == 0 && hours == 0 {
			return fmt.Sprintf("%02.f:%02.f", minutes, seconds)
		} else {
			return fmt.Sprintf("%02.f:%02.f:%02.f:%02.f", days, hours, minutes, seconds)
		}
	} else {
		days = 0
		hours = 0
		minutes = 0
		seconds = 0
		return fmt.Sprintf("%02.f:%02.f:%02.f:%02.f", days, hours, minutes, seconds)
	}
	return ""
}

func (gm *gifMaker) createFrame(timeString string) (*image.Paletted, error) {
	// palette := []color.Color{gm.fg, gm.bg}
	var palette = color.Palette{color.White, color.Black, gm.fg, gm.bg}
	backgroundHeight := 0
	backgroundWidth := 0

	for _, letter := range timeString {
		nf := truetype.NewFace(gm.font, &truetype.Options{Size: gm.fontSize, DPI: gm.dpi, Hinting: font.HintingNone})
		r, a, _ := nf.GlyphBounds(letter)
		if backgroundHeight == 0 {
			backgroundHeight = int(math.Abs(float64(r.Min.Y.Round() - r.Max.Y.Round())))
		}
		backgroundWidth += int(a.Round())
	}
	backgroundWidth = int(float64(backgroundWidth) * 1.1)
	backgroundHeight = int(float64(backgroundHeight) * 1.1)
	background := image.NewPaletted(image.Rect(0, 0, backgroundWidth, backgroundHeight), palette)
	gm.context.SetDst(background)
	gm.context.SetClip(background.Bounds())
	fontBackGroundColor := image.NewUniform(gm.bg)
	draw.Draw(background, background.Bounds(), fontBackGroundColor, image.ZP, draw.Src)
	pt := freetype.Pt(int(float64(backgroundWidth)*4.5/100), backgroundHeight*19/20)
	_, err := gm.context.DrawString(timeString, pt)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return background, nil
}

func hexToRGBA(hex string) (color.Color, error) {
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
