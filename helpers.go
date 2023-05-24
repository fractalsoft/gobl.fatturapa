package fatturapa

import "github.com/invopop/gobl/num"

func formatPercentage(p *num.Percentage) string {
	if p == nil {
		return ""
	}

	return p.Rescale(4).StringWithoutSymbol()
}

func formatAmount(a *num.Amount) string {
	if a == nil {
		return ""
	}

	return a.Rescale(2).String()
}