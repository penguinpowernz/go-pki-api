package pkiapi

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opencoff/go-pki"
)

type createUserReq struct {
	Email        string `json:"email"`
	ValidityDays int    `json:"validity_days"`
	SignWith     string `json:"sign_with"`
	Password     string `json:"password"`
}

func (svr *Server) createUser(c *gin.Context) {
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

	cur := createUserReq{}
	if err := c.BindJSON(&cur); err != nil {
		c.AbortWithError(400, err)
		return
	}

	// use CN as EmailAddress if one is not provided
	var emails []string
	if len(cur.Email) > 0 {
		emails = []string{cur.Email}
	} else if strings.Index(cn, "@") > 0 {
		emails = []string{cn}
	}

	var days = time.Duration(365)
	if cur.ValidityDays > 0 {
		days = time.Duration(cur.ValidityDays)
	}

	var ca *pki.CA
	if len(cur.SignWith) > 0 {
		ica, err := ca.FindCA(cur.SignWith)
		if err != nil {
			jsonError(c, 400, "couldn't find provided signer %s", cur.SignWith)
			return
		}
		ca = ica
	} else {
		ca = svr.ca
	}

	ci := &pki.CertInfo{
		Subject:        ca.Subject,
		Validity:       time.Hour * 24 * days,
		EmailAddresses: emails,
	}
	ci.Subject.CommonName = cn

	// TODO: add private key password
	crt, err := ca.NewClientCert(ci, cur.Password)
	if err != nil {
		jsonError(c, 500, "can't create user cert: %s", err)
		return
	}

	c.JSON(200, map[string]string{"cert": Cert(*crt.Certificate).String()})
}

func (svr *Server) deleteUser(c *gin.Context) {
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
		c.AbortWithStatus(404)
		return
	case ck.IsCA:
		c.AbortWithStatus(404)
		return
	default:
		err = svr.ca.RevokeClient(cn)
	}

	if err != nil {
		jsonError(c, 500, "failed to delete user: %s", err)
		return
	}

	c.Status(204)
}

func (svr *Server) listUsers(c *gin.Context) {
	certs, err := svr.ca.GetClients()
	if err != nil {
		jsonError(c, 500, "failed to list users: %s", err)
		return
	}

	users := []jsonableCert{}
	for _, cert := range certs {
		users = append(users, jsonableCert{cert})
	}

	c.JSON(200, users)
}
