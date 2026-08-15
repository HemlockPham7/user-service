package infrastructure

import "github.com/HemlockPham7/common-libs/pkg/jwtutils"

func CreateJWTProvider() (jwtutils.JWTGenerator, jwtutils.JWTValidator) {
	jwtGen, err := jwtutils.NewJWTGenerator("./private.pem")
	if err != nil {
		panic(err)
	}

	jwtVal, err := jwtutils.NewJWTValidator("./public.pem")
	if err != nil {
		panic(err)
	}

	return jwtGen, jwtVal
}
