package nodeenrollment

import (
	"crypto"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"reflect"
	"strings"

	"github.com/sethvargo/go-diceware/diceware"
	"golang.org/x/crypto/hkdf"
)

func IsNil(in any) bool {
	if in == nil {
		return true
	}
	if reflect.ValueOf(in).Kind() != reflect.Ptr {
		return false
	}
	return reflect.ValueOf(in).IsNil()
}

// SubjectKeyInfoAndKeyIdFromPubKey returns the PKIX-encoded public key and the
// library-specific key ID derived from it
func SubjectKeyInfoAndKeyIdFromPubKey(pubKey crypto.PublicKey) ([]byte, string, error) {
	const op = "nodeenrollment.SubjectKeyInfoAndKeyIdFromPubKey"

	if IsNil(pubKey) {
		return nil, "", fmt.Errorf("(%s) pub key is nil", op)
	}

	pubKeyPkix, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, "", fmt.Errorf("(%s) error marshaling public key: %w", op, err)
	}
	keyId, err := KeyIdFromPkix(pubKeyPkix)
	if err != nil {
		return nil, "", fmt.Errorf("(%s) error getting key id: %w", op, err)
	}
	return pubKeyPkix, keyId, nil
}

// KeyIdFromPkix derives the library-specific key ID from the PKIX-encoed public
// key
func KeyIdFromPkix(pkixKey []byte) (string, error) {
	const op = "nodeenrollment.KeyIdFromPkix"

	if IsNil(pkixKey) {
		return "", fmt.Errorf("(%s) pkix key is nil", op)
	}

	// This never returns a non-nil error (nor is there reason for it to), so
	// ignore
	reader := hkdf.New(sha256.New, pkixKey, pkixKey, pkixKey)
	gen, _ := diceware.NewGenerator(&diceware.GeneratorInput{RandReader: reader})
	words, err := gen.Generate(KeyIdNumWords)
	if err != nil {
		return "", fmt.Errorf("(%s) error generating key id: %w", op, err)
	}
	return strings.Join(words, "-"), nil
}