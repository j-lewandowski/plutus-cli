package actions

type Downloader interface {
	Download() error
}

type NBPDownloader struct {
	name   string
	source string
}

func (d NBPDownloader) Download() {}

type YahooFinanceDownloader struct {
	name   string
	source string
}

func (d YahooFinanceDownloader) Download() {}
