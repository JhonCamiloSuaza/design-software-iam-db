package models

import "github.com/golang-jwt/jwt/v5"

type RegisterInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotPasswordInput struct {
	Email string `json:"email"`
}

type ResetPasswordInput struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type ManagedUserInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	ActorType string `json:"actorType"`
	RoleName  string `json:"roleName"`
	IsActive  *bool  `json:"isActive"`
}

type Claims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type LoginUser struct {
	ID           string
	Email        string
	PasswordHash string
	IsActive     bool
}

type ProfileUser struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	ActorType string `json:"actorType"`
	IsActive  bool   `json:"isActive"`
}

type Permission struct {
	Role    string `json:"role"`
	Feature string `json:"feature"`
	Scope   string `json:"scope"`
}

type Profile struct {
	User        ProfileUser  `json:"user"`
	Permissions []Permission `json:"permissions"`
}

type ManagedUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	ActorType string `json:"actorType"`
	IsActive  bool   `json:"isActive"`
	Roles     string `json:"roles"`
}

type Role struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type Summary struct {
	Users    int `json:"users"`
	Roles    int `json:"roles"`
	Features int `json:"features"`
	Modules  int `json:"modules"`
}
