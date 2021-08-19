package client

import (
	"encoding/json"
	"fmt"
	"github.com/kovey/rpc-go/protocol"
	"github.com/kovey/rpc-go/rpc"
	"github.com/kovey/server-go/server"
	"net"
)

type Client struct {
	host        string
	port        int
	connection  *server.Connection
	config      *server.Config
	isConnected bool
}

func NewClient(host string, port int, config *server.Config) *Client {
	return &Client{host: host, port: port, config: config, isConnected: false}
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.host, c.port))
	if err != nil {
		return err
	}

	c.connection = server.NewConnection(1, conn, c.config)
	c.isConnected = true
	return nil
}

func (c *Client) Recv() (*protocol.ClientResponse, error) {
	for {
		buf, _, err := c.connection.Read()
		if err != nil {
			return nil, err
		}

		if len(buf) == 0 {
			continue
		}

		response := protocol.NewClientResponse()
		err = json.Unmarshal(buf, response)
		if err != nil {
			return nil, err
		}

		return response, nil
	}
}

func (c *Client) Send(request protocol.Request) error {
	buf, err := json.Marshal(request)
	if err != nil {
		return err
	}

	bLen := len(buf)
	lBuf := rpc.Int32ToBytes(int32(bLen))
	buf = append(lBuf, buf...)
	_, err = c.connection.Write(buf)
	return err
}

func (c *Client) Close() error {
	c.isConnected = false
	return c.connection.Close()
}
