package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
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

	resp, err := http.Get("https://www.aemet.es/xml/municipios_h/localidad_h_" + *city + ".xml")
	if err != nil {
		fmt.Println(" Network error") // https://fontawesome.com/icons/exclamation-triangle?style=solid
	}

	dec := xml.NewDecoder(resp.Body)
	dec.CharsetReader = func(charset string, input io.Reader) (reader io.Reader, e error) {
		switch charset {
		case "ISO-8859-15":
			return charmap.ISO8859_15.NewDecoder().Reader(input), nil
		default:
			return nil, errors.New("charset is not ISO-8859-15")
		}
	}

	forecast := Location{}
	err = dec.Decode(&forecast)
	if err != nil {
		fmt.Println(" Parse error")
		log.Fatal(err)
	}

	forecast.Parse()

	out := ""
	for _, f := range forecast.NextHours(*hours) {
		out += f.String() + "  "
	}

	os.Stdout.WriteString(strings.TrimSpace(out))
}
