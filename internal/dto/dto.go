package dto

type ReqForShortenHandler struct {
	URL string `json:"url"`
}

type ReqForAuth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
