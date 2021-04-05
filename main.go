package main

import (
	"flag"
	"fmt"
	"log"
	"roob.re/aemet-polybar/aemet"
	"strconv"
	"strings"
	"time"
)

func main() {
	city := flag.String("location", "08019", "City code from aemet URL (https://www.aemet.es/xml/municipios_h/localidad_h_08019.xml)")
	hours := flag.Int("n", 2, "Number of hours in the future to output")
	labels := flag.String("labels", "", "List of hours and labels to show (e.g. \":13,19h:19\")")
	separator := flag.String("separator", "  ", "Separate entries with this string")
	flag.Parse()

	forecast, err := aemet.City(*city)
	if err != nil {
		log.Fatal(err)
	}

	out := ""

	if *labels != "" {
		t := time.Now()
		for _, pair := range strings.Split(*labels, ",") {
			labelHour := strings.Split(pair, ":")
			if len(labelHour) != 2 {
				log.Printf("cannot parse label '%s'", pair)
				continue
			}

			hour, err := strconv.Atoi(labelHour[1])
			if err != nil {
				log.Printf("cannot parse hour from label '%s': %v", labelHour[1], err)
				continue
			}

			if t.Hour() > hour {
				log.Printf("hour %d is in the past, skipping", hour)
				continue
			}

			hourForecast := forecast.At(time.Date(
				t.Year(), t.Month(), t.Day(), hour, 0, 0, 0, t.Location(),
			))
			if hourForecast == nil {
				log.Printf("got nil forecast for %dh, skipping", hour)
				continue
			}

			if labelHour[0] != "" {
				out += labelHour[0] + ": "
			}
			out += hourForecast.String() + *separator
		}
	} else {
		for _, f := range forecast.NextHours(*hours) {
			out += f.String() + *separator
		}
	}

	fmt.Println(strings.TrimSuffix(out, *separator))
}
