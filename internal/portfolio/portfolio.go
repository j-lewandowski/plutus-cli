package portfolio

import (
	"fmt"
	"plutus-cli/internal/db"
)

const (
	CalculationPrecision = 1000000000000
)

type PortfolioReport struct {
	TotalInvestedInEurocents int64
	CurrentValueInEurocents  int64
	ProfitValueInEurocents   int64
	ProfitPercent            float64

	RateEURtoPLNInGrosz int64
	CurrentValueInGrosz int64
	ProfitValueInGrosz  int64
	HasExchangeRate     bool

	Warnings []string
}

func CalculatePortfolio() (*PortfolioReport, error) {
	overallAmountInEurocents, err := db.GetOverallDepositInEurocents()
	if err != nil {
		return nil, err
	}

	latestIndexPrice, err := db.GetLatestIndexPrice()
	if err != nil {
		return nil, err
	}

	latestExchangeRate, errExchangeRate := db.GetLatestExchangeRate()

	allDeposits, err := db.GetAllDeposits()
	if err != nil {
		return nil, err
	}

	report := &PortfolioReport{
		Warnings: []string{},
	}

	if latestExchangeRate.RateInGrosz == 0 && errExchangeRate != nil {
		report.Warnings = append(report.Warnings, fmt.Sprintf("Warning: Could not fetch exchange rates: %v", errExchangeRate))
	} else if latestExchangeRate.RateInGrosz > 0 {
		report.HasExchangeRate = true
		report.RateEURtoPLNInGrosz = int64(latestExchangeRate.RateInGrosz)
	}

	var totalUnitsOwnedScaled int64

	for _, deposit := range allDeposits {
		var unitsScaled int64

		historicalPrice, err := db.GetIndexPriceByDate(deposit.DepositDate)
		if err != nil {
			report.Warnings = append(report.Warnings, fmt.Sprintf("Warning: Could not find index price for date %s. Skipping calculation for this deposit.", deposit.DepositDate))
			continue
		}

		if historicalPrice.PriceInEurocents == 0 {
			continue
		}

		unitsScaled = (int64(deposit.DepositAmountInEurocents) * CalculationPrecision) / int64(historicalPrice.PriceInEurocents)
		totalUnitsOwnedScaled += unitsScaled
	}

	report.TotalInvestedInEurocents = int64(overallAmountInEurocents)

	report.CurrentValueInEurocents = (totalUnitsOwnedScaled * int64(latestIndexPrice.PriceInEurocents)) / CalculationPrecision

	report.ProfitValueInEurocents = report.CurrentValueInEurocents - report.TotalInvestedInEurocents

	if report.TotalInvestedInEurocents != 0 {
		report.ProfitPercent = (float64(report.ProfitValueInEurocents) / float64(report.TotalInvestedInEurocents)) * 100
	}

	if report.HasExchangeRate {
		report.CurrentValueInGrosz = (report.CurrentValueInEurocents * report.RateEURtoPLNInGrosz) / 100
		report.ProfitValueInGrosz = (report.ProfitValueInEurocents * report.RateEURtoPLNInGrosz) / 100
	}

	return report, nil
}
