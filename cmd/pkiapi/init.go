package main

import (
	"crypto/x509/pkix"
	"fmt"
	"log"
	"time"

	"github.com/opencoff/go-pki"
	pkiapi "github.com/penguinpowernz/go-pki-api"
	"github.com/spf13/cobra"
)

var initOpts = struct {
	country string
	org     string
	ou      string
	yrs     uint
	from    string
}{}

func init() {
	// initCmd.Flags().StringVarP(&initOpts.from, "from-json", "j", "", "Initialize from an exported JSON dump")
	initCmd.Flags().StringVarP(&initOpts.country, "country", "c", "US", "Use `C` as the country name")
	initCmd.Flags().StringVarP(&initOpts.org, "organization", "O", "", "Use `O` as the organization name")
	initCmd.Flags().StringVarP(&initOpts.ou, "organization-unit", "u", "", "Use `U` as the organization unit name")
	initCmd.Flags().UintVarP(&initOpts.yrs, "validity", "V", 5, "Issue CA root cert with `N` years validity")
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init [options] <common-name>",
	Short: "init",
	Long:  "init",
	Run: func(cmd *cobra.Command, args []string) {
		cn := args[0]

		if len(cn) == 0 {
			log.Fatalf("must specify common-name")
		}

		p := pki.Config{
			Passwd:   password,
			Validity: years(initOpts.yrs),

			Subject: pkix.Name{
				Country:            []string{initOpts.country},
				Organization:       []string{initOpts.org},
				OrganizationalUnit: []string{initOpts.ou},
				CommonName:         cn,
			},
		}
		ca, err := pki.New(&p, dbfile, true)
		if err != nil {
			log.Fatalf("failed to create PKI database: %s", err)
		}

		fmt.Printf("New CA cert:\n%s\n", pkiapi.Cert(*ca.Certificate))
	},
}

// convert duration in years to time.Duration
// 365.25 days/year * 24 hours/day
// .25 days/year = 24 hours / 4 = 6 hrs
func years(n uint) time.Duration {
	day := 24 * time.Hour
	return (6 * time.Hour) + (time.Duration(n*365) * day)
}
