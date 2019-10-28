# ⛅ aemet-polybar 

Simple module to show weather forecast from AEMET in polybar. Supports showing forecast for next `n` hours of a given city.

![DeepinScreenshot_select-area_20191028113359](https://user-images.githubusercontent.com/969721/67672215-c9963680-f977-11e9-8967-6b7e4ef15e18.png)

```ini
[module/aemet]
type = custom/script
exec = aemet
interval = 60
```
