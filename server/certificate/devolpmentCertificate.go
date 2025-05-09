package certificate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"math/big"
)

var ErrUnableToGenerateSelfSignedCertiticate = errors.New("unable to generate self sign certificate")

func CreateSelfSignedCertificate(a *tls.ClientHelloInfo) (*tls.Certificate, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return nil, errors.Join(ErrUnableToGenerateSelfSignedCertiticate, errors.New("error generating private key"))
	}

	//example from crypto/tls/generate_cert.go
	upperLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, upperLimit)
	if err != nil {
		return nil, errors.Join(ErrUnableToGenerateSelfSignedCertiticate, errors.New("error generating serial number"))
	}

	certificateTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Local Development"},
			CommonName:   "Local Development",
		},
	}

	certificateBytes, err := x509.CreateCertificate(rand.Reader, &certificateTemplate, &certificateTemplate, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}

	return &tls.Certificate{
		Certificate: [][]byte{certificateBytes},
		PrivateKey:  privateKey,
	}, nil
}
