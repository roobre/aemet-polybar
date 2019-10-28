# ⛅ aemet-polybar 

Simple module to show weather forecast from AEMET in polybar. Supports showing forecast for next `n` hours of a given city.

![DeepinScreenshot_select-area_20191028113359](https://user-images.githubusercontent.com/969721/67672215-c9963680-f977-11e9-8967-6b7e4ef15e18.png)

` 21°C   20°C`

## Configuration

```ini
[module/aemet]
type = custom/script
exec = aemet
interval = 60
```

### Dependencies

* [**Font Awesome**](https://fontawesome.com/): Icons are shamefully hardcoded to Fontawesome's in the [fontawesome.go](https://github.com/roobre/aemet-polybar/blob/master/fontawesome.go) file. This sucks hard but works. Mode modular support for other icon fonts, Font Awesome Pro, or emojis would be great. PRs welcome and stuff.

EOF. Go binaries are self-contained and need no external deps.
