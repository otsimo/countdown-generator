package main

/*
Server struct that defines the basics of countdown generator,
if serve option is static all the possible fragments in &GifLengthInMinutes length
are created into the public/gifs/ directory and each request calculates the difference
between now and the requested EndDate and serves the appropriate gif image, if dynamic option is provided
the appropriate gif image is created on the fly and is served dynamically
*/
type Server struct {
	ServeOption        string
	GifLengthInMinutes int
	MaxInterval        int
}

func main() {

}
