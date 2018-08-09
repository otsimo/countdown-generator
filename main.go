package main

import (
	"flag"

	"github.com/otsimo/countdown-generator/gifmaker"
)

var conf = gifmaker.GifServerConf{
	Port:     8090,
	FontFile: *flag.String("fontfile", "fonts/luxisr.ttf", "filename of the ttf font"),
	Dpi:      *flag.Float64("dpi", 72, "screen resolution in Dots Per Inch"),
}

func main() {
	flag.Parse()
	gifmaker.CreateGifServer(conf)

}
