package porla

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/autobrr/autobrr/pkg/jsonrpc"
)

var (
	DefaultTimeout = 60 * time.Second
)

type Client struct {
	Name      string
	Hostname  string
	cfg       Config
	rpcClient jsonrpc.Client
	http      *http.Client
	timeout   time.Duration

	log *log.Logger
}

type Config struct {
	Hostname  string
	AuthToken string

	// TLS skip cert validation
	TLSSkipVerify bool

	// HTTP Basic auth username
	BasicUser string

	// HTTP Basic auth password
	BasicPass string

	Timeout int
	Log     *log.Logger
}

func NewClient(cfg Config) *Client {
	c := &Client{
		cfg:     cfg,
		log:     log.New(io.Discard, "", log.LstdFlags),
		timeout: DefaultTimeout,
	}

	// override logger if we pass one
	if cfg.Log != nil {
		c.log = cfg.Log
	}

	if cfg.Timeout > 0 {
		c.timeout = time.Duration(cfg.Timeout) * time.Second
	}

	c.http = &http.Client{
		Timeout: c.timeout,
	}

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	if cfg.TLSSkipVerify {
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	httpClient := &http.Client{
		Timeout:   c.timeout,
		Transport: customTransport,
	}

	c.rpcClient = jsonrpc.NewClientWithOpts(cfg.Hostname+"/api/v1/jsonrpc", &jsonrpc.ClientOpts{
		Headers: map[string]string{
			"X-Porla-Token": cfg.AuthToken,
		},
		HTTPClient: httpClient,
		BasicUser:  cfg.BasicUser,
		BasicPass:  cfg.BasicPass,
	})

	return c
}
