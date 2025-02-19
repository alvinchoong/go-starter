package router

import (
	"encoding/json"
	"net/http"

	"github.com/alvinchoong/go-httphandler"
	"github.com/alvinchoong/go-httphandler/jsonresp"
)

// quoteHandler implements a proxy to an external user API service
type quoteHandler struct {
	client   *http.Client // HTTP client for external API requests
	endpoint string       // Base URL of the external API
}

// NewQuoteHandler creates a new user handler with configured HTTP client and API endpoint
func NewQuoteHandler(client *http.Client, endpoint string) *quoteHandler {
	return &quoteHandler{
		client:   client,
		endpoint: endpoint,
	}
}

// QuoteResponse represents the structure of quote data from the external API
type QuoteResponse struct {
	ID     int    `json:"id"`
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

// Get proxies the request to an external API and returns the quote data
func (h *quoteHandler) Get(r *http.Request) httphandler.Responder {
	ctx := r.Context()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, h.endpoint, nil)
	if err != nil {
		return jsonresp.Error(err, "failed to create request", http.StatusInternalServerError)
	}

	resp, err := h.client.Do(request)
	if err != nil {
		return jsonresp.Error(err, "failed to make request", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return jsonresp.Error(nil, "failed to fetch data from external API", resp.StatusCode)
	}

	var quote QuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
		return jsonresp.Error(err, "failed to decode response", http.StatusInternalServerError)
	}

	return jsonresp.Success(&quote)
}
