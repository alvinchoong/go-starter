package router

import (
	"encoding/json"
	"net/http"

	"github.com/alvinchoong/go-httphandler"
	"github.com/alvinchoong/go-httphandler/jsonresp"
	"github.com/go-chi/chi/v5"
)

// userHandler implements a proxy to an external user API service
type userHandler struct {
	client   *http.Client // HTTP client for external API requests
	endpoint string       // Base URL of the external API
}

// NewUserHandler creates a new user handler with configured HTTP client and API endpoint
func NewUserHandler(client *http.Client, endpoint string) *userHandler {
	return &userHandler{
		client:   client,
		endpoint: endpoint,
	}
}

// Mount registers user-related routes
func (h *userHandler) Mount(r chi.Router) {
	r.Get("/api/v1/users", httphandler.Handle(h.Get))
}

// User represents the structure of user data from the external API
type User struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Address  Address `json:"address"`
	Phone    string  `json:"phone"`
	Website  string  `json:"website"`
}

// Address represents the structure of a user's physical address
type Address struct {
	Street  string `json:"street"`
	Suite   string `json:"suite"`
	City    string `json:"city"`
	Zipcode string `json:"zipcode"`
}

// Get proxies the request to an external API and returns the user data
func (h *userHandler) Get(r *http.Request) httphandler.Responder {
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

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return jsonresp.Error(err, "failed to decode response", http.StatusInternalServerError)
	}

	return jsonresp.Success(&users)
}
