package mtls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

type CertGenerator struct {
	validityYears int
}

func NewCertGenerator(validityYears int) *CertGenerator {
	return &CertGenerator{
		validityYears: validityYears,
	}
}

func (cg *CertGenerator) GenerateCA(certPath, keyPath string, commonName string) error {
	if err := os.MkdirAll(filepath.Dir(certPath), 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"Morchy"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(cg.validityYears, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create CA certificate: %w", err)
	}

	if err := writeCertFile(certPath, certBytes); err != nil {
		return err
	}

	if err := writeKeyFile(keyPath, privateKey); err != nil {
		return err
	}

	return nil
}

func (cg *CertGenerator) GenerateCertificate(
	certPath, keyPath string,
	caCertPath, caKeyPath string,
	nodeID string,
	isServerCert bool,
) error {
	if err := os.MkdirAll(filepath.Dir(certPath), 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	caCertBytes, err := os.ReadFile(caCertPath)
	if err != nil {
		return fmt.Errorf("failed to read CA certificate: %w", err)
	}

	caCertBlock, _ := pem.Decode(caCertBytes)
	if caCertBlock == nil {
		return fmt.Errorf("failed to decode CA certificate PEM")
	}

	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse CA certificate: %w", err)
	}

	caKeyBytes, err := os.ReadFile(caKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read CA key: %w", err)
	}

	caKeyBlock, _ := pem.Decode(caKeyBytes)
	if caKeyBlock == nil {
		return fmt.Errorf("failed to decode CA key PEM")
	}

	caPrivateKey, err := x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse CA private key: %w", err)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %w", err)
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   nodeID,
			Organization: []string{"Morchy"},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(cg.validityYears, 0, 0),
		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{},
	}

	if isServerCert {
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
		template.DNSNames = []string{"localhost", "127.0.0.1"}
		template.IPAddresses = []net.IP{net.ParseIP("127.0.0.1")}
	} else {
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	}

	certBytes, err := x509.CreateCertificate(
		rand.Reader,
		template,
		caCert,
		&privateKey.PublicKey,
		caPrivateKey,
	)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	if err := writeCertFile(certPath, certBytes); err != nil {
		return err
	}

	if err := writeKeyFile(keyPath, privateKey); err != nil {
		return err
	}

	return nil
}

func writeCertFile(path string, certBytes []byte) error {
	certFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create certificate file: %w", err)
	}
	defer certFile.Close()

	certPEM := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}

	if err := pem.Encode(certFile, certPEM); err != nil {
		return fmt.Errorf("failed to write certificate PEM: %w", err)
	}

	return nil
}

func writeKeyFile(path string, key *rsa.PrivateKey) error {
	keyFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create key file: %w", err)
	}
	defer keyFile.Close()

	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	keyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	}

	if err := pem.Encode(keyFile, keyPEM); err != nil {
		return fmt.Errorf("failed to write key PEM: %w", err)
	}

	return nil
}
