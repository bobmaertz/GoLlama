package tools

import (
	"encoding/json"

	"github.com/bobmaertz/ollama-agent/pkg/tools/weather"
)

// Available is a map of available tools and their corresponding functions.
var Available = map[string]func(json.RawMessage) (string, error){
	// TODO: Refactor this to allow for the tools to describe themselves so
	// they can be injected into the api calls
	"get_current_weather": weather.GetCurrentWeather,
}
