package places

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jdanker/savor-api/internal/models"
)

// Tier names intent, not fields: the field masks (and their billing SKUs) stay
// server-side so iOS never learns what a field mask is.
type Tier string

const (
	TierSave       Tier = "save"       // cheap on-add fetch
	TierEnrich     Tier = "enrich"     // expensive on-demand fetch (Gemini summaries)
	TierCoordinate Tier = "coordinate" // backfill for pre-coordinate saves
)

// ParseTier rejects unknown tiers before any billed call is made.
func ParseTier(s string) (Tier, error) {
	switch t := Tier(s); t {
	case TierSave, TierEnrich, TierCoordinate:
		return t, nil
	default:
		return "", fmt.Errorf("unknown tier %q (want save, enrich, or coordinate)", s)
	}
}

// Field masks per tier — the same field sets the iOS PlaceProperty lists requested.
var tierFieldMasks = map[Tier]string{
	TierSave:       "id,displayName,rating,priceLevel,types,location,editorialSummary,websiteUri",
	TierEnrich:     "id,websiteUri,generativeSummary,reviewSummary",
	TierCoordinate: "id,location",
}

// googlePlace covers the union of all tier responses; unrequested fields
// simply stay zero. Google's shape never leaves this package.
type googlePlace struct {
	ID          string        `json:"id"`
	DisplayName localizedText `json:"displayName"`
	Rating      float64       `json:"rating"`
	PriceLevel  string        `json:"priceLevel"`
	Types       []string      `json:"types"`
	Location    *struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	EditorialSummary  localizedText `json:"editorialSummary"`
	WebsiteURI        string        `json:"websiteUri"`
	GenerativeSummary struct {
		Overview localizedText `json:"overview"`
	} `json:"generativeSummary"`
	ReviewSummary struct {
		Text localizedText `json:"text"`
	} `json:"reviewSummary"`
}

// Details fetches a place at the given tier. sessionToken should be empty except
// on tier=save — passing it closes the autocomplete billing session at Google.
func (c *Client) Details(ctx context.Context, placeID string, tier Tier, sessionToken string) (models.Place, error) {
	mask, ok := tierFieldMasks[tier]
	if !ok {
		return models.Place{}, fmt.Errorf("unknown tier %q", tier)
	}

	path := "/places/" + url.PathEscape(placeID)
	if sessionToken != "" {
		path += "?sessionToken=" + url.QueryEscape(sessionToken)
	}

	var gp googlePlace
	if err := c.do(ctx, "GET", path, mask, nil, &gp); err != nil {
		return models.Place{}, err
	}
	return translatePlace(gp), nil
}

// translatePlace converts Google's vocabulary to ours at the package boundary.
func translatePlace(gp googlePlace) models.Place {
	p := models.Place{
		PlaceID:           gp.ID,
		Name:              gp.DisplayName.Text,
		Rating:            gp.Rating,
		PriceLevel:        priceLevelToInt(gp.PriceLevel),
		Types:             gp.Types,
		EditorialSummary:  gp.EditorialSummary.Text,
		GenerativeSummary: gp.GenerativeSummary.Overview.Text,
		ReviewSummary:     gp.ReviewSummary.Text.Text,
		WebsiteURL:        gp.WebsiteURI,
	}
	if gp.Location != nil {
		p.Coordinates = &models.Coordinates{
			Lat:  gp.Location.Latitude,
			Long: gp.Location.Longitude,
		}
	}
	return p
}

// priceLevelToInt maps Google's enum strings to 1-4 ($-$$$$). Free, unspecified,
// and any future enum value all collapse to nil — unknown is never an error,
// matching the iOS @unknown default behavior.
func priceLevelToInt(s string) *int {
	var n int
	switch s {
	case "PRICE_LEVEL_INEXPENSIVE":
		n = 1
	case "PRICE_LEVEL_MODERATE":
		n = 2
	case "PRICE_LEVEL_EXPENSIVE":
		n = 3
	case "PRICE_LEVEL_VERY_EXPENSIVE":
		n = 4
	default:
		return nil
	}
	return &n
}
