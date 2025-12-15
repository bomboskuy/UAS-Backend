package models

import ("github.com/golang-jwt/jwt/v5")

type JWTClaims struct {
	UserID      string   `json:"user_id"`      
	RoleID      string   `json:"role_id"`      
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

type MetaInfo struct {
	Status 		string 		`json:"status"`
	Message 	string 		`json:"message"`
	Data 		interface{} `json:"data,omitempty"`
	Errors 		interface{} `json:"errors,omitempty"`
	Meta 		interface{} `json:"meta,omitempty"`
}

type PaginationMeta struct {
	Page 		int `json:"page"`
	Limit 		int `json:"limit"`
	TotalData 	int `json:"total_data"`
	TotalPages 	int `json:"total_pages"`
}