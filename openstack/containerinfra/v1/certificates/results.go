package certificates

import (
	"encoding/pem"

	"github.com/gophercloud/gophercloud"
)

// CertificateResult temporarily contains the response from a GenerateCertificate or ImportCertificate call.
type CertificateResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a CertificateResult and extracts a certificate resource.
func (r CertificateResult) Extract() (*ClusterCertificate, error) {
	var s struct {
		ClusterID   string `json:"cluster_uuid"`
		Certificate string `json:"pem"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}

	pemBlock, _ := pem.Decode([]byte(s.Certificate))
	certificate := &ClusterCertificate{
		ClusterID:   s.ClusterID,
		Certificate: *pemBlock,
	}
	return certificate, nil
}

// Certificate represents a certificate associated with a cluster
type ClusterCertificate struct {
	ClusterID   string
	Certificate pem.Block
}

// String returns a PEM encoded string representation of the certificate
func (c ClusterCertificate) String() string {
	return string(pem.EncodeToMemory(&c.Certificate))
}

// CreateCredentialsBundleResult temporarily contains the response from a CreateCredentialsBundle call.
type CreateCredentialsBundleResult struct {
	gophercloud.Result
}

// CredentialsBundle is a collection of certificates and supporting files necessary to communicate with a cluster.
type CredentialsBundle struct {
	ClusterID     string
	COEEndpoint   string
	Certificate   pem.Block
	PrivateKey    pem.Block
	CACertificate pem.Block
	Scripts       map[string][]byte
}
