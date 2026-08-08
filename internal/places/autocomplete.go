package places

import (
	"context"

	"github.com/jdanker/savor-api/internal/models"
)

// Same type filter the iOS AutocompleteFilter used against the SDK —
// parity target is what the app does today, nothing broader.
var autocompleteTypes = []string{"restaurant", "cafe", "bar", "bakery"}

// Field mask trimmed to exactly what models.Suggestion needs; matches (offsets)
// are deliberately excluded — highlighting was dropped for Phase 1 (decisions.md).
const autocompleteFieldMask = "suggestions.placePrediction.placeId," +
	"suggestions.placePrediction.text.text," +
	"suggestions.placePrediction.structuredFormat.mainText.text"

type autocompleteRequest struct {
	Input                string   `json:"input"`
	SessionToken         string   `json:"sessionToken"`
	IncludedPrimaryTypes []string `json:"includedPrimaryTypes"`
}

type autocompleteResponse struct {
	Suggestions []struct {
		PlacePrediction struct {
			PlaceID          string        `json:"placeId"`
			Text             localizedText `json:"text"`
			StructuredFormat struct {
				MainText localizedText `json:"mainText"`
			} `json:"structuredFormat"`
		} `json:"placePrediction"`
	} `json:"suggestions"`
}

// Autocomplete returns restaurant suggestions for a query. sessionToken groups
// keystrokes with the eventual details call for Google's session billing.
func (c *Client) Autocomplete(ctx context.Context, query, sessionToken string) ([]models.Suggestion, error) {
	body := autocompleteRequest{
		Input:                query,
		SessionToken:         sessionToken,
		IncludedPrimaryTypes: autocompleteTypes,
	}

	var resp autocompleteResponse
	if err := c.do(ctx, "POST", "/places:autocomplete", autocompleteFieldMask, body, &resp); err != nil {
		return nil, err
	}

	// Always return a non-nil slice so an empty result encodes as [] not null.
	suggestions := make([]models.Suggestion, 0, len(resp.Suggestions))
	for _, s := range resp.Suggestions {
		p := s.PlacePrediction
		if p.PlaceID == "" {
			continue // query predictions (no place) — iOS skips these too
		}
		suggestions = append(suggestions, models.Suggestion{
			PlaceID:     p.PlaceID,
			PrimaryText: p.StructuredFormat.MainText.Text,
			FullText:    p.Text.Text,
		})
	}
	return suggestions, nil
}
