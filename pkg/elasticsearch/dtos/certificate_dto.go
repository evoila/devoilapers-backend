package dtos

import "encoding/base64"

type CertificateDto struct {
	CaCrt  string `formWidget:"file" formTitle:"CA certificate:"`
	TlsCrt string `formWidget:"file" formTitle:"TLS certificate:"`
	TlsKey string `formWidget:"file" formTitle:"TLS Key:"`
}

// 
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
