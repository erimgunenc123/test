package tickers

import (
	"encoding/json"
	btcturk_constants "genericAPI/btcturk_connector/constants"
	"genericAPI/internal/common/http_utils"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// ticker has 600 rate limit per minute
type tickerService struct {
	tickers     map[string]*ticker // symbol -> ticker
	tickersLock sync.Mutex
}

var TickerService *tickerService

type ticker struct {
	symbol   string
	fetcher  func() (*http.Response, error)
	interval time.Duration
	subs     []chan *Tick
	subsLock sync.Mutex
}

func (t *ticker) listen() {
	for {
		time.Sleep(t.interval)
		res, err := t.fetcher()
		if err != nil {
			log.Printf("%s ticker failed with error: %s", t.symbol, err.Error())
			continue
		}
		if res.StatusCode == 200 {
			responseBytes, _ := io.ReadAll(res.Body)
			var tick Tick
			err = json.Unmarshal(responseBytes, &tick)
			if err != nil {
				log.Printf("%s ticker failed unmarshalling error: %s", t.symbol, err.Error())
				continue
			}
			t.subsLock.Lock()
			for _, sub := range t.subs {
				go func(s chan *Tick) {
					s <- &tick
				}(sub)
			}
			t.subsLock.Unlock()
		}
	}
}

func (t *ticker) addSub(ch chan *Tick) {
	t.subsLock.Lock()
	defer t.subsLock.Unlock()
	t.subs = append(t.subs, ch)
}

func InitTickerService() {
	TickerService = &tickerService{
		tickers: make(map[string]*ticker),
	}
}

func NewTicker(symbol string, sub chan *Tick) error {
	if ticker, ok := TickerService.tickers[symbol]; ok {
		if sub != nil {
			ticker.addSub(sub)
		}
		return nil
	}

	cl, err := http_utils.GetRequestClosure(btcturk_constants.BaseUrl+"ticker", nil, map[string]string{"pairSymbol": symbol})
	if err != nil {
		return err
	}
	t := ticker{
		symbol:   symbol,
		fetcher:  cl,
		subs:     []chan *Tick{},
		subsLock: sync.Mutex{},
		interval: btcturk_constants.TickerDefaultInterval,
	}
	if sub != nil {
		t.addSub(sub)
	}
	TickerService.tickersLock.Lock()
	defer TickerService.tickersLock.Unlock()
	TickerService.tickers[symbol] = &t
	go t.listen()
	return nil
}
