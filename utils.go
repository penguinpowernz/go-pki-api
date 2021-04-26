package pkiapi

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opencoff/go-pki"
)

func jsonError(c *gin.Context, status int, fmtm string, args ...interface{}) {
	c.AbortWithStatusJSON(status, map[string]string{"error": fmt.Sprintf(fmtm, args...)})
}

type IPList []net.IP

func NewIPList() *IPList {
	return &IPList{}
}

func ParseIPList(s string) *IPList {
	l := &IPList{}
	l.Set(s)
	return l
}

func (ipl *IPList) Set(s string) error {
	v := strings.Split(s, ",")
	ips := make([]net.IP, 0, 4)
	for _, x := range v {
		ip := net.ParseIP(x)
		if ip == nil {
			return fmt.Errorf("can't parse IP Address '%s'", s)
		}
		ips = append(ips, ip)
	}

	*ipl = append(*ipl, ips...)
	return nil
}

func (ipl *IPList) String() string {
	var x []string
	ips := []net.IP(*ipl)

	for i := range ips {
		x = append(x, ips[i].String())
	}
	return fmt.Sprintf("[%s]", strings.Join(x, ","))
}

type jsonableCert struct {
	*pki.Cert
}

func (cert jsonableCert) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"common_name":      cert.Subject.CommonName,
		"fingerprint":      fmt.Sprintf("%#x", cert.SerialNumber),
		"expired":          time.Now().After(cert.NotAfter),
		"expires_at":       cert.NotAfter.Unix(),
		"expires_at_human": cert.NotAfter.String(),
	})
}
