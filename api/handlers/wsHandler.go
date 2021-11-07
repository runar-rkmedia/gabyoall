package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/runar-rkmedia/gabyoall/logger"
	"github.com/runar-rkmedia/gabyoall/utils"
)

type WsOptions struct {
	PingInterval,
	WriteTimeout time.Duration
}

type Msg struct {
	Kind     string      `json:"kind"`
	Variant  string      `json:"variant,omitempty"`
	Contents interface{} `json:"contents,omitempty"`
}

func NewWsHandler(l logger.AppLogger, sendChannel chan Msg, options WsOptions) http.HandlerFunc {
	if options.WriteTimeout == 0 {
		options.WriteTimeout = 30 * time.Second
	}
	if options.PingInterval == 0 {
		options.PingInterval = 7 * time.Second
	}
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	clients := ClientList{
		map[string]*Client{},
		&options,
		l,
		sync.RWMutex{},
	}
	go clients.Writer(sendChannel)

	return func(w http.ResponseWriter, r *http.Request) {

		l.Debug().Msg("Client upgrading")
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			l.Error().Err(err).Msg("failed to upgrade to ws-connection")
			return
		}
		l.Debug().Msg("Client subscribing")

		clients.Subscribe(&Client{
			conn:       ws,
			options:    &options,
			l:          l,
			clientLIst: &clients,
		})
	}
}

type ClientList struct {
	clients map[string]*Client
	options *WsOptions
	l       logger.AppLogger
	sync.RWMutex
}

func (cl *ClientList) Subscribe(c *Client) {
	if c.id == "" {
		id, idErr := utils.ForceCreateUniqueId()
		if idErr != nil {
			cl.l.Error().Err(idErr).Msg("Error while creating unique id.")
		}
		c.id = id
	}
	cl.Lock()
	defer cl.Unlock()
	if cl.l.HasDebug() {
		cl.l.Debug().Str("ID", c.id).Msg("Client subscribed")
	}
	cl.clients[c.id] = c
}

func (cl *ClientList) Unsubscribe(id string) {
	if cl.l.HasDebug() {
		cl.l.Debug().Str("ID", id).Msg("Client unsubscribed")
	}
	cl.Lock()
	delete(cl.clients, id)
	cl.Unlock()
}

func (cl *ClientList) Writer(ch chan Msg) {
	ticker := time.NewTicker(cl.options.PingInterval)
	debug := cl.l.HasDebug()
	for {
		select {
		case j := <-ch:
			cl.RLock()
			length := len(cl.clients)
			if length == 0 {
				cl.l.Debug().Msg("No clients are listening at the moment")
				cl.RUnlock()
				continue
			}
			json, err := json.Marshal(j)
			if err != nil {
				cl.l.Error().Err(err).Msg("Failed to json-marshal message")
				cl.RUnlock()
				continue
			}
			if debug {
				cl.l.Debug().Int("count", length).Str("kind", j.Kind).Str("variant", j.Variant).Msg("Sending message to clients")
			}
			wg := sync.WaitGroup{}
			wg.Add(length)

			for _, client := range cl.clients {
				go func(c *Client) {
					c.write(websocket.TextMessage, json)
					wg.Done()
				}(client)
			}
			cl.RUnlock()
			wg.Wait()

		case <-ticker.C:
			cl.RLock()
			length := len(cl.clients)
			if length == 0 {
				cl.l.Debug().Msg("No clients are listening at the moment")
				cl.RUnlock()
				continue
			}
			if debug {
				cl.l.Debug().Int("count", length).Msg("Pinging clients")
			}
			wg := sync.WaitGroup{}
			wg.Add(length)
			for _, client := range cl.clients {
				go func(c *Client) {

					c.write(websocket.PingMessage, nil)
					wg.Done()
				}(client)
			}
			cl.RUnlock()
			wg.Wait()

		}
	}
}

type Client struct {
	id         string
	conn       *websocket.Conn
	options    *WsOptions
	l          logger.AppLogger
	clientLIst *ClientList
}

func (c *Client) write(kind int, j []byte) error {
	if err := c.conn.SetWriteDeadline(time.Now().Add(c.options.WriteTimeout)); err != nil {
		c.l.Error().Err(err).Msg("Failed to set write deadline")
		c.clientLIst.Unsubscribe(c.id)
		c.conn.Close()
		return err
	}
	if err := c.conn.WriteMessage(kind, j); err != nil {
		c.l.Error().Err(err).Msg("Failed to write message to the client")
		c.clientLIst.Unsubscribe(c.id)
		c.conn.Close()
		return err
	}
	return nil
}
