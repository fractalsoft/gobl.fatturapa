package fatturapa

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
)

// datiRitenuta contains all data related to the retained taxes.
type datiRitenuta struct {
	TipoRitenuta     string
	ImportoRitenuta  string
	AliquotaRitenuta string
	CausalePagamento string
}

func extractRetainedTaxes(inv *bill.Invoice) ([]*datiRitenuta, error) {
	var dr []*datiRitenuta
	var catTotals []*tax.CategoryTotal

	// First we need to find all the retained tax categories from Totals
	for _, catTotal := range inv.Totals.Taxes.Categories {
		if catTotal.Retained {
			catTotals = append(catTotals, catTotal)
		}
	}

	for _, catTotal := range catTotals {
		for _, rateTotal := range catTotal.Rates {
			rate := formatPercentage(rateTotal.Percent)
			amount := formatAmount(&rateTotal.Amount)
			codeTR, err := findCodeTipoRitenuta(catTotal.Code)
			if err != nil {
				return nil, err
			}
			codeCP, err := findCodeCausalePagamento(catTotal.Code, rateTotal.Key)
			if err != nil {
				return nil, err
			}

			dr = append(dr, &datiRitenuta{
				TipoRitenuta:     codeTR,
				ImportoRitenuta:  amount,
				AliquotaRitenuta: rate,
				CausalePagamento: codeCP,
			})
		}
	}

	return dr, nil
}

func findCodeTipoRitenuta(cat cbc.Code) (string, error) {
	taxCategory := regime.Category(cat)

	code := taxCategory.Codes[it.KeyFatturaPATipoRitenuta]

	if code == "" {
		return "", fmt.Errorf("could not find TipoRitenuta code for tax category %s", cat)
	}

	return code.String(), nil
}

func findCodeCausalePagamento(cat cbc.Code, rateKey cbc.Key) (string, error) {
	taxCategory := regime.Category(cat)

	for _, rate := range taxCategory.Rates {
		if rate.Key == rateKey {
			code := rate.Codes[it.KeyFatturaPACausalePagamento]

			if code == "" {
				return "", fmt.Errorf(
					"could not find CausalePagamento code for tax category %s and rate %s",
					cat,
					rateKey,
				)
			}

			return code.String(), nil
		}
	}

	return "", fmt.Errorf(
		"could not find CausalePagamento code for tax category %s and rate %s",
		cat,
		rateKey,
	)
}
