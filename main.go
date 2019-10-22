package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"golang.org/x/text/encoding/charmap"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	city := flag.String("location", "08019", "City code from aemet URL (https://www.aemet.es/xml/municipios_h/localidad_h_08019.xml)")
	hours := flag.Int("n", 2, "Number of hours in the future to output")
	flag.Parse()

	resp, _ := http.Get("https://www.aemet.es/xml/municipios_h/localidad_h_" + *city + ".xml")
	dec := xml.NewDecoder(resp.Body)
	dec.CharsetReader = func(charset string, input io.Reader) (reader io.Reader, e error) {
		if charset != "ISO-8859-15" {
			return nil, errors.New("charset is not ISO-8859-15")
		}
		return charmap.ISO8859_15.NewDecoder().Reader(input), nil
	}

	forecast := Location{}
	err := dec.Decode(&forecast)
	if err != nil {
		log.Fatal(err)
	}

	forecast.Parse()

	out := ""
	for _, f := range forecast.NextHours(*hours) {
		out += f.String() + "  "
	}

	os.Stdout.WriteString(strings.TrimSpace(out))
}
