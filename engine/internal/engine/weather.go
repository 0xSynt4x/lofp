package engine

// WeatherNames maps weather state IDs to display names (from GM Manual).
var WeatherNames = map[int]string{
	0:  "Sunny",
	1:  "Partly Cloudy",
	2:  "Overcast",
	3:  "Light Rain",
	4:  "Moderate Rain",
	5:  "Heavy Rain",
	6:  "Thunderstorm",
	7:  "Gale",
	8:  "Hurricane",
	9:  "Hail",
	10: "Sleet",
	11: "Snow Flurries",
	12: "Moderate Snow",
	13: "Heavy Snow",
	14: "Blizzard",
}

// GetWeatherDesc returns a weather description for a given region.
func (e *GameEngine) GetWeatherDesc(region int) string {
	if e.RegionWeather == nil {
		return ""
	}
	state, ok := e.RegionWeather[region]
	if !ok {
		state = 0
	}
	if name, ok := WeatherNames[state]; ok {
		return name
	}
	return "Clear"
}

// GetRoomWeather returns a weather line for an outdoor room, or "" for indoor.
func (e *GameEngine) GetRoomWeather(roomNum int) string {
	room := e.rooms[roomNum]
	if room == nil {
		return ""
	}
	if !isOutdoorTerrain(room.Terrain) {
		return ""
	}
	region := room.Region
	desc := e.GetWeatherDesc(region)
	if desc == "" || desc == "Sunny" || desc == "Clear" {
		return ""
	}
	return "The weather is " + desc + "."
}

// isOutdoorTerrain returns true if the terrain type is outdoors.
func isOutdoorTerrain(terrain string) bool {
	switch terrain {
	case "FOREST", "MOUNTAIN", "PLAIN", "SWAMP", "JUNGLE",
		"WASTE", "OUTDOOR_OTHER", "OUTDOOR_FLOOR", "AERIAL":
		return true
	}
	return false
}
