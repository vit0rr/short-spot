package urlshort

import (
	"net/http"

	"github.com/vit0rr/short-spot/pkg/deps"
	"github.com/vit0rr/short-spot/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type HTTP struct {
	service *Service
}

func NewHTTP(deps *deps.Deps, db *mongo.Database) *HTTP {
	return &HTTP{
		service: NewService(deps, db),
	}
}

// POST /short-url
func (h *HTTP) ShortUrl(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	url, err := h.service.ShortUrl(r.Context(), r.Body, *h.service.db.Client())
	if err != nil {
		log.Error(r.Context(), "Failed to create short URL", log.ErrAttr(err))
		return nil, err
	}

	return url, nil
}

// GET /{id}
func (h *HTTP) Redirect(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	url, err := h.service.Redirect(r.Context(), r.URL.Path[1:], *h.service.db.Client())
	if err != nil {
		log.Error(r.Context(), "Failed to redirect", log.ErrAttr(err))
		http.Error(w, "Not found", http.StatusNotFound)
		return nil, err
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
	return nil, nil
}
