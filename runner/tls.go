package runner

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

// Use certs via autocert
func automaticTLSConfig(certsPath string, hosts []string) (*tls.Config, error) {
	// Build a manager for the certificates
	certManager := &autocert.Manager{
		Cache:      autocert.DirCache(certsPath),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hosts...),
	}
	return certManager.TLSConfig(), nil
}

// Use generated self-signed TLS certs
func selfSignedTLSConfig(hosts []string) (*tls.Config, error) {
	cert, err := genSelfSignedCert(hosts)
	if err != nil {
		return nil, err
	}
	config := tls.Config{
		Certificates: []tls.Certificate{*cert},
	}
	return &config, nil
}

// Generates an X509 self signed certificate
func genSelfSignedCert(hosts []string) (*tls.Certificate, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %s", err)
	}

	notBefore := time.Now()
	// Self signed certs are for debugging, they should not be needed for long
	notAfter := notBefore.Add(time.Hour * 24)

	keyUsage := x509.KeyUsageCertSign |
		x509.KeyUsageKeyEncipherment |
		x509.KeyUsageDigitalSignature

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Bogus Certificates Inc"},
		},
		Version: 3,

		DNSNames:    []string{"localhost"},
		IPAddresses: []net.IP{net.IPv6loopback, net.IPv4(127, 0, 0, 1)},

		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	// DER encoded certificate in a byte array
	encodedCert, err := x509.CreateCertificate(
		rand.Reader, &template, &template, priv.Public(), priv,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %s", err)
	}

	outCert := tls.Certificate{
		Certificate: [][]byte{encodedCert},
		PrivateKey:  priv,
	}
	return &outCert, nil
}
