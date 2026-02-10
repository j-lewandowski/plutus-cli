package portfolio

import (
	"plutus-cli/internal/db"
)

func GetChartData(repo *db.Repository) ([]db.ChartPoint, error) {
	chartPoints, err := repo.GetChartPoints()

	if err != nil {
		return []db.ChartPoint{}, nil
	}

	return chartPoints, nil
}
