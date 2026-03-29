package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

func LoadTLSCertificate(certFile, keyFile string) (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate pair: %w", err)
	}
	return &cert, nil
}

func LoadCACertPool(caCertFile string) (*x509.CertPool, error) {
	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to parse CA certificate")
	}

	return certPool, nil
}

func CreateServerTLSConfig(certFile, keyFile, caCertFile string) (*tls.Config, error) {
	cert, err := LoadTLSCertificate(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	certPool, err := LoadCACertPool(caCertFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{*cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
		MinVersion:   tls.VersionTLS12,
	}, nil
}

func CreateClientTLSConfig(certFile, keyFile, caCertFile string) (*tls.Config, error) {
	if certFile == "" && keyFile == "" && caCertFile == "" {
		return &tls.Config{
			InsecureSkipVerify: true,
		}, nil
	}

	cert, err := LoadTLSCertificate(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	certPool, err := LoadCACertPool(caCertFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{*cert},
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS12,
	}, nil
}
