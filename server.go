package pkiapi

import (
	"net"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opencoff/go-pki"
)

type createSvrReq struct {
	DNs          []string `json:"domain_names"`
	IPs          []string `json:"ips"`
	ValidityDays int      `json:"validity_days"`
	SignWith     string   `json:"sign_with"`
	Password     string   `json:"password"`
}

func (svr *Server) createServer(c *gin.Context) {
	cn := c.Param("cn")
	if cn == "" {
		c.AbortWithStatus(400)
		return
	}

	_, err := svr.ca.Find(cn)
	if err == nil {
		c.AbortWithStatus(409)
		return
	}

	csreq := createSvrReq{}
	if err := c.BindJSON(&csreq); err != nil {
		c.AbortWithError(400, err)
		return
	}

	if strings.Index(cn, ".") > 0 {
		csreq.DNs = append(csreq.DNs, cn)
	}

	iplist := IPList{}
	(&iplist).Set(strings.Join(csreq.IPs, ","))

	var days = time.Duration(365)
	if csreq.ValidityDays > 0 {
		days = time.Duration(csreq.ValidityDays)
	}

	var ca *pki.CA
	if len(csreq.SignWith) > 0 {
		ica, err := ca.FindCA(csreq.SignWith)
		if err != nil {
			jsonError(c, 400, "couldn't find provided signer %s", csreq.SignWith)
			return
		}
		ca = ica
	} else {
		ca = svr.ca
	}

	ci := &pki.CertInfo{
		Subject:  ca.Subject,
		Validity: time.Hour * 24 * days,

		DNSNames:    []string(csreq.DNs),
		IPAddresses: []net.IP(iplist),
	}
	ci.Subject.CommonName = cn

	// TODO: add private key password
	crt, err := ca.NewClientCert(ci, csreq.Password)
	if err != nil {
		jsonError(c, 500, "can't create user cert: %s", err)
		return
	}

	c.JSON(200, map[string]string{"cert": Cert(*crt.Certificate).String()})
}

func (svr *Server) deleteServer(c *gin.Context) {
	cn := c.Param("cn")
	if cn == "" {
		c.AbortWithStatus(400)
		return
	}

	ck, err := svr.ca.Find(cn)
	if err == pki.ErrNotFound {
		c.AbortWithStatus(404)
		return
	}

	switch {
	case ck.IsServer:
		err = svr.ca.RevokeServer(cn)
	case ck.IsCA:
		c.AbortWithStatus(404)
		return
	default:
		c.AbortWithStatus(404)
		return
	}

	if err != nil {
		jsonError(c, 500, "failed to delete server: %s", err)
		return
	}

	c.Status(204)
}

func (svr *Server) listServers(c *gin.Context) {
	certs, err := svr.ca.GetServers()
	if err != nil {
		jsonError(c, 500, "failed to list servers: %s", err)
		return
	}

	svrs := []jsonableCert{}
	for _, cert := range certs {
		svrs = append(svrs, jsonableCert{cert})
	}

	c.JSON(200, svrs)
}
