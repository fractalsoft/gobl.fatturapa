package fatturapa

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
)

const (
	statoLiquidazioneDefault   = "LN"
	regimeFiscaleDefault       = "RF01"
	euCitizenTaxCodeDefault    = "0000000"
	nonEUCitizenTaxCodeDefault = "99999999999"
)

// Supplier is the party that issues the invoice
type Supplier struct {
	DatiAnagrafici *DatiAnagrafici
	Sede           *Address
	IscrizioneREA  *IscrizioneREA `xml:",omitempty"`
}

// Customer is the party that receives the invoice
type Customer struct {
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

// IscrizioneREA contains information related to the company registration details (REA)
type IscrizioneREA struct {
	// Initials of the province where the company's Registry Office is located
	Ufficio string
	// Company's REA registration number
	NumeroREA string
	// Company's share capital
	CapitaleSociale string `xml:",omitempty"`
	// Indication of whether the Company is in liquidation or not.
	// Possible values: LS (in liquidation), LN (not in liquidation)
	StatoLiquidazione string
}

func newCedentePrestatore(inv *bill.Invoice) (*Supplier, error) {
	s := inv.Supplier

	address, err := newAddress(s)
	if err != nil {
		return nil, err
	}

	return &Supplier{
		DatiAnagrafici: &DatiAnagrafici{
			IdFiscaleIVA: &TaxID{
				IdPaese:  s.TaxID.Country.String(),
				IdCodice: s.TaxID.Code.String(),
			},
			Anagrafica:    newAnagrafica(s),
			RegimeFiscale: findCodeRegimeFiscale(inv),
		},
		Sede:          address,
		IscrizioneREA: newIscrizioneREA(s),
	}, nil
}

func newCessionarioCommittente(inv *bill.Invoice) (*Customer, error) {
	c := inv.Customer

	address, err := newAddress(c)
	if err != nil {
		return nil, err
	}

	da := &DatiAnagrafici{
		Anagrafica: newAnagrafica(c),
	}

	if c.TaxID == nil || c.TaxID.Country == "" {
		return nil, fmt.Errorf(
			"missing customer TaxID. at least the country code " +
				"must be present under Invoice.Customer.TaxID")
	}

	if isCodiceFiscale(c.TaxID) {
		da.CodiceFiscale = c.TaxID.Code.String()
	} else if isEUCountry(c.TaxID.Country) {
		da.IdFiscaleIVA = customerFiscaleIVA(c.TaxID, euCitizenTaxCodeDefault)
	} else {
		da.IdFiscaleIVA = customerFiscaleIVA(c.TaxID, nonEUCitizenTaxCodeDefault)
	}

	return &Customer{
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

	regimeFiscale := ss.Meta[it.KeyFatturaPARegimeFiscale]

	if regimeFiscale != "" {
		return regimeFiscale
	}

	return regimeFiscaleDefault
}

func customerFiscaleIVA(taxID *tax.Identity, fallBack string) *TaxID {
	idCodice := taxID.Code.String()

	if idCodice == "" {
		idCodice = fallBack
	}

	return &TaxID{
		IdPaese:  taxID.Country.String(),
		IdCodice: idCodice,
	}
}

func newIscrizioneREA(supplier *org.Party) *IscrizioneREA {
	if supplier.Registration == nil {
		return nil
	}

	capital := supplier.Registration.Capital
	var capitalFormatted string

	if capital == nil {
		capitalFormatted = ""
	} else {
		capitalFormatted = capital.Rescale(2).String()
	}

	return &IscrizioneREA{
		Ufficio:           supplier.Registration.Office,
		NumeroREA:         supplier.Registration.Entry,
		CapitaleSociale:   capitalFormatted,
		StatoLiquidazione: statoLiquidazioneDefault,
	}
}

func isCodiceFiscale(taxID *tax.Identity) bool {
	if taxID.Country != l10n.IT {
		return false
	}

	return len(taxID.Code.String()) == 16
}
