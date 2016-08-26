package deviceserver

import (
	"encoding/json"
	"time"

	"github.com/square/go-jose"
)

// JwtSigner is the main object for simplified JWT operations
type JwtSigner struct {
	signer jose.Signer
}

// TokenFromPSK generates an JWT with signed OrgClaim
func TokenFromPSK(psk string, orgID int) (token string, err error) {

	signer := JwtSigner{}
	err = signer.Init(jose.HS256, []byte(psk))
	if err != nil {
		return "", err
	}

	// the lifetime should be shorter, but think I'm hitting some timezone issues at the moment
	orgClaim := OrgClaim{
		OrgID: orgID,
		Exp:   time.Now().Add(60 * time.Minute).Unix(),
	}

	serialized, err := signer.MarshallSignSerialize(orgClaim)
	if err != nil {
		return "", err
	}

	return serialized, nil
}

// Init creates JOSE signer
func (s *JwtSigner) Init(alg jose.SignatureAlgorithm, signingKey interface{}) error {
	var err error
	s.signer, err = jose.NewSigner(jose.HS256, signingKey)
	return err
}

// MarshallSignSerialize returns a compacted serialised JWT from a claims structure
func (s *JwtSigner) MarshallSignSerialize(in interface{}) (string, error) {
	claimJSON, err := json.Marshal(in)
	if err != nil {
		return "", err
	}

	object, err := s.signer.Sign(claimJSON)
	if err != nil {
		return "", err
	}

	serialized, err := object.CompactSerialize()
	return serialized, err
}

// ParseVerify performs signature validation and returns byte string
func ParseVerify(serialized []byte, signingKey interface{}) ([]byte, error) {
	object, err := jose.ParseSigned(string(serialized))
	if err != nil {
		return nil, err
	}

	output, err := object.Verify(signingKey)
	if err != nil {
		return nil, err
	}

	return output, nil
}
