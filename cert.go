package pkiapi

import (
	"crypto/x509"
	"fmt"

	"github.com/opencoff/go-pki"
)

type Cert x509.Certificate

func (z Cert) String() string {
	c := x509.Certificate(z)
	s, err := pki.CertificateText(&c)
	if err != nil {
		s = fmt.Sprintf("can't stringify %x (%s)", c.SerialNumber, err)
	}
	return s
}
