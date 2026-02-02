package sync

import (
	"encoding/json"
	"fmt"
	"net/http"
	"plutus-cli/internal/db"
	"sync"
	"time"
)

var ErrNoDeposits = fmt.Errorf("Couldn't perform sync. No deposits found in the database. Please add a deposit first using the 'add' command.")

type Downloader interface {
	SyncData() error
	GetName() string
}

type NBPDownloader struct {
	name       string
	source     string
	repo       *db.Repository
	HttpClient http.Client
}

func NewNBPDownloader(name string, source string, repo *db.Repository) *NBPDownloader {
	return &NBPDownloader{
		name:   name,
		source: source,
		repo:   repo,
		HttpClient: http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (d NBPDownloader) GetName() string {
	return d.name
}

type NBPAPIResult struct {
	Currency string `json:"currency"`
	Code     string `json:"code"`
	Rates    []struct {
		EffectiveDate string      `json:"effectiveDate"`
		Mid           json.Number `json:"mid"`
	} `json:"rates"`
}

func (d NBPDownloader) SyncData() error {
	lastDeposit, err := d.repo.GetFirstDeposit()
	if err != nil {
		return err
	}

	if lastDeposit == (db.Deposit{}) {
		return ErrNoDeposits
	}

	missingDays := DaysUntilToday(lastDeposit.DepositDate)

	data, err := d.DownloadData(missingDays[0], missingDays[len(missingDays)-1])
	if err != nil {
		return err
	}

	userRates := []db.CurrencyRate{}
	for _, rate := range data.Rates {
		userRate := db.CurrencyRate{}
		userRate.From(
			db.NewRateParams{
				Date: rate.EffectiveDate,
				Rate: rate.Mid.String(),
			},
		)

		userRates = append(userRates, userRate)
	}

	populatedRates := PopulateMissingData(missingDays, userRates)

	err = d.repo.AddRates(populatedRates)
	if err != nil {
		return err
	}

	return nil
}

func (d NBPDownloader) DownloadData(startDate time.Time, endDate time.Time) (NBPAPIResult, error) {
	fmt.Println("Downloading NBP data...")

	startDateStr, endDateStr := startDate.Format(time.DateOnly), endDate.Format(time.DateOnly)

	url := fmt.Sprintf("%s/exchangerates/rates/a/eur/%s/%s?format=json", d.source, startDateStr, endDateStr)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return NBPAPIResult{}, fmt.Errorf("Couldn't create request: %v", err)
	}

	req.Header.Set("User-Agent", "PlutusCLI/1.0")

	res, err := d.HttpClient.Do(req)
	if err != nil {
		return NBPAPIResult{}, fmt.Errorf("Network error occurred: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return NBPAPIResult{}, fmt.Errorf("NBP API Error: %d %s", res.StatusCode, res.Status)
	}

	var data NBPAPIResult
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return NBPAPIResult{}, fmt.Errorf("Couldn't parse JSON: %v", err)
	}

	return data, nil
}

type YahooFinanceDownloader struct {
	name       string
	source     string
	ticker     string
	repo       *db.Repository
	HttpClient *http.Client
}

func NewYahooFinanceDownloader(name string, source string, repo *db.Repository) *YahooFinanceDownloader {
	return &YahooFinanceDownloader{
		name:   name,
		source: source,
		repo:   repo,
		ticker: "P500.DE",
		HttpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (d YahooFinanceDownloader) GetName() string {
	return d.name
}

func (d YahooFinanceDownloader) SyncData() error {
	lastDeposit, err := d.repo.GetFirstDeposit()
	if err != nil {
		return err
	}

	if lastDeposit == (db.Deposit{}) {
		return ErrNoDeposits
	}

	missingDays := DaysUntilToday(lastDeposit.DepositDate)

	data, err := d.DownloadData(missingDays[0], missingDays[len(missingDays)-1])
	if err != nil {
		return err
	}

	userIndexPriceList := []db.IndexPrice{}
	result := data.Chart.Result[0]

	if len(result.Indicators.Quote) > 0 {
		quotes := result.Indicators.Quote[0]

		for i, timestamp := range result.Timestamp {
			var avgPrice float64

			hasFullOHLC := i < len(quotes.Open) && quotes.Open[i] != nil &&
				i < len(quotes.High) && quotes.High[i] != nil &&
				i < len(quotes.Low) && quotes.Low[i] != nil &&
				i < len(quotes.Close) && quotes.Close[i] != nil

			hasOpen := i < len(quotes.Open) && quotes.Open[i] != nil

			if hasFullOHLC {
				avgPrice = (*quotes.Open[i] + *quotes.High[i] + *quotes.Low[i] + *quotes.Close[i]) / 4
			} else if hasOpen {
				avgPrice = *quotes.Open[i]
			} else {
				continue
			}

			// Heuristic to handle unit mismatch from Yahoo Finance
			// Data sometimes drops by factor of 100 (e.g. 1150 -> 11.50)
			// Thank you @YahooFinance ;)
			if avgPrice < 200 {
				avgPrice = avgPrice * 100
			}

			userRate := db.IndexPrice{IsReal: true}
			userRate.From(db.NewIndexPriceParams{
				Date:             time.Unix(timestamp, 0).Format("2006-01-02"),
				PriceInEurocents: fmt.Sprintf("%f", avgPrice),
			})

			userIndexPriceList = append(userIndexPriceList, userRate)
		}
	}

	populatedPrices := PopulateMissingData(missingDays, userIndexPriceList)

	err = d.repo.AddIndexPrices(populatedPrices)
	if err != nil {
		return err
	}

	return nil
}

type YahooQuote struct {
	Open  []*float64 `json:"open"`
	High  []*float64 `json:"high"`
	Low   []*float64 `json:"low"`
	Close []*float64 `json:"close"`
}

type YahooChartResponse struct {
	Chart struct {
		Result []struct {
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []YahooQuote `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
	} `json:"chart"`
}

func (d YahooFinanceDownloader) DownloadData(startDate time.Time, endDate time.Time) (YahooChartResponse, error) {
	url := fmt.Sprintf("%s/v8/finance/chart/%s?period1=%d&period2=%d&interval=1d", d.source, d.ticker, startDate.Unix(), endDate.Unix())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return YahooChartResponse{}, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	res, err := d.HttpClient.Do(req)
	if err != nil {
		return YahooChartResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return YahooChartResponse{}, fmt.Errorf("Yahoo API Error: %d %s", res.StatusCode, res.Status)
	}

	var data YahooChartResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return YahooChartResponse{}, err
	}

	if len(data.Chart.Result) == 0 {
		return YahooChartResponse{}, fmt.Errorf("No data found")
	}

	return data, nil
}

func RunSync(repo *db.Repository) error {
	downloaders := []Downloader{
		NewNBPDownloader("NBP Downloader", "https://api.nbp.pl/api", repo),
		NewYahooFinanceDownloader("Yahoo Finance Downloader", "https://query1.finance.yahoo.com", repo),
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(downloaders))

	for _, downloader := range downloaders {
		wg.Add(1)
		go func(d Downloader) {
			defer wg.Done()
			if err := d.SyncData(); err != nil {
				errChan <- fmt.Errorf("%s failed: %w", d.GetName(), err)
			}
		}(downloader)
	}

	wg.Wait()
	close(errChan)

	var combinedError error
	for err := range errChan {
		if combinedError == nil {
			combinedError = err
		} else {
			fmt.Println("Error:", err)
		}
	}

	if combinedError != nil {
		return combinedError
	}

	fmt.Println("Sync completed successfully.")
	return nil
}
