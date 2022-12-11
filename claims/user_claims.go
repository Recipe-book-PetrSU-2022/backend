package claims

import "github.com/golang-jwt/jwt"

type UserClaims struct {
	IntUserId     uint   `json:"id"`
	StrUserName   string `json:"username"`
	IntUserRights int    `json:"rights"`
	jwt.StandardClaims
}
