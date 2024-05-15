package main

import (
	"fmt"

	sdi "github.com/invopop/gobl.fatturapa/sdi"
	"github.com/spf13/cobra"
)

type serverOpts struct {
	*rootOpts
}

func server(o *rootOpts) *serverOpts {
	return &serverOpts{rootOpts: o}
}

func (c *serverOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server [host] [port]",
		Short: "Server for communication with SdI in selected environment",
		RunE:  c.runE,
	}

	cmd.Flags().BoolP("verbose", "v", false, "Logs all requests into the console")

	// f := cmd.Flags()

	return cmd
}

func (c *serverOpts) runE(cmd *cobra.Command, args []string) error {
	host := inputHost(args)
	port := inputPort(args)

	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Server start: %s:%s\n", host, port)
	}

	// err := sdi.DemoListener(host, port)
	config := &sdi.ServerConfig{
		Host:    host,
		Port:    port,
		Verbose: verbose,
	}

	err = sdi.RunServer(config)
	if err != nil {
		return fmt.Errorf("sending error: %w", err)
	}

	return nil
}

func inputHost(args []string) string {
	if len(args) > 0 && args[0] != "-" {
		return args[0]
	}
	return ""
}

func inputPort(args []string) string {
	if len(args) > 1 && args[1] != "-" {
		return args[1]
	}
	return ""
}
