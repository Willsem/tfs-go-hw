package subscribe

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/willsem/tfs-go-hw/course_project/internal/config"
	"github.com/willsem/tfs-go-hw/course_project/internal/domain"
)

const (
	Candle1m candleType = "candles_trade_1m"
	Candle5m candleType = "candles_trade_5m"
	Candle1h candleType = "candles_trade_1h"

	sizeChan = 100
)

type KrakenSubscribeService struct {
	ws     *websocket.Conn
	config config.Kraken

	writeMutex *sync.Mutex
	tickers    chan domain.TickerInfo
	alerts     chan string
	success    chan struct{}

	cancelListen func()
}

func connectWebsocket(config config.Kraken) (*websocket.Conn, error) {
	ws, _, err := websocket.DefaultDialer.Dial(config.SocketUrl, nil)
	if err != nil {
		return nil, err
	}
	ws.SetReadDeadline(time.Time{})

	heartbeatEvent := event{
		Event: "subscribe",
		Feed:  "heartbeat",
	}
	ws.WriteJSON(heartbeatEvent)

	return ws, nil
}

func New(config config.Kraken) (*KrakenSubscribeService, error) {
	ws, err := connectWebsocket(config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	service := &KrakenSubscribeService{
		ws:     ws,
		config: config,

		writeMutex: &sync.Mutex{},
		tickers:    make(chan domain.TickerInfo, sizeChan),
		alerts:     make(chan string, sizeChan),
		success:    make(chan struct{}, sizeChan),

		cancelListen: cancel,
	}

	go service.listenSocket(ctx)

	select {
	case <-service.success:
		return service, nil
	case message := <-service.alerts:
		return nil, errors.New(message)
	}
}

func (service *KrakenSubscribeService) Close() error {
	service.cancelListen()
	close(service.tickers)
	close(service.alerts)
	close(service.success)
	return service.ws.Close()
}

func (service *KrakenSubscribeService) GetChan() <-chan domain.TickerInfo {
	return service.tickers
}

func (service *KrakenSubscribeService) Subscribe(ticker string, candle candleType) error {
	service.writeMutex.Lock()
	defer service.writeMutex.Unlock()

	return service.sendEvent(event{
		Event:      "subscribe",
		Feed:       string(candle),
		ProductIds: []string{ticker},
	})
}

func (service *KrakenSubscribeService) Unsubscribe(ticker string, candle candleType) error {
	service.writeMutex.Lock()
	defer service.writeMutex.Unlock()

	return service.sendEvent(event{
		Event:      "unsubscribe",
		Feed:       string(candle),
		ProductIds: []string{ticker},
	})
}

func (service *KrakenSubscribeService) sendEvent(event event) error {
	err := service.ws.WriteJSON(event)
	if err != nil {
		return err
	}

	select {
	case <-service.success:
		return nil
	case message := <-service.alerts:
		return errors.New(message)
	}
}

func (service *KrakenSubscribeService) listenSocket(ctx context.Context) {
	var resp event

	for {
		resp = event{}

		select {
		case <-ctx.Done():
			return
		default:
			err := service.ws.ReadJSON(&resp)
			if err != nil {
				service.ws, _ = connectWebsocket(service.config)
				select {
				case <-service.success:
				case <-service.alerts:
				}
			}
			service.writeResponse(resp)
		}
	}
}

func (service *KrakenSubscribeService) writeResponse(response event) {
	switch {
	case response.Event == "subscribed" || response.Event == "unsubscribed":
		service.success <- struct{}{}

	case response.Event == "alert":
		service.alerts <- response.Message

	case response.Feed == string(Candle1m) ||
		response.Feed == string(Candle5m) ||
		response.Feed == string(Candle1h) ||
		response.Feed == string(Candle1m)+"_snapshot" ||
		response.Feed == string(Candle5m)+"_snapshot" ||
		response.Feed == string(Candle1h)+"_snapshot":
		tickerInfo := domain.TickerInfo{
			Feed:      response.Feed,
			ProductId: response.ProductId,
			Candle:    response.Candle,
		}
		service.tickers <- tickerInfo
	}
}
