package users

type User struct {
	Name  string `bson:"name"`
	Email string `bson:"email"`
}
