package dtos

import "encoding/base64"

type CertificateDto struct {
	CaCrt  string
	TlsCrt string
	TlsKey string
}

// Convert from CertificateDto containing base64 string to CertificateDto with decoded strin
func (certDto *CertificateDto) EncodeFromBase64ToString() (*CertificateDto, error) {
	caCrt, err := base64.StdEncoding.DecodeString(certDto.CaCrt)
	if err != nil {
		return nil, err
	}
	tlsCrt, err := base64.StdEncoding.DecodeString(certDto.TlsCrt)
	if err != nil {
		return nil, err
	}
	tlsKey, err := base64.StdEncoding.DecodeString(certDto.TlsKey)
	if err != nil {
		return nil, err
	}
	return &CertificateDto{
		CaCrt:  string(caCrt),
		TlsCrt: string(tlsCrt),
		TlsKey: string(tlsKey),
	}, nil
}
