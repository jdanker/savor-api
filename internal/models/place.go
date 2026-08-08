package models

// Place is the provider-neutral place shape returned by GET /places/{id}.
// Fields are tier-scoped: a coordinate-tier response only populates Coordinates,
// so most fields use omitempty rather than forcing iOS to decode zero values.
type Place struct {
	PlaceID string  `json:"placeID"`
	Name    string  `json:"name,omitempty"`
	Rating  float64 `json:"rating,omitempty"`
	// *int, not int: 0 would be indistinguishable from "Google didn't say".
	// nil = unknown/free; 1-4 maps to $-$$$$ like the iOS PriceLevel conversion.
	PriceLevel        *int         `json:"priceLevel,omitempty"`
	Types             []string     `json:"types,omitempty"`
	Coordinates       *Coordinates `json:"coordinates,omitempty"`
	EditorialSummary  string       `json:"editorialSummary,omitempty"`
	GenerativeSummary string       `json:"generativeSummary,omitempty"`
	ReviewSummary     string       `json:"reviewSummary,omitempty"`
	WebsiteURL        string       `json:"websiteURL,omitempty"`
}

// Coordinates is a pointer field on Place so "not requested" is nil, not (0, 0) —
// the iOS side validates coordinates instead of storing the null-island sentinel.
type Coordinates struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

// Suggestion is one autocomplete result: POST /places/autocomplete → []Suggestion.
type Suggestion struct {
	PlaceID     string `json:"placeID"`
	PrimaryText string `json:"primaryText"`
	FullText    string `json:"fullText"`
}

// Photo is one resolved photo: GET /places/{id}/photos → []Photo.
// Attributions ride along because Google's terms require displaying them.
type Photo struct {
	URI          string             `json:"uri"`
	WidthPx      int                `json:"widthPx"`
	HeightPx     int                `json:"heightPx"`
	Attributions []PhotoAttribution `json:"attributions,omitempty"`
}

// PhotoAttribution credits the photo author (Google terms require display).
type PhotoAttribution struct {
	DisplayName string `json:"displayName"`
	URI         string `json:"uri,omitempty"`
	PhotoURI    string `json:"photoURI,omitempty"`
}
