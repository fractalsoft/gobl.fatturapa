package fatturapa

import (
	"errors"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/it"
)

const (
	RegimeFiscaleDefault = "RF01"
)

// Party contains data related to the party
type Party struct {
	DatiAnagrafici *DatiAnagrafici
	Sede           *Address
}

// DatiAnagrafici contains information related to an individual or company
type DatiAnagrafici struct {
	IdFiscaleIVA *TaxID `xml:",omitempty"`
	// CodiceFiscale is the Italian fiscal code, distinct from TaxID
	CodiceFiscale string `xml:",omitempty"`
	Anagrafica    *Anagrafica
	// RegimeFiscale identifies the tax system to be applied
	// Has the form RFXX where XX is numeric; required only for the supplier
	RegimeFiscale string `xml:",omitempty"`
}

// Anagrafica contains further party information
type Anagrafica struct {
	// Name of the organization
	Denominazione string
	// Name of the person
	Nome string `xml:",omitempty"`
	// Surname of the person
	Cognome string `xml:",omitempty"`
	// Title of the person
	Titolo string `xml:",omitempty"`
	// EORI (Economic Operator Registration and Identification) code
	CodEORI string `xml:",omitempty"`
}

func newCedentePrestatore(inv *bill.Invoice) (*Party, error) {
	s := inv.Supplier

	address, err := newAddress(s)
	if err != nil {
		return nil, err
	}

	return &Party{
		DatiAnagrafici: &DatiAnagrafici{
			IdFiscaleIVA: &TaxID{
				IdPaese:  s.TaxID.Country.String(),
				IdCodice: s.TaxID.Code.String(),
			},
			Anagrafica:    newAnagrafica(s),
			RegimeFiscale: findCodeRegimeFiscale(inv),
		},
		Sede: address,
	}, nil
}

func newCessionarioCommittente(inv *bill.Invoice) (*Party, error) {
	c := inv.Customer

	address, err := newAddress(c)
	if err != nil {
		return nil, err
	}

	da := &DatiAnagrafici{
		Anagrafica: newAnagrafica(c),
	}

	// Apply TaxID or fiscal code. At least one of them is required.
	// FatturaPA only evaluates TaxID if both are present
	if c.TaxID != nil {
		da.IdFiscaleIVA = &TaxID{
			IdPaese:  c.TaxID.Country.String(),
			IdCodice: c.TaxID.Code.String(),
		}
	} else {
		for _, id := range c.Identities {
			if id.Type == "CF" {
				da.CodiceFiscale = id.Code.String()
			}
		}

		if da.CodiceFiscale == "" {
			return nil, errors.New("customer has no TaxID or fiscal code")
		}
	}

	return &Party{
		DatiAnagrafici: da,
		Sede:           address,
	}, nil
}

func newAnagrafica(party *org.Party) *Anagrafica {
	a := Anagrafica{
		Denominazione: party.Name,
	}

	if len(party.People) > 0 {
		name := party.People[0].Name

		a.Nome = name.Given
		a.Cognome = name.Surname
		a.Titolo = name.Prefix
	}

	return &a
}

func findCodeRegimeFiscale(inv *bill.Invoice) string {
	ss := inv.ScenarioSummary()

	return ss.Meta[it.KeyFatturaPARegimeFiscale]
}
