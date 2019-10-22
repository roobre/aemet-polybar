package main

import (
	"fmt"
	"strings"
	"time"
)

const dateFormat = "2006-01-02"

type Location struct {
	Name          string `xml:"nombre"`
	Region        string `xml:"provincia"`
	ElaboratedStr string `xml:"elaborado"`

	Forecasts      []DailyForecast `xml:"prediccion>dia" json:"-"`
	DailyForecasts map[string]*DailyForecast
}

type DailyForecast struct {
	DateStr string `xml:"fecha,attr"`

	SkyStates []struct {
		Hour        int    `xml:"periodo,attr"`
		Description string `xml:"descripcion,attr"`
		State       string `xml:",chardata"`
	} `xml:"estado_cielo" json:"-"`
	HourlySkyState map[int]string

	Precipitations []struct {
		Hour   int     `xml:"periodo,attr"`
		Amount float32 `xml:",chardata"`
	} `xml:"precipitations" json:"-"`
	HourlyPrecipitation map[int]float32

	POPs []struct {
		Hour       int `xml:"periodo,attr"`
		POPPercent int `xml:",chardata"`
	} `xml:"prob_precipitacion" json:"-"`
	HourlyPOP map[int]int

	Temperatures []struct {
		Hour        int `xml:"periodo,attr"`
		Temperature int `xml:",chardata"`
	} `xml:"temperatura" json:"-"`
	HourlyTemperature map[int]int

	ThermalFeels []struct {
		Hour        int `xml:"periodo,attr"`
		ThermalFeel int `xml:",chardata"`
	} `xml:"sens_termica" json:"-"`
	HourlyThermalFeel map[int]int

	Humidities []struct {
		Hour            int `xml:"periodo,attr"`
		HumidityPercent int `xml:",chardata"`
	} `xml:"humedad_relativa" json:"-"`
	HourlyHumidity map[int]int
}

func (l *Location) NextHours(n int) []ParsedForecast {
	forecasts := make([]ParsedForecast, 0, n)

	for i := 1; i <= n; i++ {
		day := time.Now().Truncate(24 * time.Hour).Format(dateFormat)
		hour := time.Now().Add(-15 * time.Minute).Truncate(time.Hour).Add(time.Duration(i) * time.Hour)

		f := l.DailyForecasts[day]
		forecasts = append(forecasts, ParsedForecast{
			//Elaborated:      l.ElaboratedStr,
			Location: l.Name,
			DateStr:  f.DateStr,
			Hour:     hour.Hour(),

			SkyState:        f.HourlySkyState[hour.Hour()],
			Precipitation:   f.HourlyPrecipitation[hour.Hour()],
			POPPercent:      f.HourlyPOP[hour.Hour()],
			Temperature:     f.HourlyTemperature[hour.Hour()],
			ThermalFeel:     f.HourlyThermalFeel[hour.Hour()],
			HumidityPercent: f.HourlyHumidity[hour.Hour()],
		})
	}

	return forecasts
}

func (l *Location) Parse() {
	l.DailyForecasts = map[string]*DailyForecast{}

	for i := range l.Forecasts {
		f := &l.Forecasts[i]
		f.Parse()
		l.DailyForecasts[f.DateStr] = f
	}
}

func (f *DailyForecast) Parse() {
	f.HourlySkyState = make(map[int]string, len(f.SkyStates))
	for _, ss := range f.SkyStates {
		f.HourlySkyState[ss.Hour] = ss.State
	}

	f.HourlyPrecipitation = make(map[int]float32, len(f.Precipitations))
	for _, p := range f.Precipitations {
		f.HourlyPrecipitation[p.Hour] = p.Amount
	}

	f.HourlyPOP = make(map[int]int, len(f.POPs))
	for _, pop := range f.POPs {
		f.HourlyPOP[pop.Hour] = pop.POPPercent
	}

	f.HourlyTemperature = make(map[int]int, len(f.Temperatures))
	for _, t := range f.Temperatures {
		f.HourlyTemperature[t.Hour] = t.Temperature
	}

	f.HourlyThermalFeel = make(map[int]int, len(f.ThermalFeels))
	for _, tf := range f.ThermalFeels {
		f.HourlyThermalFeel[tf.Hour] = tf.ThermalFeel
	}

	f.HourlyHumidity = make(map[int]int, len(f.Humidities))
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
	return fmt.Sprintf("%s %d%s°C", statusFonts[strings.TrimSuffix(pf.SkyState, "n")], pf.Temperature, fdstr)
}
