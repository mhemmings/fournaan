package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// naanFile holds the bytes for the Peep Show gif.
var naanFile []byte

// impact is the Impact font face.
var impact font.Face

// numbers is a map on number strings to their human readable strings.
var numbers = map[string]string{
	"0":  "NO",
	"1":  "ONE",
	"2":  "TWO",
	"3":  "THREE",
	"4":  "FOUR",
	"5":  "FIVE",
	"6":  "SIX",
	"7":  "SEVEN",
	"8":  "EIGHT",
	"9":  "NINE",
	"10": "TEN",
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run performs setup then listens for http requests.
func run() error {
	var err error
	naanFile, err = ioutil.ReadFile("fournaan.gif")
	if err != nil {
		return err
	}

	fontBytes, err := ioutil.ReadFile("impact.ttf")
	if err != nil {
		return err
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return err
	}
	impact = truetype.NewFace(f, &truetype.Options{
		Size: 40,
	})

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/img/", imageGenerateHandler)

	return http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
}

func imageGenerateHandler(w http.ResponseWriter, req *http.Request) {
	text := strings.ToUpper(req.URL.Path[5:])

	name := strings.ToUpper(strings.TrimSpace(req.URL.Query().Get("name")))
	if name == "" {
		name = "JEREMY"
	}
	numQ := strings.ToUpper(strings.TrimSpace(req.URL.Query().Get("num")))
	num, ok := numbers[numQ]
	if !ok {
		num = "FOUR"
	}

	g, err := gif.DecodeAll(bytes.NewReader(naanFile))
	if err != nil {
		http.Error(w, "error loading gif", http.StatusInternalServerError)
		return
	}

	// Add text to each frame of the gif.
	for i, frame := range g.Image {
		// Find the closest colour in the palette to white, and set it to white
		// so that the text doesn't appear off-white.
		white := frame.Palette.Index(color.White)
		frame.Palette[white] = color.White

		var str string
		if i < 2 {
			continue
		}
		if i < 18 {
			str = fmt.Sprintf("%s %s, %s?", num, text, name)
		} else if i < 27 {
			str = fmt.Sprintf("%s?", num)
		} else {
			str = "THAT'S INSANE"
		}

		bounds, _ := font.BoundString(impact, str)
		d := &font.Drawer{
			Dst:  frame,
			Src:  image.White,
			Face: impact,
			Dot: fixed.Point26_6{
				fixed.Int26_6(
					(frame.Bounds().Max.X/2 - bounds.Max.X.Round()/2) * 64,
				),
				fixed.Int26_6(
					(frame.Bounds().Max.Y - 10) * 64,
				),
			},
		}
		d.DrawString(str)
	}

	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("Content-Type", "image/gif")

	err = gif.EncodeAll(w, g)
	if err != nil && !errors.Is(err, syscall.EPIPE) {
		http.Error(w, "error encoding gif", http.StatusInternalServerError)
		return
	}
}
