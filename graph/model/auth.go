package model

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginReturn struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

type VerifyToken struct {
	Token string `json:"token"`
}

type ChangePassword struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}
