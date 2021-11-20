package subscribe

import (
	"testing"

	"github.com/willsem/tfs-go-hw/course_project/internal/config"
)

func TestSubscribeServiceOk(t *testing.T) {
	service, err := New(config.Kraken{
		SocketUrl: "wss://demo-futures.kraken.com/ws/v1",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("success")

	err = service.Subscribe("PI_XBTUS", Candle1m)
	t.Log(err)

	err = service.Subscribe("PI_XBTUSD", Candle1m)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Success subsribed")

	count := 2
	for info := range service.GetChan() {
		t.Log(info)
		count--
		if count == 0 {
			break
		}
	}

	err = service.Unsubscribe("PI_XBTUSD", Candle1m)
	if err != nil {
		t.Fatal(err)
	}

	err = service.Close()
	if err != nil {
		t.Fatal(err)
	}
}
