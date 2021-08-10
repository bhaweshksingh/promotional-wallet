package config

type DataRefresherConfig struct {
	tickerIntervalInSec int
}

func newDataRefresherConfig() DataRefresherConfig {
	return DataRefresherConfig{
		tickerIntervalInSec: getInt("TICKER_INTERVAL_IN_SEC"),
	}
}

func (cc DataRefresherConfig) GetTickerIntervalInSec() int {
	return cc.tickerIntervalInSec
}
