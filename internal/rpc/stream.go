package rpc

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type (
	//Stream which used to read/write messages
	Stream interface {
		Read() (Response, error)
		Write(Request) error
		Close() error
	}

	//Codec specific how to encode/decode stream data
	Codec interface {
		Decode([]byte) (Response, error)
		Encode(Request) ([]byte, error)
	}

	websocketStream struct {
		conn  *websocket.Conn
		codec Codec
	}
)

//NewWebsocketStream create a new websocket stream with specific codec
func NewWebsocketStream(addr string, codec Codec) (Stream, error) {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, errors.WithMessagef(err, "websocket conn create fail")
	}

	return &websocketStream{
		conn:  conn,
		codec: codec,
	}, nil
}

func (ws *websocketStream) Read() (Response, error) {
	_, msg, err := ws.conn.ReadMessage()
	if err != nil {
		return nil, NewStreamError(err)
	}

	ret, err := ws.codec.Decode(msg)
	if err != nil {
		if _, ok := err.(*MsgError); ok {
			return nil, err
		}
		return nil, NewMsgError(msg, err)
	}
	return ret, nil
}

func (ws *websocketStream) Write(req Request) error {
	msg, err := ws.codec.Encode(req)
	if err != nil {
		if _, ok := err.(*MsgError); ok {
			return err
		}
		return NewMsgError(msg, err)
	}
	if err := ws.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		return NewStreamError(err)
	}
	return nil
}

func (ws *websocketStream) Close() error {
	return ws.conn.Close()
}
