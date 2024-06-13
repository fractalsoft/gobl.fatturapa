package sdi_test

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	sdi "github.com/invopop/gobl.fatturapa/sdi"
	"github.com/stretchr/testify/assert"
)

func handlerFunc(env *sdi.SDIEnvelope) {
	if env.Body.FileSubmissionMetadata != nil {
		log.Printf("parsing MetadatiInvioFile:\n")
	}
	if env.Body.NonDeliveryNotificationMessage != nil {
		log.Printf("parsing NotificaMancataConsegna:\n")
	}
	if env.Body.InvoiceTransmissionCertificate != nil {
		log.Printf("parsing AttestazioneTrasmissioneFattura:\n")
	}
}

func TestParseMessage(t *testing.T) {
	t.Run("parse MetadatiInvioFile", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(os.Stderr)
		}()

		message := `<?xml version='1.0' encoding='UTF-8'?>` +
			`<soapenv:Envelope xmlns:soapenv='http://schemas.xmlsoap.org/soap/envelope/' xmlns:typ='http://www.fatturapa.gov.it/sdi/ws/trasmissione/v1.0/types'>` +
			`<soapenv:Header/>` +
			`<soapenv:Body>` +
			`<ns3:MetadatiInvioFile xmlns:ns2="http://www.w3.org/2000/09/xmldsig#" xmlns:ns3="http://www.fatturapa.gov.it/sdi/messaggi/v1.0" versione="1.0">` +
			`<IdentificativoSdI>29218239</IdentificativoSdI>` +
			`<NomeFile>ESB85905495_00010.xml</NomeFile>` +
			`<CodiceDestinatario>WSBKWM</CodiceDestinatario>` +
			`<Formato>FPA12</Formato>` +
			`<TentativiInvio>1</TentativiInvio>` +
			`<MessageId>176121330</MessageId>` +
			`</ns3:MetadatiInvioFile>` +
			`</soapenv:Body>` +
			`</soapenv:Envelope>`
		reader := strings.NewReader(message)

		sdi.ParseMessage(io.NopCloser(reader), handlerFunc)

		assert.Contains(t, buf.String(), "parsing MetadatiInvioFile")
	})

	t.Run("parse NotificaMancataConsegna", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(os.Stderr)
		}()

		message := `<?xml version='1.0' encoding='UTF-8'?>` +
			`<soapenv:Envelope xmlns:soapenv='http://schemas.xmlsoap.org/soap/envelope/' xmlns:typ='http://www.fatturapa.gov.it/sdi/ws/trasmissione/v1.0/types'>` +
			`<soapenv:Header/>` +
			`<soapenv:Body>` +
			`<ns3:NotificaMancataConsegna xmlns:ns3="http://www.fatturapa.gov.it/sdi/messaggi/v1.0" xmlns:ns2="http://www.w3.org/2000/09/xmldsig#" versione="1.0">` +
			`<IdentificativoSdI>29218239</IdentificativoSdI>` +
			`<NomeFile>ESB85905495_00010.xml</NomeFile>` +
			`<DataOraRicezione>2024-05-31T14:54:02.000+02:00</DataOraRicezione>` +
			`<Descrizione>Non è stato possibile recapitare la fattura/e al destinatario.Sono in corso le necessarie verifiche,al termine delle quali si procederà ad un nuovo tentativo di trasmissione. Si rimanda pertanto ad un momento successivo l'invio della ricevuta di consegna.</Descrizione>` +
			`<MessageId>176130653</MessageId>` +
			`<Note/>` +
			`</ns3:NotificaMancataConsegna>` +
			`</soapenv:Body>` +
			`</soapenv:Envelope>`
		reader := strings.NewReader(message)

		sdi.ParseMessage(io.NopCloser(reader), handlerFunc)

		assert.Contains(t, buf.String(), "parsing NotificaMancataConsegna")
	})

	t.Run("parse AttestazioneTrasmissioneFattura", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(os.Stderr)
		}()

		message := `<?xml version='1.0' encoding='UTF-8'?>` +
			`<soapenv:Envelope xmlns:soapenv='http://schemas.xmlsoap.org/soap/envelope/' xmlns:typ='http://www.fatturapa.gov.it/sdi/ws/trasmissione/v1.0/types'>` +
			`<soapenv:Header/>` +
			`<soapenv:Body>` +
			`<ns3:AttestazioneTrasmissioneFattura xmlns:ns3="http://www.fatturapa.gov.it/sdi/messaggi/v1.0" xmlns:ns2="http://www.w3.org/2000/09/xmldsig#" versione="1.0">` +
			`<IdentificativoSdI>29218239</IdentificativoSdI>` +
			`<NomeFile>ESB85905495_00010.xml</NomeFile>` +
			`<DataOraRicezione>2024-05-31T14:54:02.000+02:00</DataOraRicezione>` +
			`<Destinatario>` +
			`<Codice>WSBKWM</Codice>` +
			`<Descrizione>Amministrazione di test - Ufficio_test</Descrizione>` +
			`</Destinatario>` +
			`<MessageId>176197456</MessageId>` +
			`<Note>Fattura</Note>` +
			`<HashFileOriginale>bc0c40728a04f06d52412d946a939540583cbe0ea0edaa5c0ba8097fc0519d16</HashFileOriginale>` +
			`</ns3:AttestazioneTrasmissioneFattura>` +
			`</soapenv:Body>` +
			`</soapenv:Envelope>`
		reader := strings.NewReader(message)

		sdi.ParseMessage(io.NopCloser(reader), handlerFunc)

		assert.Contains(t, buf.String(), "parsing AttestazioneTrasmissioneFattura")
	})
}
