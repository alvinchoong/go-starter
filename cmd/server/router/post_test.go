package router_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-starter/cmd/server/router"
	"go-starter/internal/mocks"
	"go-starter/internal/models"
	"go-starter/internal/pkg/ptr"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_PostHandler_Create(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2025, 1, 17, 23, 51, 43, 0, time.UTC)
	fixedUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	testCases := []struct {
		desc       string
		mockFunc   func(*mocks.Querier)
		wantStatus int
		wantBody   string
	}{
		{
			desc: "success",
			mockFunc: func(m *mocks.Querier) {
				m.On("CreatePost", mock.Anything, mock.Anything, mock.Anything).
					Return(models.Post{
						ID:          fixedUUID,
						Title:       "Post title",
						Description: ptr.Ref("Post description"),
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"id":"550e8400-e29b-41d4-a716-446655440000","title":"Post title","description":"Post description","created_at":"2025-01-17T23:51:43Z","updated_at":"2025-01-17T23:51:43Z"}`,
		},
		{
			desc: "fail",
			mockFunc: func(m *mocks.Querier) {
				m.On("CreatePost", mock.Anything, mock.Anything, mock.Anything).
					Return(models.Post{}, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// Given:
			mockQ := &mocks.Querier{}
			tc.mockFunc(mockQ)
			h := router.NewPostHandler(nil, mockQ)
			r := httptest.NewRequest(http.MethodPost, "/api/v1/posts", nil)
			w := httptest.NewRecorder()

			// When:
			h.Create(r, router.CreatePostParams{
				Title:       "Post title",
				Description: "Post description",
			}).Respond(w, r)

			got := w.Result()
			defer got.Body.Close()
			gotBodyBytes, err := io.ReadAll(got.Body)
			require.NoError(t, err)

			// Then:
			assert.Equal(t, tc.wantStatus, got.StatusCode)
			assert.JSONEq(t, tc.wantBody, string(gotBodyBytes))
		})
	}
}

func Test_PostHandler_List(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2025, 1, 18, 0, 13, 2, 0, time.UTC)
	fixedUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	testCases := []struct {
		desc       string
		mockFunc   func(*mocks.Querier)
		wantStatus int
		wantBody   string
	}{
		{
			desc: "success | one result",
			mockFunc: func(m *mocks.Querier) {
				m.On("ListPosts", mock.Anything, mock.Anything).
					Return([]models.Post{{
						ID:          fixedUUID,
						Title:       "Post title",
						Description: ptr.Ref("Post description"),
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					}}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   `[{"id":"550e8400-e29b-41d4-a716-446655440000","title":"Post title","description":"Post description","created_at":"2025-01-18T00:13:02Z","updated_at":"2025-01-18T00:13:02Z"}]`,
		},
		{
			desc: "success | no results",
			mockFunc: func(m *mocks.Querier) {
				m.On("ListPosts", mock.Anything, mock.Anything).
					Return([]models.Post{}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   `[]`,
		},
		{
			desc: "fail",
			mockFunc: func(m *mocks.Querier) {
				m.On("ListPosts", mock.Anything, mock.Anything).
					Return([]models.Post{}, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// Given:
			mockQ := &mocks.Querier{}
			tc.mockFunc(mockQ)
			h := router.NewPostHandler(nil, mockQ)
			r := httptest.NewRequest(http.MethodGet, "/api/v1/posts", nil)
			w := httptest.NewRecorder()

			// When:
			h.List(r).Respond(w, r)

			got := w.Result()
			defer got.Body.Close()
			gotBodyBytes, err := io.ReadAll(got.Body)
			require.NoError(t, err)

			// Then:
			assert.Equal(t, tc.wantStatus, got.StatusCode)
			assert.JSONEq(t, tc.wantBody, string(gotBodyBytes))
		})
	}
}

func Test_PostHandler_Get(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2025, 1, 18, 0, 13, 2, 0, time.UTC)
	fixedUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	testCases := []struct {
		desc       string
		given      string
		mockFunc   func(*mocks.Querier)
		wantStatus int
		wantBody   string
	}{
		{
			desc:  "success",
			given: fixedUUID.String(),
			mockFunc: func(m *mocks.Querier) {
				m.On("GetPost", mock.Anything, mock.Anything, fixedUUID).
					Return(models.Post{
						ID:          fixedUUID,
						Title:       "Post title",
						Description: ptr.Ref("Post description"),
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"id":"550e8400-e29b-41d4-a716-446655440000","title":"Post title","description":"Post description","created_at":"2025-01-18T00:13:02Z","updated_at":"2025-01-18T00:13:02Z"}`,
		},
		{
			desc: "not found",
			mockFunc: func(m *mocks.Querier) {
				m.On("GetPost", mock.Anything, mock.Anything, fixedUUID).
					Return(models.Post{}, pgx.ErrNoRows)
			},
			given:      fixedUUID.String(),
			wantStatus: http.StatusNotFound,
			wantBody:   `{"error":"Post not found"}`,
		},
		{
			desc:       "invalid uuid",
			mockFunc:   func(m *mocks.Querier) {},
			given:      "invalid-uuid",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"Invalid ID format"}`,
		},
		{
			desc: "db error",
			mockFunc: func(m *mocks.Querier) {
				m.On("GetPost", mock.Anything, mock.Anything, fixedUUID).
					Return(models.Post{}, errors.New("db error"))
			},
			given:      fixedUUID.String(),
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// Given:
			mockQ := &mocks.Querier{}
			tc.mockFunc(mockQ)
			h := router.NewPostHandler(nil, mockQ)
			r := httptest.NewRequest(http.MethodGet, "/api/v1/posts/"+tc.given, nil)
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.given)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			// When:
			h.Get(r).Respond(w, r)

			got := w.Result()
			defer got.Body.Close()
			gotBodyBytes, err := io.ReadAll(got.Body)
			require.NoError(t, err)

			// Then:
			assert.Equal(t, tc.wantStatus, got.StatusCode)
			assert.JSONEq(t, tc.wantBody, string(gotBodyBytes))
		})
	}
}

func Test_PostHandler_Update(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2025, 1, 18, 0, 13, 2, 0, time.UTC)
	fixedUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	testCases := []struct {
		desc       string
		given      string
		mockFunc   func(*mocks.Querier)
		input      router.UpdatePostParams
		wantStatus int
		wantBody   string
	}{
		{
			desc:  "success",
			given: fixedUUID.String(),
			mockFunc: func(m *mocks.Querier) {
				m.On("UpdatePost", mock.Anything, mock.Anything, models.UpdatePostParams{
					ID:          fixedUUID,
					Title:       "Updated title",
					Description: ptr.Ref("Updated description"),
				}).Return(models.Post{
					ID:          fixedUUID,
					Title:       "Updated title",
					Description: ptr.Ref("Updated description"),
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
				}, nil)
			},
			input: router.UpdatePostParams{
				Title:       "Updated title",
				Description: ptr.Ref("Updated description"),
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"id":"550e8400-e29b-41d4-a716-446655440000","title":"Updated title","description":"Updated description","created_at":"2025-01-18T00:13:02Z","updated_at":"2025-01-18T00:13:02Z"}`,
		},
		{
			desc:       "invalid uuid",
			given:      "invalid-uuid",
			mockFunc:   func(m *mocks.Querier) {},
			input:      router.UpdatePostParams{},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"Invalid ID format"}`,
		},
		{
			desc:  "db error",
			given: fixedUUID.String(),
			mockFunc: func(m *mocks.Querier) {
				m.On("UpdatePost", mock.Anything, mock.Anything, mock.Anything).
					Return(models.Post{}, errors.New("db error"))
			},
			input:      router.UpdatePostParams{},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// Given:
			mockQ := &mocks.Querier{}
			tc.mockFunc(mockQ)
			h := router.NewPostHandler(nil, mockQ)
			r := httptest.NewRequest(http.MethodPut, "/api/v1/posts/"+tc.given, nil)
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.given)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			// When:
			h.Update(r, tc.input).Respond(w, r)

			got := w.Result()
			defer got.Body.Close()
			gotBodyBytes, err := io.ReadAll(got.Body)
			require.NoError(t, err)

			// Then:
			assert.Equal(t, tc.wantStatus, got.StatusCode)
			assert.JSONEq(t, tc.wantBody, string(gotBodyBytes))
		})
	}
}

func Test_PostHandler_Delete(t *testing.T) {
	t.Parallel()

	fixedUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	testCases := []struct {
		desc       string
		given      string
		mockFunc   func(*mocks.Querier)
		wantStatus int
		wantBody   string
	}{
		{
			desc:  "success",
			given: fixedUUID.String(),
			mockFunc: func(m *mocks.Querier) {
				m.On("DeletePost", mock.Anything, mock.Anything, fixedUUID).
					Return(int64(1), nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   "",
		},
		{
			desc:  "not found",
			given: fixedUUID.String(),
			mockFunc: func(m *mocks.Querier) {
				m.On("DeletePost", mock.Anything, mock.Anything, fixedUUID).
					Return(int64(0), nil)
			},
			wantStatus: http.StatusNotFound,
			wantBody:   `{"error":"Post not found"}`,
		},
		{
			desc:       "invalid uuid",
			given:      "invalid-uuid",
			mockFunc:   func(m *mocks.Querier) {},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"Invalid ID format"}`,
		},
		{
			desc: "fail",
			mockFunc: func(m *mocks.Querier) {
				m.On("DeletePost", mock.Anything, mock.Anything, fixedUUID).
					Return(int64(0), errors.New("db error"))
			},
			given:      fixedUUID.String(),
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// Given:
			mockQ := &mocks.Querier{}
			tc.mockFunc(mockQ)
			h := router.NewPostHandler(nil, mockQ)
			r := httptest.NewRequest(http.MethodDelete, "/api/v1/posts/"+tc.given, nil)
			w := httptest.NewRecorder()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.given)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			// When:
			resp := h.Delete(r)
			if resp != nil {
				resp.Respond(w, r)
			}

			got := w.Result()
			defer got.Body.Close()
			gotBodyBytes, err := io.ReadAll(got.Body)

			require.NoError(t, err)

			// Then:
			assert.Equal(t, tc.wantStatus, got.StatusCode)
			assert.Equal(t, tc.wantBody, strings.TrimSpace(string(gotBodyBytes)))
		})
	}
}
