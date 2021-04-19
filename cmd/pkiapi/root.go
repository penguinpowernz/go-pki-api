package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var password, dbfile string

func init() {
	if v := os.Getenv("PKI_PASSWORD"); v != "" {
		password = v
	}

	if fn := os.Getenv("PKI_PASSWORD_FILE"); fn != "" {
		data, err := ioutil.ReadFile(fn)
		if err != nil {
			log.Fatalf("couldn't read password file %s: %s", fn, err)
		}
		password = string(data)
	}

	rootCmd.PersistentFlags().StringVarP(&dbfile, "db", "d", os.Getenv("PKI_DB"), "the path to the PKI database")
}

var rootCmd = &cobra.Command{
	Use:   "pkiapi",
	Short: "pkiapi",
	Long:  "pkiapi",
	Run: func(cmd *cobra.Command, args []string) {
		runCmd.Run(cmd, args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
