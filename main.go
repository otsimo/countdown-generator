package main

import (
	"flag"

	"github.com/otsimo/countdown-generator/gifmaker"
)

var conf = gifmaker.GifServerConf{
	Port: 8090,
}

func main() {
	ff := flag.String("fontfile", "fonts/luxisr.ttf", "filename of the ttf font")
	dp := flag.Float64("dpi", 60, "screen resolution in Dots Per Inch")
	flag.Parse()
	conf.FontFile = *ff
	conf.Dpi = *dp
	gifmaker.CreateGifServer(conf)

}
