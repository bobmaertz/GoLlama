package weather

import (
	"encoding/json"
	"fmt"
)

type arguments struct {
	Location string
}

// GetCurrentWeather is a stub function that simulates getting the current weather.
func GetCurrentWeather(args json.RawMessage) (string, error) {

	// This is a cool way to handle different types of
	// arguments. It allows you to define a struct
	// Credit to Thorsten Ball's post https://ampcode.com/how-to-build-an-agent
	// for this idea
	details := arguments{}
	err := json.Unmarshal(args, &details)
	if err != nil {
		return "", err
	}

	fmt.Println(details)
	return "sunny", nil
}
