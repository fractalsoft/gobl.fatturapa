package fatturapa

import (
	"strconv"
	"strings"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

type DatiBeniServizi struct {
	DettaglioLinee []DettaglioLinee
	DatiRiepilogo  []DatiRiepilogo
}

type DettaglioLinee struct {
	NumeroLinea    string
	Descrizione    string
	Quantita       string
	PrezzoUnitario string
	PrezzoTotale   string
	AliquotaIVA    string
}

type DatiRiepilogo struct {
	AliquotaIVA       string
	ImponibileImporto string
	Imposta           string
	EsigibilitaIVA    string
}

func newDatiBeniServizi(inv bill.Invoice) DatiBeniServizi {
	return DatiBeniServizi{
		DettaglioLinee: newDettaglioLinee(inv),
		DatiRiepilogo:  newDatiRiepilogo(inv),
	}
}

func newDettaglioLinee(inv bill.Invoice) []DettaglioLinee {
	var dl []DettaglioLinee

	for _, line := range inv.Lines {
		desc := ""
		vatRate := ""
		l := len(line.Notes)

		for i, note := range line.Notes {
			if note.Key == cbc.NoteKeyGoods {
				desc += strings.TrimSpace(note.Text)

				if desc[len(desc)-1] != '.' {
					desc += "."
				}

				if i < l-1 {
					desc += " "
				}
			}
		}

		for _, tax := range line.Taxes {
			if tax.Category == common.TaxCategoryVAT {
				vatRate = tax.Percent.String()
				break
			}
		}

		dl = append(dl, DettaglioLinee{
			NumeroLinea:    strconv.Itoa(line.Index),
			Descrizione:    desc,
			Quantita:       line.Quantity.String(),
			PrezzoUnitario: line.Item.Price.String(),
			PrezzoTotale:   line.Sum.String(),
			AliquotaIVA:    vatRate,
		})
	}

	return dl
}

func newDatiRiepilogo(inv bill.Invoice) []DatiRiepilogo {
	var dr []DatiRiepilogo
	var vatRates []*tax.RateTotal

	for _, cat := range inv.Totals.Taxes.Categories {
		if cat.Code == common.TaxCategoryVAT {
			vatRates = cat.Rates
		}
	}

	for _, rate := range vatRates {
		dr = append(dr, DatiRiepilogo{
			AliquotaIVA:       rate.Percent.String(),
			ImponibileImporto: rate.Base.String(),
			Imposta:           rate.Amount.String(),
			EsigibilitaIVA:    "I", // TODO
		})
	}

	return dr
}
