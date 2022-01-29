package model

type User struct {
	ID       string `json:"_id" bson:"_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type NewUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
