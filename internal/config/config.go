package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port               string
	GooglePlacesAPIKey string
	Environment        string
}

func Load() (Config, error) {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port := ":" + portStr

	apikey := os.Getenv("GOOGLE_PLACES_API_KEY")
	if apikey == "" {
		return Config{}, fmt.Errorf("GOOGLE_PLACES_API_KEY required")
	}
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	return Config{
		Port:               port,
		GooglePlacesAPIKey: apikey,
		Environment:        env,
	}, nil

}
