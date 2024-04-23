package nats

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
)

// Msg alias this here to avoid importing the nats-go package when wanting to use it
// in other pkgs
type Msg = nats.Msg

// Subscription alias here - same reason as for Msg
type Subscription = nats.Subscription

// Config represents the NATS client configuration
type Config struct {
	URL     string
	Name    string
	User    string
	Pass    string
	TLS     TLSConfig
	Options Options
}

// TLSConfig represents the TLS part of NATS client configuration
type TLSConfig struct {
	Enabled bool
	Cert    string
	Key     string
	CA      string
}

// Options represents the NATS client options configuration
type Options struct {
	// ReconnectBufSize specifies the buffer size of messages kept while busy reconnecting.
	ReconnectBufSize int
}

type client struct {
	cfg Config
	nc  *nats.Conn
}

// Client is a custom wrapper on top of nats-go pkg
type Client interface {
	Connect() error
	GetConn() *nats.Conn
	IsConnected() bool
	HeadersSupported() bool
	DisconnectErrHandler(_ *nats.Conn, err error)
	ReconnectHandler(_ *nats.Conn)
	Subscribe(queue string, handler nats.MsgHandler) (*nats.Subscription, error)
	Unsubscribe(sub *nats.Subscription) error
	QueueSubscribe(queue, name string, handler nats.MsgHandler) (*nats.Subscription, error)
	SubscribeSync(queue string) (*nats.Subscription, error)
	QueueSubscribeSync(queue, name string) (*nats.Subscription, error)
	Publish(subj string, data []byte) error
	PublishWithRetries(subject string, data []byte, retries int) (int, error)
	PublishMsg(msg *Msg) error
	PublishMsgWithRetries(msg *Msg, retries int) (int, error)
	RequestMsg(msg *Msg, timeout time.Duration) (*Msg, error)
	RequestMsgWithRetries(msg *Msg, timeout time.Duration, retries int) (*Msg, int, error)
	Close()
}

// NewClient creates a new NATS client
func NewClient(cfg Config) Client {
	return &client{cfg: cfg}
}

// Connect starts a network connection to the NATS server
func (client *client) Connect() error {
	var err error
	options := []nats.Option{
		nats.Name(client.cfg.Name),
		nats.Timeout(10 * time.Second),
		nats.DisconnectErrHandler(client.DisconnectErrHandler),
		nats.ReconnectHandler(client.ReconnectHandler),
		nats.MaxReconnects(-1),
		// nats.PingInterval(20*time.Second),
		// nats.MaxPingsOutstanding(5),
		// nats.NoEcho(), // Do not receive published messages back even if subscribed
		// nats.NoReconnect(), // Do not reconnect on network failure
	}
	// set reconnect buffer size
	if client.cfg.Options.ReconnectBufSize == 0 {
		// set a default 5MB buffer size
		options = append(options, nats.ReconnectBufSize(5*1024*1024))
	} else {
		// set the value defined in config
		// a negative value will represent buffer size of 0
		options = append(options, nats.ReconnectBufSize(client.cfg.Options.ReconnectBufSize))
	}
	if client.cfg.User != "" && client.cfg.Pass != "" {
		options = append(options, nats.UserInfo(client.cfg.User, client.cfg.Pass))
	}
	if client.cfg.TLS.Enabled {
		options = append(
			options,
			nats.ClientCert(client.cfg.TLS.Cert, client.cfg.TLS.Key),
			nats.RootCAs(client.cfg.TLS.CA),
		)
	}
	client.nc, err = nats.Connect(client.cfg.URL, options...)
	return err
}

// GetConn returns current NATS connection
func (client client) GetConn() *nats.Conn {
	return client.nc
}

// IsConnected checks if connection is was made
func (client client) IsConnected() bool {
	if client.nc == nil {
		return false
	}

	return client.nc.IsConnected()
}

// HeadersSupported returns if NATS server supports headers
func (client client) HeadersSupported() bool {
	return client.nc.HeadersSupported()
}

func (client client) DisconnectErrHandler(_ *nats.Conn, err error) {
	// handle disconnect error event
	// @todo: What should we add here?
}

func (client client) ReconnectHandler(_ *nats.Conn) {
	// handle reconnect event
	// @todo: What should we add here?
}

// Subscribe will subscribe async on the given queue
func (client client) Subscribe(queue string, handler nats.MsgHandler) (*nats.Subscription, error) {
	return client.nc.Subscribe(queue, handler)
}

// QueueSubscribe returns an async queue subscriber on the give subject (queue)
func (client client) QueueSubscribe(queue, name string, handler nats.MsgHandler) (*nats.Subscription, error) {
	return client.nc.QueueSubscribe(queue, name, handler)
}

// SubscribeSync will subscribe sync on the given queue
func (client client) SubscribeSync(queue string) (*nats.Subscription, error) {
	return client.nc.SubscribeSync(queue)
}

// QueueSubscribeSync returns a sync queue subscriber on the give subject (queue)
func (client client) QueueSubscribeSync(queue, name string) (*nats.Subscription, error) {
	return client.nc.QueueSubscribeSync(queue, name)
}

// Unsubscribe from a given subject
func (client client) Unsubscribe(subscription *nats.Subscription) error {
	return subscription.Unsubscribe()
}

// Publish publishes a slice of bytes to the give subject (queue)
func (client client) Publish(subject string, data []byte) error {
	return client.nc.Publish(subject, data)
}

// PublishWithRetries publishes a slice of bytes to the give subject (queue) using retries with exponential backoff
// The exponential backoff is calculated based on the number of retries.
// E.g. retries = 3
// - main call will be made with a delay of 0 seconds
// - retry 1 will be made with a delay of 1 second
// - retry 2 will be made with a delay of 2 seconds
func (client client) PublishWithRetries(subject string, data []byte, retries int) (int, error) {
	var err error
	i := 0
	for {
		err = client.Publish(subject, data)
		if err == nil {
			return i, nil
		}
		i++
		if i >= retries {
			return i, err
		}
		time.Sleep(time.Duration(i) * time.Second)
	}
}

// PublishMsg publishes a Msg structure
func (client client) PublishMsg(msg *Msg) error {
	return client.nc.PublishMsg(msg)
}

// PublishMsgWithRetries publishes a Msg structure using retries
// The exponential backoff is calculated based on the number of retries.
// E.g. retries = 3
// - main call will be made with a delay of 0 seconds
// - retry 1 will be made with a delay of 1 second
// - retry 2 will be made with a delay of 2 seconds
func (client client) PublishMsgWithRetries(msg *Msg, retries int) (int, error) {
	var err error
	i := 0
	for {
		err = client.PublishMsg(msg)
		if err == nil {
			return i, nil
		}
		i++
		if i >= retries {
			return i, err
		}
		time.Sleep(time.Duration(i) * time.Second)
	}
}

// RequestMsg wrapper for RequestMsg
func (client client) RequestMsg(msg *Msg, timeout time.Duration) (*Msg, error) {
	return client.nc.RequestMsg(msg, timeout)
}

// RequestMsgWithRetries wrapper for RequestMsg using retries
// The exponential backoff is calculated based on the number of retries.
// E.g. retries = 3
// - main call will be made with a delay of 0 seconds
// - retry 1 will be made with a delay of 1 second
// - retry 2 will be made with a delay of 2 seconds
func (client client) RequestMsgWithRetries(msg *Msg, timeout time.Duration, retries int) (*Msg, int, error) {
	var err error
	i := 0
	for {
		var responseMsg *Msg
		responseMsg, err = client.RequestMsg(msg, timeout)
		if err == nil {
			return responseMsg, i, nil
		}
		i++
		if i >= retries {
			return nil, i, err
		}
		time.Sleep(time.Duration(i) * time.Second)
	}
}

// Close terminates the connection to the NATS server and releases all blocking calls
func (client client) Close() {
	_ = client.nc.Flush()
	client.nc.Close()
}

// ReadMsg gets next message from subscription
func ReadMsg(ctx context.Context, sub *nats.Subscription) (*Msg, error) {
	return sub.NextMsgWithContext(ctx)
}

// NewMsg wrapper for NewMsg
func NewMsg(subject string) *Msg {
	return nats.NewMsg(subject)
}

// NewMsgWithHeaders wrapper for NewMsg and setting headers afterwards
func NewMsgWithHeaders(subject string, headers map[string]string) *Msg {
	msg := nats.NewMsg(subject)
	for k, v := range headers {
		msg.Header.Set(k, v)
	}

	return msg
}
