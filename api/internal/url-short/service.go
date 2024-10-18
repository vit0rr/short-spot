package urlshort

import (
	"context"
	"encoding/json"
	"io"

	"github.com/sqids/sqids-go"
	"github.com/vit0rr/short-spot/pkg/deps"
	"github.com/vit0rr/short-spot/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	deps *deps.Deps
	db   *mongo.Database
}

type Response struct {
	Message      string      `json:"message"`
	ShortenedUrl interface{} `json:"shortenedUrl,omitempty"`
}

func NewService(deps *deps.Deps, db *mongo.Database) *Service {
	return &Service{
		deps: deps,
		db:   db,
	}
}

// Get short URL
func (s *Service) ShortUrl(c context.Context, b io.ReadCloser, dbclient mongo.Client) (*Response, error) {
	sq, _ := sqids.New(sqids.Options{
		MinLength: 10,
	})

	var body Body
	err := json.NewDecoder(b).Decode(&body)
	if err != nil {
		log.Error(c, "Failed to decode request body", log.ErrAttr(err))
		return &Response{
			Message: "Failed to decode request body",
		}, err
	}
	defer b.Close()

	if body.OriginalUrl == "" {
		log.Error(c, "Invalid URL data", log.ErrAttr(err))
		return &Response{
			Message: "Invalid URL data",
		}, err
	}

	id, err := sq.Encode([]uint64{1, 2, 3})
	if err != nil {
		log.Error(c, "Failed to encode URL", log.ErrAttr(err))
		return &Response{
			Message: "Failed to encode URL",
		}, err
	}

	shortenedUrl := "http://localhost:8080/" + id

	coll := dbclient.Database("shortspot").Collection("shorturls")
	_, err = coll.InsertOne(c, Urls{
		OriginalUrl:  body.OriginalUrl,
		ShortenedUrl: shortenedUrl,
		Id:           id,
	}, nil)
	if err != nil {
		log.Error(c, "Failed to insert short URL into database", log.ErrAttr(err))
		return &Response{
			Message: "Failed to insert short URL into database",
		}, err
	}

	return &Response{
		Message:      "Short URL created",
		ShortenedUrl: shortenedUrl,
	}, nil
}

// Redirect to original URL
func (s *Service) Redirect(c context.Context, id string, dbclient mongo.Client) (string, error) {
	coll := dbclient.Database("shortspot").Collection("shorturls")
	var url Urls
	err := coll.FindOne(c, bson.M{"id": id}).Decode(&url)
	if err != nil {
		log.Error(c, "Failed to fetch short URL from database", log.ErrAttr(err))
		return "", err
	}

	return url.OriginalUrl, nil
}
