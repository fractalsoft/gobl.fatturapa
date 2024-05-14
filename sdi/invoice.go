package sdi

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/tiaguinho/gosoap"
)

// SendInvoice sends invoice content to SdI
func SendInvoice(ctx context.Context, content []byte, config Config) error {
	soapEndpoint := config.SOAPReceiveFileEndpoint()
	if config.Verbose {
		log.Println(soapEndpoint)
	}

	opts := []HTTPClientOptFunc{
		WithCaCertPool(config.CACert),
		WithContext(ctx),
		WithTimeout(15000),
	}
	if config.Verbose {
		opts = append(opts, WithDebugClient())
	}

	httpClient := NewHTTPClient(
		opts...,
	).Build()

	soap, err := gosoap.SoapClient(soapEndpoint, httpClient)
	if err != nil {
		log.Fatalf("SoapClient error: %s", err)
	}

	// TODO File name taken from function arguments
	fileName := "invoice.xml"
	fileBody := base64.StdEncoding.EncodeToString(content)

	params := gosoap.Params{
		"NomeFile": fileName,
		"File":     fileBody,
	}

	operation := "RicevutaConsegna"
	res, err := soap.Call(operation, params)
	if err != nil {
		log.Fatalf("Call error: %s", err)
	}

	// TODO Delete when communication will be correct
	fmt.Println("Response Payload:", string(res.Payload))
	fmt.Println("Response Header:", string(res.Header))
	fmt.Println("Response Body:", string(res.Body))

	return nil
}
