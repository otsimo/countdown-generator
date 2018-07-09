package main

import (
	"flag"

	"github.com/otsimo/countdown-generator/gifmaker"
)

var conf = gifmaker.GifServerConf{
	Port: 8090,
}

func main() {
	flag.Parse()
	gifmaker.CreateGifServer(conf)

}
