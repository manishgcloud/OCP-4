package tls

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
)

// CertInterface contains cert.
type CertInterface interface {
	// Cert returns the certificate.
	Cert() []byte
}

// CertKeyInterface contains a private key and the associated cert.
type CertKeyInterface interface {
	CertInterface
	// Key returns the private key.
	Key() []byte
}

// CertKey contains the private key and the cert.
type CertKey struct {
	CertRaw  []byte
	KeyRaw   []byte
	FileList []*asset.File
}

// Cert returns the certificate.
func (c *CertKey) Cert() []byte {
	return c.CertRaw
}

// Key returns the private key.
func (c *CertKey) Key() []byte {
	return c.KeyRaw
}

// Files returns the files generated by the asset.
func (c *CertKey) Files() []*asset.File {
	return c.FileList
}

// CertFile returns the certificate file.
func (c *CertKey) CertFile() *asset.File {
	return c.FileList[1]
}

func (c *CertKey) generateFiles(filenameBase string) {
	c.FileList = []*asset.File{
		{
			Filename: assetFilePath(filenameBase + ".key"),
			Data:     c.KeyRaw,
		},
		{
			Filename: assetFilePath(filenameBase + ".crt"),
			Data:     c.CertRaw,
		},
	}
}

// Load is a no-op because TLS assets are not written to disk.
func (c *CertKey) Load(asset.FileFetcher) (bool, error) {
	return false, nil
}

// AppendParentChoice dictates whether the parent's cert is to be added to the
// cert.
type AppendParentChoice bool

const (
	// AppendParent indicates that the parent's cert should be added.
	AppendParent AppendParentChoice = true
	// DoNotAppendParent indicates that the parent's cert should not be added.
	DoNotAppendParent AppendParentChoice = false
)

// SignedCertKey contains the private key and the cert that's
// signed by the parent CA.
type SignedCertKey struct {
	CertKey
}

// Generate generates a cert/key pair signed by the specified parent CA.
func (c *SignedCertKey) Generate(
	cfg *CertCfg,
	parentCA CertKeyInterface,
	filenameBase string,
	appendParent AppendParentChoice,
) error {
	var key *rsa.PrivateKey
	var crt *x509.Certificate
	var err error

	caKey, err := PemToPrivateKey(parentCA.Key())
	if err != nil {
		return errors.Wrap(err, "failed to parse rsa private key")
	}

	caCert, err := PemToCertificate(parentCA.Cert())
	if err != nil {
		return errors.Wrap(err, "failed to parse x509 certificate")
	}

	key, crt, err = GenerateSignedCertificate(caKey, caCert, cfg)
	if err != nil {
		return errors.Wrap(err, "failed to generate signed cert/key pair")
	}

	c.KeyRaw = PrivateKeyToPem(key)
	c.CertRaw = CertToPem(crt)

	if appendParent {
		c.CertRaw = bytes.Join([][]byte{c.CertRaw, CertToPem(caCert)}, []byte("\n"))
	}

	c.generateFiles(filenameBase)

	return nil
}

// SelfSignedCertKey contains the private key and the cert that's self-signed.
type SelfSignedCertKey struct {
	CertKey
}

// Generate generates a cert/key pair signed by the specified parent CA.
func (c *SelfSignedCertKey) Generate(
	cfg *CertCfg,
	filenameBase string,
) error {
	key, crt, err := GenerateSelfSignedCertificate(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to generate self-signed cert/key pair")
	}

	c.KeyRaw = PrivateKeyToPem(key)
	c.CertRaw = CertToPem(crt)

	c.generateFiles(filenameBase)

	return nil
}
