package common

import (
	"net"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/bet"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/protocol"

	log "github.com/sirupsen/logrus"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopLapse     time.Duration
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Fatalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) Send(bet *bet.Bet) int {
	msg := protocol.EncodeBet(bet)
	c.createClientSocket()
	bytesSent, err := c.conn.Write(msg)

	if err != nil {
		log.Error("action: receive_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return 0
	}
	for bytesSent < len(msg) {
		newBytesSent, err := c.conn.Write(msg[bytesSent:])
		bytesSent += newBytesSent

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return 0
		}
	}
	bytesReceived := make([]byte, 1)
	_, err = c.conn.Read(bytesReceived)
	c.conn.Close()

	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return 0
	}
	if bytesReceived[0] == byte(1) {
		log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
			bet.Dni,
			bet.Number,
		)
		return 1
	} else {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return 0
	}

}
