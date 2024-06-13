package sdi

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"strings"

	resty "github.com/go-resty/resty/v2"
)

// SDIEnvelope defines messages received by
type SDIEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		FileSubmissionMetadata         *FileSubmissionMetadata         `xml:"MetadatiInvioFile,omitempty"`
		NonDeliveryNotificationMessage *NonDeliveryNotificationMessage `xml:"NotificaMancataConsegna,omitempty"`
		InvoiceTransmissionCertificate *InvoiceTransmissionCertificate `xml:"AttestazioneTrasmissioneFattura,omitempty"`
	} `xml:"Body"`
}

// parseMultipartResponse parses a multipart HTTP response and deserializes the content into the provided structure
func parseMultipartResponse(resp *resty.Response, response interface{}) error {
	mediaType, params, err := mime.ParseMediaType(resp.Header().Get("Content-Type"))
	if err != nil {
		return err
	}

	if !strings.HasPrefix(mediaType, "multipart/related") {
		return fmt.Errorf("unexpected content type: %s", mediaType)
	}

	reader := strings.NewReader(string(resp.Body()))
	mr := multipart.NewReader(reader, params["boundary"])

	for {
		part, err := mr.NextPart()
		if err != nil {
			break
		}
		if part.Header.Get("Content-Type") != "application/xop+xml; charset=utf-8; type=\"text/xml\"" {
			continue
		}

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(part)
		if err != nil {
			return fmt.Errorf("multipart reading error: %s", err)
		}
		xmlData := buf.String()

		err = xml.Unmarshal([]byte(xmlData), &response)
		if err != nil {
			return fmt.Errorf("parsing xml error: %s", err)
		}
	}
	return nil
}

func ParseMessage(body io.ReadCloser) error {
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	env := new(SDIEnvelope)
	err = xml.Unmarshal(data, env)
	if err != nil {
		return err
	}
	if env.Body.FileSubmissionMetadata != nil {
		log.Printf("parsing MetadatiInvioFile:\n")
	}
	if env.Body.NonDeliveryNotificationMessage != nil {
		log.Printf("parsing NotificaMancataConsegna:\n")
	}
	if env.Body.InvoiceTransmissionCertificate != nil {
		log.Printf("parsing AttestazioneTrasmissioneFattura:\n")
	}

	return nil
}
