/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

// Package vc-rest VC Issuer, Verifier and Holder REST API.
//
// Terms Of Service:
//
//	Schemes: http, https
//	Version: 0.1.0
//	License: SPDX-License-Identifier: Apache-2.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package main

import (
	"github.com/spf13/cobra"
	"github.com/trustbloc/logutil-go/pkg/log"

	"github.com/trustbloc/vcs/cmd/vc-rest/startcmd"
)

var logger = log.New("vc-rest")

func main() {
	rootCmd := &cobra.Command{
		Use: "vc-rest",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	rootCmd.AddCommand(startcmd.GetStartCmd())

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal("Failed to run vc-rest", log.WithError(err))
	}
}
