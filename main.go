package main

import (
	"encoding/xml"
	"errors"
	"golang.org/x/text/encoding/charmap"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	resp, _ := http.Get("https://www.aemet.es/xml/municipios_h/localidad_h_08019.xml")
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
	for _, f := range forecast.NextHours(2) {
		out += f.String() + "	 "
	}

	os.Stdout.WriteString(strings.TrimSpace(out))
}
