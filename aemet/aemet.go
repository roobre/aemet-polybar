package aemet

import (
	"encoding/xml"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "2006-01-02"

type Client struct {
	HttpClient *http.Client
}

func (c *Client) City(cityCode string) (*Location, error) {
	resp, err := http.Get("https://www.aemet.es/xml/municipios_h/localidad_h_" + url.QueryEscape(cityCode) + ".xml")
	if err != nil {
		return nil, err
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

	loc := &Location{}
	err = dec.Decode(&loc)
	if err != nil {
		return nil, err
	}

	loc.parse()
	return loc, nil
}

var defaultClient Client = Client{HttpClient: http.DefaultClient}

func City(cityCode string) (*Location, error) {
	return defaultClient.City(cityCode)
}

type Location struct {
	Name          string `xml:"nombre"`
	Region        string `xml:"provincia"`
	ElaboratedStr string `xml:"elaborado"`

	Forecasts      []DailyForecast `xml:"prediccion>dia" json:"-"`
	DailyForecasts map[string]*DailyForecast
}

func (l *Location) NextHours(n int) []*ParsedForecast {
	forecasts := make([]*ParsedForecast, 0, n)

	for i := 1; i <= n; i++ {
		hour := time.Now().Add(-15 * time.Minute).Truncate(time.Hour).Add(time.Duration(i) * time.Hour)
		forecasts = append(forecasts, l.At(hour))
	}

	return forecasts
}

func (l *Location) At(t time.Time) *ParsedForecast {
	day := t.Format(dateFormat)
	f := l.DailyForecasts[day]

	if f == nil {
		return nil
	}

	hour := t.Hour()
	return &ParsedForecast{
		Location: l.Name,
		DateStr:  f.DateStr,
		Hour:     hour,

		SkyState:        f.HourlySkyState[hour],
		Precipitation:   f.HourlyPrecipitation[hour],
		POPPercent:      f.HourlyPOP[hour],
		Temperature:     f.HourlyTemperature[hour],
		ThermalFeel:     f.HourlyThermalFeel[hour],
		HumidityPercent: f.HourlyHumidity[hour],
	}
}

func (l *Location) parse() {
	l.DailyForecasts = map[string]*DailyForecast{}

	for i := range l.Forecasts {
		f := &l.Forecasts[i]
		f.parse()
		l.DailyForecasts[f.DateStr] = f
	}
}

type DailyForecast struct {
	DateStr string `xml:"fecha,attr"`

	SkyStates []struct {
		Hour        int    `xml:"periodo,attr"`
		Description string `xml:"descripcion,attr"`
		State       string `xml:",chardata"`
	} `xml:"estado_cielo" json:"-"`
	HourlySkyState [24]string

	Precipitations []struct {
		Hour   int     `xml:"periodo,attr"`
		Amount float32 `xml:",chardata"`
	} `xml:"precipitations" json:"-"`
	HourlyPrecipitation [24]float32

	POPs []struct {
		Period     string `xml:"periodo,attr"`
		POPPercent int    `xml:",chardata"`
	} `xml:"prob_precipitacion" json:"-"`
	HourlyPOP [24]int

	Temperatures []struct {
		Hour        int `xml:"periodo,attr"`
		Temperature int `xml:",chardata"`
	} `xml:"temperatura" json:"-"`
	HourlyTemperature [24]int

	ThermalFeels []struct {
		Hour        int `xml:"periodo,attr"`
		ThermalFeel int `xml:",chardata"`
	} `xml:"sens_termica" json:"-"`
	HourlyThermalFeel [24]int

	Humidities []struct {
		Hour            int `xml:"periodo,attr"`
		HumidityPercent int `xml:",chardata"`
	} `xml:"humedad_relativa" json:"-"`
	HourlyHumidity [24]int
}

func (f *DailyForecast) parse() {
	for _, ss := range f.SkyStates {
		f.HourlySkyState[ss.Hour] = ss.State
	}

	for _, p := range f.Precipitations {
		f.HourlyPrecipitation[p.Hour] = p.Amount
	}

	for _, pop := range f.POPs {
		if len(pop.Period) != 4 {
			continue
		}
		begin, err := strconv.Atoi(pop.Period[:2])
		if err != nil {
			continue
		}
		end, err := strconv.Atoi(pop.Period[2:])
		if err != nil {
			continue
		}

		hour := begin % 24
		for hour != end {
			f.HourlyPOP[hour] = pop.POPPercent
			hour = (hour + 1) % 24
		}
	}

	for _, t := range f.Temperatures {
		f.HourlyTemperature[t.Hour] = t.Temperature
	}

	for _, tf := range f.ThermalFeels {
		f.HourlyThermalFeel[tf.Hour] = tf.ThermalFeel
	}

	for _, h := range f.Humidities {
		f.HourlyHumidity[h.Hour] = h.HumidityPercent
	}
}

type ParsedForecast struct {
	Elaborated      time.Time
	Location        string
	DateStr         string
	Hour            int
	SkyState        string
	Precipitation   float32
	POPPercent      int
	Temperature     int
	ThermalFeel     int
	HumidityPercent int
}

func (pf *ParsedForecast) String() string {
	fdstr := ""
	if feeldiff := pf.Temperature - pf.ThermalFeel; feeldiff != 0 {
		fdstr = fmt.Sprintf("(%+d)", feeldiff)
	}

	tempstr := ""
	if pf.ThermalFeel <= 17 || pf.ThermalFeel >= 22 {
		tempstr = strings.TrimRight(fmt.Sprintf(" %d°C %s", pf.Temperature, fdstr), " ")
	}

	popstr := ""
	if pf.POPPercent > 10 {
		popstr = fmt.Sprintf(" %d%%", pf.POPPercent)
	}

	return statusFonts[strings.TrimSuffix(pf.SkyState, "n")] + popstr + tempstr
}
