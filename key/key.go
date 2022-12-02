package key

import (
	"encoding/json"
	"errors"
	"github.com/sujit-baniya/pkg/paseto"
)

func Generate(secret string, payload []byte) (string, error) {
	claims := paseto.CustomClaim(payload)
	pv4 := paseto.NewPV4Local()
	symK, err := paseto.NewSymmetricKey([]byte(secret), paseto.Version4)
	if err != nil {
		return "", err
	}

	token, err := pv4.Encrypt(symK, claims)
	if err != nil {
		return "", err
	}
	return token, nil
}

func Validate(secret, token string) ValidatedKey {
	pv4 := paseto.NewPV4Local()
	symK, err := paseto.NewSymmetricKey([]byte(secret), paseto.Version4)
	if err != nil {
		return ValidatedKey{Error: err, Valid: false}
	}
	tk := pv4.Decrypt(token, symK)
	if tk.Err() != nil {
		return ValidatedKey{Error: tk.Err(), Valid: false}
	}

	if tk.HasFooter() {
		return ValidatedKey{Error: errors.New("footer was not passed to the library"), Valid: false}
	}

	var cc paseto.CustomClaim
	if err := tk.ScanClaims(&cc); err != nil {
		return ValidatedKey{Error: err, Valid: false}
	}
	return ValidatedKey{Payload: cc, Valid: true}
}

type ValidatedKey struct {
	Payload paseto.CustomClaim
	Valid   bool
	Error   error
}

func (v *ValidatedKey) Unmarshal(result interface{}) error {
	return json.Unmarshal(v.Payload, &result)
}
