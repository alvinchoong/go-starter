package router

import (
	"errors"
	"net/http"

	"go-starter/internal/models"

	"github.com/alvinchoong/go-httphandler"
	"github.com/alvinchoong/go-httphandler/jsonresp"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// postHandler implements CRUD operations for blog posts
type postHandler struct {
	db      models.DBTX
	querier models.Querier
}

// NewPostHandler creates a new post handler with database connection and query interface
func NewPostHandler(db models.DBTX, q models.Querier) *postHandler {
	return &postHandler{
		db:      db,
		querier: q,
	}
}

// Mount registers all post-related routes
func (h *postHandler) Mount(r chi.Router) {
	r.Post("/api/v1/posts", httphandler.HandleWithInput(h.Create))
	r.Get("/api/v1/posts", httphandler.Handle(h.List))
	r.Get("/api/v1/posts/{id}", httphandler.Handle(h.Get))
	r.Put("/api/v1/posts/{id}", httphandler.HandleWithInput(h.Update))
	r.Delete("/api/v1/posts/{id}", httphandler.Handle(h.Delete))
}

// CreatePostParams defines the required fields for creating a new post
type CreatePostParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Creates a new blog post
func (h *postHandler) Create(r *http.Request, params CreatePostParams) httphandler.Responder {
	ctx := r.Context()

	post, err := h.querier.CreatePost(ctx, h.db, models.CreatePostParams{
		ID:          uuid.New(),
		Title:       params.Title,
		Description: &params.Description,
	})
	if err != nil {
		return jsonresp.InternalServerError(err)
	}

	return jsonresp.Success(&post)
}

// List retrieves a list of all blog posts
func (h *postHandler) List(r *http.Request) httphandler.Responder {
	ctx := r.Context()

	posts, err := h.querier.ListPosts(ctx, h.db)
	if err != nil {
		return jsonresp.InternalServerError(err)
	}
	return jsonresp.Success(&posts)
}

// Get retrieves a single blog post by ID
func (h *postHandler) Get(r *http.Request) httphandler.Responder {
	ctx := r.Context()

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return jsonresp.Error(err, "Invalid ID format", http.StatusBadRequest)
	}

	post, err := h.querier.GetPost(ctx, h.db, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return jsonresp.Error(err, "Post not found", http.StatusNotFound)
		}
		return jsonresp.InternalServerError(err)
	}

	return jsonresp.Success(&post)
}

// UpdatePostParams defines the required fields for updating a post
type UpdatePostParams struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

// Update ipdates an existing blog post by ID
func (h *postHandler) Update(r *http.Request, input UpdatePostParams) httphandler.Responder {
	ctx := r.Context()

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return jsonresp.Error(err, "Invalid ID format", http.StatusBadRequest)
	}

	post, err := h.querier.UpdatePost(ctx, h.db, models.UpdatePostParams{
		ID:          id,
		Title:       input.Title,
		Description: input.Description,
	})
	if err != nil {
		return jsonresp.InternalServerError(err)
	}

	return jsonresp.Success(&post)
}

// Delete deletes a single blog post by ID
func (h *postHandler) Delete(r *http.Request) httphandler.Responder {
	ctx := r.Context()

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return jsonresp.Error(err, "Invalid ID format", http.StatusBadRequest)
	}

	rows, err := h.querier.DeletePost(ctx, h.db, id)
	if err != nil {
		return jsonresp.InternalServerError(err)
	}
	if rows == 0 {
		return jsonresp.Error(nil, "Post not found", http.StatusNotFound)
	}

	return nil
}
