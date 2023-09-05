package common

import (
	"encoding/binary"
	"math/rand"
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
	Agency        string
	BatchSize     int
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
func (c *Client) sendExact(msg []byte) {
	bytesSent, err := c.conn.Write(msg)

	if err != nil {
		log.Error("action: send_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}
	for bytesSent < len(msg) {
		newBytesSent, err := c.conn.Write(msg[bytesSent:])
		bytesSent += newBytesSent

		if err != nil {
			log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}
	}

}

func (c *Client) readExact(amountOfBytesToRead int) []byte {
	bytesReceived := make([]byte, amountOfBytesToRead)
	numberOfBytesRead, err := c.conn.Read(bytesReceived)

	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return bytesReceived
	}
	for numberOfBytesRead < amountOfBytesToRead {
		newBytesRead, err := c.conn.Read(bytesReceived)
		numberOfBytesRead += newBytesRead

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return bytesReceived
		}
	}

	return bytesReceived
}

func (c *Client) retrieveServerACK() int {
	len_b := c.readExact(2)
	length := binary.BigEndian.Uint16(len_b)
	msg := c.readExact(int(length))

	return int(uint64(uint8(msg[0])))
}

func (c *Client) NotifyBetsSent() int {
	msg := protocol.EncodeBetsSent()
	c.createClientSocket()
	c.sendExact(msg)

	serverACK := c.retrieveServerACK()

	c.conn.Close()

	return serverACK
}

func (c *Client) SendBets(bets []bet.Bet, batch_number int) int {

	msg := protocol.EncodeBets(bets)

	c.createClientSocket()
	c.sendExact(msg)

	serverACK := c.retrieveServerACK()

	c.conn.Close()

	if serverACK == 2 {
		log.Infof("action: apuesta_enviada | result: success | batch_number: %d | bets: %d",
			batch_number,
			len(bets),
		)

		return 1
	} else {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: unkown answer: %d",
			c.config.ID,
			serverACK,
		)
		return 0
	}
}

func (c *Client) AskForWinners(agencyID int) {
	msg := protocol.EncodeWinnersQuery(agencyID)

	for {
		c.createClientSocket()
		// log.Infof("I am in the ask for winners loop")
		// log.Infof("I will write: %x\n", msg)
		c.sendExact(msg)

		header := c.readExact(2)
		length := binary.BigEndian.Uint16(header)
		//log.Infof("ASK FOR WINNERS: HEADER: %x length: %d", header, length)

		payload := c.readExact(int(length))

		switch payload[0] {
		case byte(0):
			c.conn.Close()
			time.Sleep(time.Duration(rand.Intn(100)+1) * time.Millisecond)

		case byte(3):
			c.conn.Close()
			winnersQuantity := protocol.DecodeWinners(payload[1:])
			log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d", winnersQuantity)
			return

		default:
			log.Error("action: winners | result: fail | error: unkown answer")
			c.conn.Close()
		}
	}

}
