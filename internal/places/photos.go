package places

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jdanker/savor-api/internal/models"
)

// Same max size the iOS FetchPhotoRequest asked for.
const (
	photoMaxWidthPx  = 800
	photoMaxHeightPx = 600
)

type googlePhoto struct {
	Name               string `json:"name"` // "places/{id}/photos/{ref}"
	WidthPx            int    `json:"widthPx"`
	HeightPx           int    `json:"heightPx"`
	AuthorAttributions []struct {
		DisplayName string `json:"displayName"`
		URI         string `json:"uri"`
		PhotoURI    string `json:"photoUri"`
	} `json:"authorAttributions"`
}

// photoMedia is the media endpoint's response when skipHttpRedirect=true:
// the CDN URI as JSON instead of a 302 to the image bytes.
type photoMedia struct {
	PhotoURI string `json:"photoUri"`
}

// Photos returns up to max photo URIs for a place. Two-step: a photos-only
// details call for the photo names, then one media call per name. iOS downloads
// the actual bytes straight from the returned CDN URIs — image payloads never
// transit this service.
func (c *Client) Photos(ctx context.Context, placeID string, max int) ([]models.Photo, error) {
	var gp struct {
		Photos []googlePhoto `json:"photos"`
	}
	if err := c.do(ctx, "GET", "/places/"+url.PathEscape(placeID), "id,photos", nil, &gp); err != nil {
		return nil, err
	}

	if len(gp.Photos) > max {
		gp.Photos = gp.Photos[:max]
	}

	photos := make([]models.Photo, 0, len(gp.Photos))
	for _, g := range gp.Photos {
		mediaPath := fmt.Sprintf("/%s/media?maxWidthPx=%d&maxHeightPx=%d&skipHttpRedirect=true",
			g.Name, photoMaxWidthPx, photoMaxHeightPx)

		var media photoMedia
		if err := c.do(ctx, "GET", mediaPath, "", nil, &media); err != nil {
			// One bad photo shouldn't sink the set — return what resolved.
			continue
		}

		p := models.Photo{
			URI:      media.PhotoURI,
			WidthPx:  g.WidthPx,
			HeightPx: g.HeightPx,
		}
		for _, a := range g.AuthorAttributions {
			p.Attributions = append(p.Attributions, models.PhotoAttribution{
				DisplayName: a.DisplayName,
				URI:         a.URI,
				PhotoURI:    a.PhotoURI,
			})
		}
		photos = append(photos, p)
	}
	return photos, nil
}
