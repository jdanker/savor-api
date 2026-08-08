package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"

	"github.com/jdanker/savor-api/internal/places"
)

// PlacesHandler holds the one dependency the place endpoints need. Methods on a
// struct (rather than package-level funcs) so the client is injected once at
// wire-up instead of reached for globally.
type PlacesHandler struct {
	client *places.Client
}

func NewPlacesHandler(client *places.Client) *PlacesHandler {
	return &PlacesHandler{client: client}
}

// uuidRe validates searchSessionID before it's forwarded anywhere. Cheap gate:
// malformed input is rejected before any billed Google call.
var uuidRe = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

type autocompleteBody struct {
	Query           string `json:"query"`
	SearchSessionID string `json:"searchSessionID"`
}

// Autocomplete handles POST /places/autocomplete.
func (h *PlacesHandler) Autocomplete(w http.ResponseWriter, r *http.Request) {
	var body autocompleteBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if body.Query == "" {
		writeError(w, http.StatusBadRequest, "query is required")
		return
	}
	if !uuidRe.MatchString(body.SearchSessionID) {
		writeError(w, http.StatusBadRequest, "searchSessionID must be a UUID")
		return
	}

	// Structured so session tokens can be grepped later for
	// autocomplete-to-details ratio analysis (billing).
	slog.Info("autocomplete", "searchSessionID", body.SearchSessionID)

	suggestions, err := h.client.Autocomplete(r.Context(), body.Query, body.SearchSessionID)
	if err != nil {
		writeUpstreamError(w, "autocomplete", err)
		return
	}
	writeJSON(w, http.StatusOK, suggestions)
}

// Details handles GET /places/{id}?tier=save|enrich|coordinate.
func (h *PlacesHandler) Details(w http.ResponseWriter, r *http.Request) {
	placeID := r.PathValue("id")

	tier, err := places.ParseTier(r.URL.Query().Get("tier"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Session policy is server-side: only the save tier closes the autocomplete
	// billing session; other tiers never forward the token even if sent.
	sessionToken := ""
	if tier == places.TierSave {
		sessionToken = r.URL.Query().Get("searchSessionID")
		if sessionToken != "" && !uuidRe.MatchString(sessionToken) {
			writeError(w, http.StatusBadRequest, "searchSessionID must be a UUID")
			return
		}
		slog.Info("details save", "searchSessionID", sessionToken, "placeID", placeID)
	}

	place, err := h.client.Details(r.Context(), placeID, tier, sessionToken)
	if err != nil {
		writeUpstreamError(w, "details", err)
		return
	}
	writeJSON(w, http.StatusOK, place)
}

// Photos handles GET /places/{id}/photos?max=3.
func (h *PlacesHandler) Photos(w http.ResponseWriter, r *http.Request) {
	placeID := r.PathValue("id")

	max := 3
	if s := r.URL.Query().Get("max"); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 || n > 10 {
			writeError(w, http.StatusBadRequest, "max must be an integer between 1 and 10")
			return
		}
		max = n
	}

	photos, err := h.client.Photos(r.Context(), placeID, max)
	if err != nil {
		writeUpstreamError(w, "photos", err)
		return
	}
	writeJSON(w, http.StatusOK, photos)
}

// writeJSON is the single success-encoding path for all handlers.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeError returns {"error": ...} — one error shape for iOS to decode.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// writeUpstreamError logs the real Google error but returns a generic 502 —
// upstream error bodies could leak key/quota details to the client.
func writeUpstreamError(w http.ResponseWriter, op string, err error) {
	slog.Error(op+" upstream failure", "error", err)
	writeError(w, http.StatusBadGateway, "upstream place lookup failed")
}
