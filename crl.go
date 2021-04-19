package pkiapi

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (svr *Server) generateCRL(c *gin.Context) {
	sdays := c.Param("days")
	days, err := strconv.Atoi(sdays)
	if err != nil {
		jsonError(c, 400, "days must be specified as an integer")
		return
	}

	pem, err := svr.ca.CRL(days)
	if err != nil {
		jsonError(c, 500, "failed to create CRL: %s", err)
		return
	}

	c.String(200, "utf8", pem)
}
