package aemet

import "strings"

func statusIcon(status string) string {
	if icon := statusFonts[status]; icon != "" {
		return icon
	}

	return statusFonts[strings.TrimSuffix(status, "n")]
}

var statusFonts = map[string]string{
	"11":  "\uf185", // Despejado
	"11n": "\uf186", // Despejado
	"12":  "\ueef0",
	"12n": "\ueeef",
	"13":  "\ueef0",
	"13n": "\ueeef",
	"14":  "\uf0c2",
	"15":  "\uf0c2",
	"16":  "\uf0c2",
	"17":  "\uf0c2",
	"23":  "\uef1c",
	"24":  "\uef1c",
	"25":  "\uef1c",
	"26":  "\uef1c",
	"27":  "\uef1c",
	"33":  "\ue201",
	"34":  "\ue201",
	"35":  "\ue201",
	"36":  "\ue201",
	"43":  "\uf0c2",
	"44":  "\uf0c2",
	"45":  "\uf0c2",
	"46":  "\uf0c2",
	"51":  "\uef2c",
	"52":  "\uef2c",
	"53":  "\uef2c",
	"54":  "\uef2c",
	"61":  "\uef2c",
	"62":  "\uef2c",
	"63":  "\uef2c",
	"64":  "\uef2c",
	"71":  "\uf0c2",
	"72":  "\uf0c2",
	"73":  "\uf0c2",
	"74":  "\uf0c2",
}
