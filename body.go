package fatturapa

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
)

const (
	TipoDocumentoDefault = "TD01"
)

// FatturaElettronicaBody contains all invoice data apart from the parties
// involved, which are contained in FatturaElettronicaHeader.
type FatturaElettronicaBody struct {
	DatiGenerali    DatiGenerali
	DatiBeniServizi DatiBeniServizi
	DatiPagamento   DatiPagamento `xml:",omitempty"`
}

// DatiGenerali contains general data about the invoice such as retained taxes,
// invoice number, invoice date, document type, etc.
type DatiGenerali struct {
	DatiGeneraliDocumento DatiGeneraliDocumento
}

type DatiGeneraliDocumento struct {
	TipoDocumento string
	Divisa        string
	Data          string
	Numero        string
	Causale       []string
}

type DatiPagamento struct {
	CondizioniPagamento string
	DettaglioPagamento  []DettaglioPagamento
}

type DettaglioPagamento struct {
	ModalitaPagamento     string
	DataScadenzaPagamento string `xml:",omitempty"`
	ImportoPagamento      string
}

func newFatturaElettronicaBody(inv bill.Invoice) (*FatturaElettronicaBody, error) {
	return &FatturaElettronicaBody{
		DatiGenerali: DatiGenerali{
			DatiGeneraliDocumento: DatiGeneraliDocumento{
				TipoDocumento: TipoDocumentoDefault,
				Divisa:        string(inv.Currency),
				Data:          inv.IssueDate.String(),
				Numero:        inv.Code,
				Causale:       extractInvoiceReasons(inv),
			},
		},
		DatiBeniServizi: newDatiBeniServizi(inv),
		DatiPagamento: DatiPagamento{
			CondizioniPagamento: "TP02", // TODO
			DettaglioPagamento: []DettaglioPagamento{
				{
					ModalitaPagamento: "MP05", // TODO
					ImportoPagamento:  inv.Totals.Due.String(),
				},
			},
		},
	}, nil
}

func extractInvoiceReasons(inv bill.Invoice) []string {
	// find inv.Notes with NoteKey as cbc.NoteKeyReason
	var reasons []string

	for _, note := range inv.Notes {
		if note.Key == cbc.NoteKeyReason {
			reasons = append(reasons, note.Text)
		}
	}

	return reasons
}
