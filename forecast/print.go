package forecast

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/jessfraz/weather/geocode"
	"github.com/jessfraz/weather/icons"
	"github.com/mitchellh/colorstring"
)

// UnitMeasures are the location specific terms for weather data.
type UnitMeasures struct {
	Degrees       string
	Speed         string
	Length        string
	Precipitation string
	LongDate      string
	Hour          string
}

var (
	// UnitFormats describe each regions UnitMeasures.
	UnitFormats = map[string]UnitMeasures{
		"us": {
			Degrees:       "°F",
			Speed:         "mph",
			Length:        "miles",
			Precipitation: "in/hr",
			LongDate:      "January 2 at 3:04pm MST",
			Hour:          "3:04pm MST",
		},
		"si": {
			Degrees:       "°C",
			Speed:         "m/s",
			Length:        "kilometers",
			Precipitation: "mm/h",
			LongDate:      "2006-01-02 15:04:05 EET",
			Hour:          "15:04 EET",
		},
		"ca": {
			Degrees:       "°C",
			Speed:         "km/h",
			Length:        "kilometers",
			Precipitation: "mm/h",
			LongDate:      "January 2 at 3:04pm MST",
			Hour:          "3:04pm MST",
		},
		// deprecated, use "uk2" in stead
		"uk": {
			Degrees:       "°C",
			Speed:         "mph",
			Length:        "kilometers",
			Precipitation: "mm/h",
			LongDate:      "January 2 at 15:04 MST",
			Hour:          "15:04 MST",
		},
		"uk2": {
			Degrees:       "°C",
			Speed:         "mph",
			Length:        "miles",
			Precipitation: "mm/h",
			LongDate:      "January 2 at 15:04 MST",
			Hour:          "15:04 MST",
		},
	}
	// Directions contain all the combinations of N,S,E,W
	Directions = []string{
		"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW",
	}
)

func epochFormat(seconds int64, unitsFormat UnitMeasures) string {
	epochTime := time.Unix(0, seconds*int64(time.Second))
	return epochTime.Format(unitsFormat.LongDate)
}

func epochFormatDate(seconds int64) string {
	epochTime := time.Unix(0, seconds*int64(time.Second))
	return epochTime.Format("January 2 (Monday)")
}

func epochFormatTime(seconds int64, unitsFormat UnitMeasures) string {
	epochTime := time.Unix(0, seconds*int64(time.Second))
	return epochTime.Format(unitsFormat.Hour)
}

func getIcon(iconStr string) (icon string, err error) {
	color := "blue"
	// steralize the icon string name
	iconStr = strings.Replace(strings.Replace(iconStr, "-", "", -1), "_", "", -1)

	switch iconStr {
	case "clear":
		icon = icons.Clear
	case "clearday":
		color = "yellow"
		icon = icons.Clearday
	case "clearnight":
		color = "light_yellow"
		icon = icons.Clearnight
	case "clouds":
		icon = icons.Clouds
	case "cloudy":
		icon = icons.Cloudy
	case "cloudsnight":
		color = "light_yellow"
		icon = icons.Cloudsnight
	case "fog":
		icon = icons.Fog
	case "haze":
		icon = icons.Haze
	case "hazenight":
		color = "light_yellow"
		icon = icons.Hazenight
	case "partlycloudyday":
		color = "yellow"
		icon = icons.Partlycloudyday
	case "partlycloudynight":
		color = "light_yellow"
		icon = icons.Partlycloudynight
	case "rain":
		icon = icons.Rain
	case "sleet":
		icon = icons.Sleet
	case "snow":
		color = "white"
		icon = icons.Snow
	case "thunderstorm":
		color = "black"
		icon = icons.Thunderstorm
	case "tornado":
		color = "black"
		icon = icons.Tornado
	case "wind":
		color = "black"
		icon = icons.Wind
	}

	return colorstring.Color("[" + color + "]" + icon), nil
}

func getBearingDetails(degrees float64) string {
	index := int(math.Mod((degrees+11.25)/22.5, 16))
	return Directions[index]
}

func printCommon(weather Weather, unitsFormat UnitMeasures) error {
	if weather.Humidity > 0 {
		humidity := colorstring.Color(fmt.Sprintf("[white]%v%s", weather.Humidity*100, "%"))
		if weather.Humidity > 0.20 {
			fmt.Printf("  Ick! The humidity is %s\n", humidity)
		} else {
			fmt.Printf("  The humidity is %s\n", humidity)
		}
	}

	if weather.PrecipIntensity > 0 {
		precInt := colorstring.Color(fmt.Sprintf("[white]%v %s", weather.PrecipIntensity, unitsFormat.Precipitation))
		fmt.Printf("  The precipitation intensity of %s is %s\n", colorstring.Color("[white]"+weather.PrecipType), precInt)
	}

	if weather.PrecipProbability > 0 {
		prec := colorstring.Color(fmt.Sprintf("[white]%v%s", weather.PrecipProbability*100, "%"))
		fmt.Printf("  The precipitation probability is %s\n", prec)
	}

	if weather.NearestStormDistance > 0 {
		dist := colorstring.Color(fmt.Sprintf("[white]%v %s %v", weather.NearestStormDistance, unitsFormat.Length, getBearingDetails(weather.NearestStormBearing)))
		fmt.Printf("  The nearest storm is %s away\n", dist)
	}

	if weather.WindSpeed > 0 {
		wind := colorstring.Color(fmt.Sprintf("[white]%v %s %v", weather.WindSpeed, unitsFormat.Speed, getBearingDetails(weather.WindBearing)))
		fmt.Printf("  The wind speed is %s\n", wind)
	}

	if weather.CloudCover > 0 {
		cloudCover := colorstring.Color(fmt.Sprintf("[white]%v%s", weather.CloudCover*100, "%"))
		fmt.Printf("  The cloud coverage is %s\n", cloudCover)
	}

	if weather.Visibility < 10 {
		visibility := colorstring.Color(fmt.Sprintf("[white]%v %s", weather.Visibility, unitsFormat.Length))
		fmt.Printf("  The visibility is %s\n", visibility)
	}

	if weather.Pressure > 0 {
		pressure := colorstring.Color(fmt.Sprintf("[white]%v %s", weather.Pressure, "mbar"))
		fmt.Printf("  The pressure is %s\n\n", pressure)
	}

	return nil
}

// PrintCurrent pretty prints the current forecast data.
func PrintCurrent(forecast Forecast, geolocation geocode.Geocode, ignoreAlerts bool, hideIcon bool) error {
	unitsFormat := UnitFormats[forecast.Flags.Units]

	if !hideIcon {
		icon, err := getIcon(forecast.Currently.Icon)
		if err != nil {
			return err
		}

		fmt.Println(icon)
	}

	location := colorstring.Color(fmt.Sprintf("[green]%s in %s", geolocation.City, geolocation.Region))
	fmt.Printf("\nCurrent weather is %s in %s for %s\n", colorstring.Color("[cyan]"+forecast.Currently.Summary), location, colorstring.Color("[cyan]"+epochFormat(forecast.Currently.Time, unitsFormat)))

	temp := colorstring.Color(fmt.Sprintf("[magenta]%v%s", forecast.Currently.Temperature, unitsFormat.Degrees))
	feelslike := colorstring.Color(fmt.Sprintf("[magenta]%v%s", forecast.Currently.ApparentTemperature, unitsFormat.Degrees))
	if temp == feelslike {
		fmt.Printf("The temperature is %s\n\n", temp)
	} else {
		fmt.Printf("The temperature is %s, but it feels like %s\n\n", temp, feelslike)
	}

	if !ignoreAlerts {
		for _, alert := range forecast.Alerts {
			if alert.Title != "" {
				fmt.Println(colorstring.Color("[red]" + alert.Title))
			}
			if alert.Description != "" {
				fmt.Print(colorstring.Color("[red]" + alert.Description))
			}
			fmt.Println("\t\t\t" + colorstring.Color("[red]Created: "+epochFormat(alert.Time, unitsFormat)))
			fmt.Println("\t\t\t" + colorstring.Color("[red]Expires: "+epochFormat(alert.Expires, unitsFormat)) + "\n")
		}
	}

	return printCommon(forecast.Currently, unitsFormat)
}

// PrintDaily pretty prints the daily forecast data.
func PrintDaily(forecast Forecast, days int) error {
	unitsFormat := UnitFormats[forecast.Flags.Units]

	// Ignore the current day as it's printed before
	for index, daily := range forecast.Daily.Data[1:] {
		// only do the amount of days they request
		if index == days {
			break
		}

		fmt.Println(colorstring.Color("[magenta]" + epochFormatDate(daily.Time)))

		tempMax := colorstring.Color(fmt.Sprintf("[blue]%v%s", daily.TemperatureMax, unitsFormat.Degrees))
		tempMin := colorstring.Color(fmt.Sprintf("[blue]%v%s", daily.TemperatureMin, unitsFormat.Degrees))
		feelsLikeMax := colorstring.Color(fmt.Sprintf("[cyan]%v%s", daily.ApparentTemperatureMax, unitsFormat.Degrees))
		feelsLikeMin := colorstring.Color(fmt.Sprintf("[cyan]%v%s", daily.ApparentTemperatureMin, unitsFormat.Degrees))
		fmt.Printf("The temperature high is %s, feels like %s around %s,\n", tempMax, feelsLikeMax, epochFormatTime(daily.TemperatureMaxTime, unitsFormat))
		fmt.Printf("and low is %s, feels like %s around %s\n\n", tempMin, feelsLikeMin, epochFormatTime(daily.TemperatureMinTime, unitsFormat))

		printCommon(daily, unitsFormat)
	}

	return nil
}
