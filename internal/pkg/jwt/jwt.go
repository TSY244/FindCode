package jwt

import "github.com/golang-jwt/jwt/v4"

const (
	SIGNING_KEY = "234564567890dxfcgvhbjnkml56rtyguhjnk,././78o123//..VGYBHNJKM@#$%^&*"
)

const (
	RoleAdmin = 1
	RoleUser  = 2
)

type Info struct {
	Role     int
	UserName string
	Token    string
}

type InfoClaims struct {
	Info Info
	jwt.RegisteredClaims
}

func NewJwt(info Info) (string, error) {
	infoClaims := InfoClaims{
		Info:             info,
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, infoClaims)
	ss, err := token.SignedString([]byte(SIGNING_KEY))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func ParseJwt(ss string) (Info, error) {
	t, err := jwt.ParseWithClaims(ss, &InfoClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SIGNING_KEY), nil
	})
	if err != nil {
		return Info{}, err
	}
	if claims, ok := t.Claims.(*InfoClaims); ok && t.Valid {
		return claims.Info, nil
	}
	return Info{}, nil
}
