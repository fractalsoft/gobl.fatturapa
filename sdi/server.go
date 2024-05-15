package sdi

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	soap "github.com/globusdigital/soap"
)

// ServerConfig defines soap server configuration options
type ServerConfig struct {
	Host    string
	Port    string
	Verbose bool
}

// RunServer sets up a server for receiving invoices from SdI
func RunServer(config *ServerConfig) error {
	soapServer := soap.NewServer()
	if config.Verbose {
		fmt.Printf("%+v\n", soapServer)
		fmt.Printf("Soap Version: %s\n", soapServer.SoapVersion)
	}

	pathTo := "/RicezioneFatture"
	action := "RiceviFatture" // Receive Invoices
	tagName := "someRequest"

	soapServer.RegisterHandler(
		pathTo,
		action,
		tagName,
		func() interface{} {
			return nil
		},
		// func(request interface{}, w http.ResponseWriter, httpRequest *http.Request) (response interface{}, err error) {
		nil,
	)

	// Interrupt signal handling
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sig
		fmt.Println("Received interrupt signal, shutting down...")
		// err := soapServer.Close() // nolint:errcheck
		// if err != nil {
		// 	fmt.Printf("Error closing SOAP server: %v\n", err)
		// }
		os.Exit(0)
	}()

	err := http.ListenAndServe(config.Host+":"+config.Port, soapServer)
	if err != nil {
		fmt.Println("Exit with error:", err)
		return err
	}

	return nil
}
