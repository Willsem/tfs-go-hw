package subscribe

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/willsem/tfs-go-hw/course_project/internal/config"
	"golang.org/x/net/websocket"
)

const (
	pingSeconds = 59
)

type KrakenSubscribeService struct {
	websocket    *websocket.Conn
	mutex        *sync.Mutex
	tickerChan   chan TickerInfo
	cancelCtx    func()
	countTickers int
}

func New(config config.Kraken) (*KrakenSubscribeService, error) {
	ws, err := websocket.Dial(config.SocketUrl, "", "/")
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	err = websocket.JSON.Receive(ws, &response)
	if err != nil {
		return nil, err
	}

	responseEvent, ok := response["event"]
	if !ok {
		return nil, fmt.Errorf("connection error")
	}

	if responseEvent.(string) != "info" {
		return nil, fmt.Errorf("connection error")
	}

	ctx, cancel := context.WithCancel(context.Background())

	service := &KrakenSubscribeService{
		websocket:    ws,
		mutex:        &sync.Mutex{},
		tickerChan:   make(chan TickerInfo),
		cancelCtx:    cancel,
		countTickers: 0,
	}

	go service.listenSocket(ctx)
	go service.pingSocket(ctx)

	return service, nil
}

func (service *KrakenSubscribeService) Close() error {
	service.mutex.Lock()
	defer service.mutex.Unlock()
	service.cancelCtx()
	close(service.tickerChan)
	return service.websocket.Close()
}

func (service *KrakenSubscribeService) GetChan() <-chan TickerInfo {
	return service.tickerChan
}

func (service *KrakenSubscribeService) Subscribe(ticker string) error {
	event := Event{
		Event:      "subscribe",
		Feed:       "ticker",
		ProductIds: []string{ticker},
	}
	var response map[string]interface{}

	service.mutex.Lock()
	websocket.JSON.Send(service.websocket, event)
	err := websocket.JSON.Receive(service.websocket, &response)
	service.mutex.Unlock()

	if err != nil {
		return err
	}

	responseEvent, ok := response["event"]
	if !ok {
		return fmt.Errorf("unknown response")
	}

	if responseEvent.(string) != "subscribed" {
		message, ok := response["message"].(string)
		if !ok {
			return fmt.Errorf("unknown error response:" + responseEvent.(string))
		}

		return fmt.Errorf(message)
	}

	service.countTickers++
	return nil
}

func (service *KrakenSubscribeService) Unsubscribe(ticker string) error {
	event := Event{
		Event:      "unsubscribe",
		Feed:       "ticker",
		ProductIds: []string{ticker},
	}
	var response map[string]interface{}

	service.mutex.Lock()
	websocket.JSON.Send(service.websocket, event)
	err := websocket.JSON.Receive(service.websocket, &response)
	service.mutex.Unlock()

	if err != nil {
		return err
	}

	responseEvent, ok := response["event"]
	if !ok {
		return fmt.Errorf("unknown response")
	}

	if responseEvent.(string) != "unsubscribed" {
		message, ok := response["message"]
		if !ok {
			return fmt.Errorf("unknown error response:" + responseEvent.(string))
		}

		return fmt.Errorf(message.(string))
	}

	service.countTickers--
	return nil
}

func (service *KrakenSubscribeService) listenSocket(ctx context.Context) {
	data := TickerInfo{}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if service.countTickers > 0 {
				service.mutex.Lock()
				err := websocket.JSON.Receive(service.websocket, &data)
				service.mutex.Unlock()

				if err != nil {
					continue
				}

				service.tickerChan <- data
			}
		}
	}
}

func (service *KrakenSubscribeService) pingSocket(ctx context.Context) {
	nothing := struct{}{}

	timeout, cancel := context.WithTimeout(context.Background(), time.Duration(pingSeconds*time.Second))

	for {
		select {
		case <-ctx.Done():
			cancel()
			return
		case <-timeout.Done():
			timeout, cancel = context.WithTimeout(context.Background(), time.Duration(pingSeconds*time.Second))
			service.mutex.Lock()
			websocket.JSON.Send(service.websocket, nothing)
			websocket.Message.Receive(service.websocket, &nothing)
			service.mutex.Unlock()
		}
	}
}
