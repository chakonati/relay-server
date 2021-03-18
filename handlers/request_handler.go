package handlers

import (
	"context"
	"log"
	"net/http"
	"server/decoders"
	"server/encoders"
	"server/session"

	"gitlab.com/xdevs23/go-reflectutil/reflectutil"

	"nhooyr.io/websocket"
)

const bufSize = 512

type Request struct {
	Id     int64         `key:"id"`
	Action string        `key:"action"`
	Data   []interface{} `key:"data"`
}

type Response struct {
	Id   int64         `key:"id"`
	Data []interface{} `key:"data"`
}

type Handler struct {
	Session        session.Session
	context        context.Context
	decoder        decoders.Decoder
	encoder        encoders.Encoder
	messageType    websocket.MessageType
	requestChannel chan *Request
}

func (h *Handler) Handle() {
	h.context = context.Background()
	requestChannel := make(chan *Request)
	go h.startReceiving(requestChannel)

	messageType, byt, err := h.Session.Conn.Read(h.context)
	if err != nil {
		log.Println(err)
		return
	}
	if messageType == websocket.MessageBinary {
		h.decoder = decoders.Msgpack{}
		h.encoder = encoders.Msgpack{}
	} else {
		h.decoder = decoders.JSON{}
		h.encoder = encoders.JSON{}
	}

	h.requestChannel = requestChannel
	h.messageType = messageType

	request := Request{}
	if err = h.decoder.Unmarshal(byt, &request); err != nil {
		log.Println(err)
		return
	}

	if request.Action == "hello" {
		log.Println("hello received")
		err := h.reply(&Response{
			Id:   request.Id,
			Data: []interface{}{"hi!"},
		})
		if err != nil {
			log.Println(err)
		}
	}

	for {
		var msgType websocket.MessageType
		msgType, byt, err = h.Session.Conn.Read(h.context)
		if err != nil {
			log.Println(err)
			break
		}
		if len(byt) == 0 {
			continue
		}
		if msgType != h.messageType {
			log.Println("message type mismatch")
			continue
		}
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}

		request := Request{}
		if err = h.decoder.Unmarshal(byt, &request); err != nil {
			log.Println(err)
			break
		}

		h.requestChannel <- &request
	}

	close(requestChannel)

	if err != nil {
		_ = h.Session.Conn.Close(http.StatusInternalServerError, "error")
	} else {
		_ = h.Session.Conn.Close(http.StatusOK, "OK")
	}

	log.Println("Closed connection")
}

func (h *Handler) reply(response *Response) error {
	bytes, err := h.encoder.Marshal(response)
	if err != nil {
		return err
	}

	err = h.Session.Conn.Write(h.context, h.messageType, bytes)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) startReceiving(requestChan chan *Request) {
	for {
		request, hasMore := <-requestChan
		if !hasMore {
			break
		}

		actionHandler := &ActionHandler{
			handler: h,
			request: request,
		}
		log.Println("handling action", request.Action)

		go func() {
			result, err := reflectutil.CallPath(actionHandler, request.Action, request.Data...)
			if err != nil {
				log.Println(err)
				return
			}

			err = h.reply(&Response{
				Id:   request.Id,
				Data: result,
			})
			if err != nil {
				log.Println(err)
			}
		}()
	}
}
