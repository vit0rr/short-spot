package users

import (
	"context"
	"encoding/json"
	"io"

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
	Message    string      `json:"message"`
	InsertedID interface{} `json:"insertedID,omitempty"`
}

func NewService(deps *deps.Deps, db *mongo.Database) *Service {
	return &Service{
		deps: deps,
		db:   db,
	}
}

func (s *Service) List(c context.Context, dbclient mongo.Client) ([]User, error) {
	coll := dbclient.Database("shortspot").Collection("users")
	cursor, err := coll.Find(c, bson.M{})
	if err != nil {
		log.Error(c, "Failed to fetch users from database", log.ErrAttr(err))
		return nil, err
	}

	var users []User
	if err := cursor.All(c, &users); err != nil {
		log.Error(c, "Failed to decode users", log.ErrAttr(err))
		return nil, err
	}

	return users, nil
}

// POST /users/create
func (s *Service) Create(c context.Context, b io.ReadCloser, dbclient mongo.Client) (*Response, error) {
	var user User
	err := json.NewDecoder(b).Decode(&user)
	if err != nil {
		log.Error(c, "Failed to decode request body", log.ErrAttr(err))
		return nil, err
	}
	defer b.Close()

	if user.Email == "" || user.Name == "" {
		log.Error(c, "Invalid user data", log.ErrAttr(err))
		return &Response{
			Message: "Invalid user data",
		}, nil
	}

	coll := dbclient.Database("shortspot").Collection("users")
	result, err := coll.InsertOne(c, user)
	if err != nil {
		log.Error(c, "Failed to insert user into database", log.ErrAttr(err))
		return nil, err
	}

	return &Response{
		Message:    "User created successfully",
		InsertedID: result.InsertedID,
	}, nil
}
