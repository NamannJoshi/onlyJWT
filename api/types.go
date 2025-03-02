package api

import "time"

type CreateUserReq struct {
	FullName string	`json:"fullName"`
	Email string	`json:"email"`
	IsAdmin bool	`json:"isAdmin"`
	Number int	`json:"number"`
	Password string `json:"password"`
}
type CreateUserRes struct {
	FullName string `json:"fullName"`
	Email string `json:"email"`
	IsAdmin bool `json:"isAdmin"`
	Number int `json:"number"`
}

type LoginUserReq struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRes struct {
	AccessToken string `json:"accessToken"`
	UserRes CreateUserRes `json:"userInfo"`
}

type User struct {
	ID        int       `json:"id"`
	FullName  string    `json:"fullName"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"isAdmin"`
	Number    int       `json:"number"`
	Password string			`json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewUser(u CreateUserReq) *User{
	return &User{
		FullName: u.FullName,
		Email: u.Email,
		IsAdmin: u.IsAdmin,
		Number: u.Number,
	}
}