package models

type Place struct {
	PlaceID     string
	Name        string
	Rating      float64
	PriceLevel  string
	Types       []string
	Coordinates struct {
		Long float64
		Lat  float64
	}
	GenerativeSummary string
	URL               string
	ReviewSummary     string
	PhotoReference    string
}
