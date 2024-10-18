package urlshort

type Urls struct {
	OriginalUrl  string `bson:"original_url"`
	ShortenedUrl string `bson:"shortened_url"`
	Id           string `bson:"id"`
}

type Body struct {
	OriginalUrl string `json:"originalUrl"`
}
