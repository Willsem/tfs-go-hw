package subscribe

import (
	"testing"
	"time"

	"github.com/willsem/tfs-go-hw/course_project/internal/config"
)

func TestSubscribeServiceOk(t *testing.T) {
	service, err := New(config.Kraken{
		SocketUrl: "wss://demo-futures.kraken.com/ws/v1",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = service.Subscribe("PI_XBTUS")
	t.Log(err)

	err = service.Subscribe("PI_XBTUSD")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for info := range service.GetChan() {
			t.Log(info)
		}
	}()

	time.Sleep(2 * time.Second)
	err = service.Unsubscribe("PI_XBTUSD")
	if err != nil {
		t.Fatal(err)
	}

	err = service.Close()
	if err != nil {
		t.Fatal(err)
	}
}
