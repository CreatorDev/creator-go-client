package deviceserver

import (
	"encoding/json"
	"time"

	"github.com/square/go-jose"
)

type JwtSigner struct {
	signer jose.Signer
}

func TokenFromPSK(psk string) (token string, err error) {

	signer := JwtSigner{}
	err = signer.Init(jose.HS256, []byte(psk))
	if err != nil {
		return "", err
	}

	// the lifetime should be shorter, but think I'm hitting some timezone issues at the moment
	orgClaim := OrgClaim{
		OrgID: 0,
		Exp:   time.Now().Add(60 * time.Minute).Unix(),
	}

	serialized, err := signer.MarshallSignSerialize(orgClaim)
	if err != nil {
		return "", err
	}

	return serialized, nil
}

func (s *JwtSigner) Init(alg jose.SignatureAlgorithm, signingKey interface{}) error {
	var err error
	s.signer, err = jose.NewSigner(jose.HS256, signingKey)
	return err
}

func (s *JwtSigner) MarshallSignSerialize(in interface{}) (string, error) {
	claimJson, err := json.Marshal(in)
	if err != nil {
		return "", err
	}

	object, err := s.signer.Sign(claimJson)
	if err != nil {
		return "", err
	}

	serialized, err := object.CompactSerialize()
	return serialized, err
}

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
