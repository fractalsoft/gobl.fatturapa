package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"

	sdi "github.com/invopop/gobl.fatturapa/sdi"
	"github.com/spf13/cobra"
)

type transmitOpts struct {
	*rootOpts
	config *sdi.Config
}

func transmit(o *rootOpts) *transmitOpts {
	return &transmitOpts{rootOpts: o}
}

func (c *transmitOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transmit [file] [environment]",
		Short: "Transmit a FatturaPA XML file to SdI in selected environment",
		RunE:  c.runE,
	}

	cmd.Flags().Bool("verbose", false, "Logs all requests into the console")
	cmd.Flags().String("ca-cert", "", "Path to a file containing the CA certificate")
	_ = cmd.MarkFlagRequired("ca-cert")
	// f := cmd.Flags()

	return cmd
}

func (c *transmitOpts) runE(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	env := inputEnvironment(args)
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return err
	}
	cert, err := cmd.Flags().GetString("ca-cert")
	if err != nil {
		return err
	}
	certPool, err := loadCert(cert)
	if err != nil {
		return err
	}

	input, err := openInput(cmd, args)
	if err != nil {
		return err
	}
	defer input.Close() // nolint:errcheck

	config, err := parseEnv(env)
	if err != nil {
		return err
	}

	config.Verbose = verbose
	config.CACert = certPool

	c.config = config

	if verbose {
		fmt.Printf("Environment: %s\n", c.config.Environment)
	}

	data, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	err = sdi.SendInvoice(ctx, data, *c.config)
	if err != nil {
		return fmt.Errorf("sending error: %w", err)
	}

	return nil
}

func parseEnv(env string) (*sdi.Config, error) {
	switch env {
	case "dev":
		return &sdi.DevelopmentSdIConfig, nil
	case "test":
		return &sdi.TestSdIConfig, nil
	case "prod":
		return &sdi.ProductionSdIConfig, nil
	default:
		return nil, fmt.Errorf("wrong environment: %s", env)
	}
}

func loadCert(path string) (*x509.CertPool, error) {
	pubPEM, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(pubPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: " + err.Error())
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(cert)

	return caCertPool, nil
}
