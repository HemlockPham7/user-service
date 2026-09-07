package infrastructure

import "github.com/HemlockPham7/common-libs/pkg/jwtutils"

// CreateJWTProvider creates JWT generator and validator providers
// using the application's RSA private and public keys.
//
// Returns:
//   - A JWT generator and validator configured with the application's RSA key pair.
//
// Panics:
//   - If the private or public key cannot be read or parsed.
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
