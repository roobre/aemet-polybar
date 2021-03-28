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
	"strconv"
	"strings"
	"time"
)

func main() {
	city := flag.String("location", "08019", "City code from aemet URL (https://www.aemet.es/xml/municipios_h/localidad_h_08019.xml)")
	hours := flag.Int("n", 2, "Number of hours in the future to output")
	labels := flag.String("labels", "", "List of hours and labels to show (e.g. \":13,19h:19\")")
	separator := flag.String("separator", "  ", "Separate entries with this string.")
	flag.Parse()

	resp, err := http.Get("https://www.aemet.es/xml/municipios_h/localidad_h_" + *city + ".xml")
	if err != nil {
		fmt.Println(" net") // https://fontawesome.com/icons/exclamation-triangle?style=solid
		log.Fatal(err)
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
		fmt.Println(" fmt")
		log.Fatal(err)
	}

	forecast.Parse()

	out := ""

	if *labels != "" {
		for _, pair := range strings.Split(*labels, ",") {
			labelHour := strings.Split(pair, ":")
			if len(labelHour) != 2 {
				continue
			}

			hour, err := strconv.Atoi(labelHour[1])
			if err != nil {
				continue
			}

			if time.Now().Hour() > hour {
				continue
			}

			hourForecast := forecast.At(hour)
			if hourForecast == nil {
				continue
			}

			if labelHour[0] != "" {
				out += labelHour[0] + ": "
			}
			out += hourForecast.String() + *separator
		}
	} else {
		for _, f := range forecast.NextHours(*hours) {
			out += f.String() + "  "
		}
	}

	os.Stdout.WriteString(strings.TrimSpace(out))
}
