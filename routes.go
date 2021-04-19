package pkiapi

import (
	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/opencoff/go-pki"
)

func NewServer(ca *pki.CA) *Server {
	return &Server{ca: ca}
}

type Server struct {
	ca *pki.CA
}

func (svr *Server) SetupRoutes(r gin.IRouter) {
	r.GET("/servers", svr.listServers)
	r.GET("/servers/:cn/export", svr.ExportCert)
	r.POST("/servers/:cn", svr.createServer)
	r.DELETE("/servers/:cn", svr.deleteServer)

	r.GET("/users", svr.listUsers)
	r.GET("/users/:cn/export", svr.ExportCert)
	r.POST("/users/:cn", svr.createUser)
	r.DELETE("/users/:cn", svr.deleteUser)

	r.GET("/crl/:days", svr.generateCRL)
}

func (svr *Server) ExportCert(c *gin.Context) {
	cn := c.Param("cn")
	if cn == "" {
		c.AbortWithStatus(400)
		return
	}

	chain := c.Query("chain") != ""
	withCA := c.Query("ca") != ""

	cert, err := svr.ca.Find(cn)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	var pem []byte
	var key []byte
	if cert.IsCA && chain {
		cas, err := svr.ca.ChainFor(cert)
		if err != nil {
			jsonError(c, 500, "can't find cert chain: %s", err)
			return
		}

		var cw bytes.Buffer
		for i := range cas {
			ck := cas[i]
			cw.Write(ck.PEM())
		}

		pem = cw.Bytes()
		_, key = cert.PEM()
	} else {
		pem, key = cert.PEM()
	}

	data := map[string]string{
		"key": string(key),
		"pem": string(pem),
	}

	if withCA {
		data["ca"] = string(svr.ca.PEM())
	}

	c.JSON(200, data)
}
