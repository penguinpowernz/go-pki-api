package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/opencoff/go-pki"
	pkiapi "github.com/penguinpowernz/go-pki-api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run",
	Long:  "run",
	Run: func(cmd *cobra.Command, args []string) {

		p := pki.Config{
			Passwd: password,
		}
		ca, err := pki.New(&p, dbfile, false)
		if err != nil {
			log.Fatalf("failed to open the PKI database: %s", err)
		}
		defer ca.Close()

		svr := pkiapi.NewServer(ca)
		api := gin.Default()
		svr.SetupRoutes(api)

		api.Run()
	},
}
