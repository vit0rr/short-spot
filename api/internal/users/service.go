package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/vit0rr/short-spot/pkg/deps"
	"github.com/vit0rr/short-spot/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	deps *deps.Deps
}

func NewService(deps *deps.Deps) *Service {
	return &Service{
		deps: deps,
	}
}

func (s *Service) List(c context.Context, dbclient mongo.Client) ([]map[string]interface{}, error) {
	coll := dbclient.Database("shortspot").Collection("users")
	cursor, err := coll.Find(c, map[string]interface{}{})
	if err != nil {
		log.Error(c, "Failed to fetch users from database", log.ErrAttr(err))
		return nil, err
	}

	var users []map[string]interface{}
	for cursor.Next(context.Background()) {
		var user map[string]interface{}
		err := cursor.Decode(&user)
		if err != nil {
			log.Error(c, "Failed to decode user", log.ErrAttr(err))
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// POST /users/create
func (s *Service) Create(r *http.Request, h *HTTP) ([]map[string]interface{}, error) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Error(r.Context(), "Failed to decode request body", log.ErrAttr(err))
		return nil, err
	}

	defer r.Body.Close()

	coll := h.service.deps.DBClient.Database("shortspot").Collection("users")
	result, err := coll.InsertOne(r.Context(), user)
	if err != nil {
		log.Error(r.Context(), "Failed to insert user into database", log.ErrAttr(err))
		return nil, err
	}

	return []map[string]interface{}{
		{
			"message":    "User created successfully",
			"insertedID": result.InsertedID,
		},
	}, nil
}
