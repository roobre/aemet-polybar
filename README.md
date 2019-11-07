# ⛅ aemet-polybar 

Simple module to show weather forecast from AEMET in polybar. Supports showing forecast for next `n` hours of a given city.

![DeepinScreenshot_select-area_20191028113359](https://user-images.githubusercontent.com/969721/67672215-c9963680-f977-11e9-8967-6b7e4ef15e18.png)

## Usage

aemet-polybar supports two usage modes. The first one will output the forecast for the next `n` hours. This is the default mode, with `n = 2`

```bash
$ aemet -n 2
 21°C   20°C
```

In addition, aemet supports passing multiple labels and times, to show the forecast for those concrete times.
Currently the label specification is a bit custom-ish and hacky, and obeys the following syntax:

`aemet -l <LABEL>:<TIME>[,<LABEL>:<TIME>]...`

For example:

| Label      | Output |
|------------|--------|
|`🏠:18`     |`🏠:  21°C`|
|`9h:9,🏠:18`|`9h:  18°C  🏠:  21°C`|

Icons and/or emoji are usually good labels, although any string can be used as `<LABEL>` as long as it does not contain a colon (`:`) or a comma(`,`). `<TIME>` must be an integer representing the hour you want to see the forecast for. Multiple pairs of `<LABEL>:<TIME>` must be separated with a comma (`,`).

`aemet` will play smart and automatically hide forecasts for times that are in the past.

## Output

The full output format is:

`{icon} ({POP}%) {temperature} (+{thermalfeel})`

however, to avoid cluttering your bar with non-relevant information, `aemet` will hide values that it considers not remarkable. More precisely:

* Probability of precipitation (POP) will only be show when it is greater than 10%
* Thermal Feel will only be shown when it is different than the actual temperature
* Temperature will only be shown when it is outside of the (17, 22) range.

Neither these defaults nor the output format are configurable at the moment. PRs implementing them using flags and `text/template` are welcome.

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
