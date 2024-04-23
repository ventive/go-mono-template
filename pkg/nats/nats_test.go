package nats

import (
	"fmt"
	"testing"
	"time"

	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
)

func connectToMockedServer(t *testing.T) Client {
	c := NewClient(Config{
		URL:     fmt.Sprintf("nats://%s:%d", natsserver.DefaultTestOptions.Host, natsserver.DefaultTestOptions.Port),
		Name:    "unit-tests",
		User:    "",
		Pass:    "",
		TLS:     TLSConfig{},
		Options: Options{ReconnectBufSize: -1},
	})
	err := c.Connect()
	if err != nil {
		t.Fatalf("unexpected error when establish the connection to NATS, error = %v", err)
	}
	return c
}

func TestClient_Connect(t *testing.T) {
	// mock NATS server
	s := natsserver.RunDefaultServer()
	defer s.Shutdown()

	type fields struct {
		cfg Config
		nc  *nats.Conn
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"success",
			fields{cfg: Config{
				URL:  fmt.Sprintf("nats://%s:%d", natsserver.DefaultTestOptions.Host, natsserver.DefaultTestOptions.Port),
				Name: "unit-tests",
				User: "",
				Pass: "",
				TLS:  TLSConfig{},
			}},
			false,
		},
		{"success with options",
			fields{cfg: Config{
				URL:  fmt.Sprintf("nats://%s:%d", natsserver.DefaultTestOptions.Host, natsserver.DefaultTestOptions.Port),
				Name: "unit-tests",
				User: "",
				Pass: "",
				TLS:  TLSConfig{},
				Options: Options{
					ReconnectBufSize: 1024,
				},
			}},
			false,
		},
		{"failure",
			fields{cfg: Config{
				URL:  fmt.Sprintf("nats://%s:%d", "invalidhost", natsserver.DefaultTestOptions.Port),
				Name: "unit-tests",
				User: "",
				Pass: "",
				TLS:  TLSConfig{},
			}},
			true,
		},
		{"failure with options",
			fields{cfg: Config{
				URL:  fmt.Sprintf("nats://%s:%d", "invalidhost", natsserver.DefaultTestOptions.Port),
				Name: "unit-tests",
				User: "",
				Pass: "",
				TLS:  TLSConfig{},
				Options: Options{
					ReconnectBufSize: 1024,
				},
			}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &client{
				cfg: tt.fields.cfg,
				nc:  tt.fields.nc,
			}
			if err := client.Connect(); (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_HeadersSupported(t *testing.T) {
	// mock NATS server
	s := natsserver.RunDefaultServer()
	defer s.Shutdown()

	c := connectToMockedServer(t)
	if c.HeadersSupported() != c.GetConn().HeadersSupported() {
		t.Errorf("got = %v, want = %v", c.HeadersSupported(), c.GetConn().HeadersSupported())
	}
}

func TestClient_Subscribe(t *testing.T) {
	// mock NATS server
	s := natsserver.RunDefaultServer()
	defer s.Shutdown()

	c := connectToMockedServer(t)
	// subscribe async
	sub, err := c.Subscribe("unit-tests", func(msg *nats.Msg) {})
	if err != nil {
		t.Fatalf("unepected error = %v", err)
	}
	// subscription should be async
	if sub.Type() != nats.AsyncSubscription {
		t.Errorf("got = %v, want = %v", sub.Type(), nats.AsyncSubscription)
	}
}

func TestClient_QueueSubscribe(t *testing.T) {
	// mock NATS server
	s := natsserver.RunDefaultServer()
	defer s.Shutdown()

	c := connectToMockedServer(t)
	// queue subscribe async
	sub, err := c.QueueSubscribe("unit-tests", "group-unit-tests", func(msg *nats.Msg) {})
	if err != nil {
		t.Fatalf("unepected error = %v", err)
	}
	// subscription should be async
	if sub.Type() != nats.AsyncSubscription {
		t.Errorf("got = %v, want = %v", sub.Type(), nats.AsyncSubscription)
	}
}

func TestClient_SubscribeSync(t *testing.T) {
	// mock NATS server
	s := natsserver.RunDefaultServer()
	defer s.Shutdown()

	c := connectToMockedServer(t)
	// subscribe sync
	sub, err := c.SubscribeSync("unit-tests")
	if err != nil {
		t.Fatalf("unepected error = %v", err)
	}
	// subscription should be sync
	if sub.Type() != nats.SyncSubscription {
		t.Errorf("got = %v, want = %v", sub.Type(), nats.SyncSubscription)
	}
}

func TestClient_QueueSubscribeSync(t *testing.T) {
	// mock NATS server
	s := natsserver.RunDefaultServer()
	defer s.Shutdown()

	c := connectToMockedServer(t)
	// queue subscribe sync
	sub, err := c.QueueSubscribeSync("unit-tests", "group-unit-tests")
	if err != nil {
		t.Fatalf("unepected error = %v", err)
	}
	// subscription should be sync
	if sub.Type() != nats.SyncSubscription {
		t.Errorf("got = %v, want = %v", sub.Type(), nats.SyncSubscription)
	}
}

func TestClient_Unsubscribe(t *testing.T) {
	// mock NATS server
	s := natsserver.RunDefaultServer()
	defer s.Shutdown()

	c := connectToMockedServer(t)
	sub, err := c.Subscribe("unit-tests", func(msg *nats.Msg) {})
	if err != nil {
		t.Fatalf("unepected error = %v", err)
	}
	err = c.Unsubscribe(sub)
	if err != nil {
		t.Fatalf("unepected error = %v", err)
	}
	// @todo: assert somehow Unsubscribe() works
}

func TestClient_Publish(t *testing.T) {
	// mock NATS server
	s := natsserver.RunDefaultServer()
	defer s.Shutdown()

	c := connectToMockedServer(t)
	// publish a slice of bytes
	err := c.Publish("unit-test", []byte("abcd-tests"))
	if err != nil {
		t.Fatalf("unepected error = %v", err)
	}
	// no. of messages should be 1
	if c.GetConn().OutMsgs != 1 {
		t.Errorf("got = %v, want = %v", c.GetConn().OutMsgs, 1)
	}
}

func TestClient_PublishWithRetries(t *testing.T) {
	t.Run("all retries are done", func(t *testing.T) {
		// mock NATS server
		s := natsserver.RunDefaultServer()
		c := connectToMockedServer(t)
		// close connection for retries
		s.Shutdown()
		// wait for server to shutdown
		time.Sleep(500 * time.Millisecond)

		retries := 3
		got, err := c.PublishWithRetries("unit-test", []byte("abcd-tests"), retries)
		if err != nats.ErrReconnectBufExceeded {
			t.Fatalf("unepected error = %v", err)
		}
		if got != retries {
			t.Errorf("got = %v, want = %v", got, retries)
		}
		// no. of messages should be 0
		if c.GetConn().OutMsgs != 0 {
			t.Errorf("got = %v, want = %v", c.GetConn().OutMsgs, 0)
		}
	})

	t.Run("no retries are need", func(t *testing.T) {
		// mock NATS server
		s := natsserver.RunDefaultServer()
		defer s.Shutdown()
		c := connectToMockedServer(t)

		got, err := c.PublishWithRetries("unit-test", []byte("abcd-tests"), 3)
		if err != nil {
			t.Fatalf("unepected error = %v", err)
		}
		if got != 0 {
			t.Errorf("got = %v, want = %v", got, 0)
		}
		// no. of messages should be 0
		if c.GetConn().OutMsgs != 1 {
			t.Errorf("got = %v, want = %v", c.GetConn().OutMsgs, 1)
		}
	})
}

func TestClient_PublishMsg(t *testing.T) {
	// mock NATS server
	s := natsserver.RunDefaultServer()
	defer s.Shutdown()

	c := connectToMockedServer(t)
	// publish a message
	err := c.PublishMsg(&nats.Msg{Subject: "unit-test", Data: []byte("abcd-tests")})
	if err != nil {
		t.Fatalf("unepected error = %v", err)
	}
	// no. of messages should be 1
	if c.GetConn().OutMsgs != 1 {
		t.Errorf("got = %v, want = %v", c.GetConn().OutMsgs, 1)
	}
}

func TestClient_PublishMsgWithRetries(t *testing.T) {
	t.Run("all retries are done", func(t *testing.T) {
		// mock NATS server
		s := natsserver.RunDefaultServer()
		c := connectToMockedServer(t)
		// close connection for retries
		s.Shutdown()
		// wait for server to shutdown
		time.Sleep(500 * time.Millisecond)

		retries := 3
		got, err := c.PublishMsgWithRetries(&nats.Msg{Subject: "unit-test", Data: []byte("abcd-tests")}, retries)
		if err != nats.ErrReconnectBufExceeded {
			t.Fatalf("unepected error = %v", err)
		}
		if got != retries {
			t.Errorf("got = %v, want = %v", got, retries)
		}
		// no. of messages should be 0
		if c.GetConn().OutMsgs != 0 {
			t.Errorf("got = %v, want = %v", c.GetConn().OutMsgs, 0)
		}
	})

	t.Run("no retries are need", func(t *testing.T) {
		// mock NATS server
		s := natsserver.RunDefaultServer()
		defer s.Shutdown()
		c := connectToMockedServer(t)

		retries := 3
		got, err := c.PublishMsgWithRetries(&nats.Msg{Subject: "unit-test", Data: []byte("abcd-tests")}, retries)
		if err != nil {
			t.Fatalf("unepected error = %v", err)
		}
		if got != 0 {
			t.Errorf("got = %v, want = %v", got, 0)
		}
		// no. of messages should be 0
		if c.GetConn().OutMsgs != 1 {
			t.Errorf("got = %v, want = %v", c.GetConn().OutMsgs, 1)
		}
	})
}

func TestClient_Close(t *testing.T) {
	// mock NATS server
	s := natsserver.RunDefaultServer()
	defer s.Shutdown()

	c := connectToMockedServer(t)
	// close connection
	c.Close()
	// should not be connected
	if c.IsConnected() != false {
		t.Errorf("got = %v, want = %v", c.IsConnected(), false)
	}
}
