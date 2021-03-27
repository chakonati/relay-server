package handlers

import (
	"context"
	"io"
	"log"
	"net/http"
	"server/decoders"
	"server/encoders"
	"server/session"
	"sync"
	"time"

	"github.com/pkg/errors"

	"gitlab.com/xdevs23/go-reflectutil/reflectutil"

	"nhooyr.io/websocket"
)

const bufSize = 512

type MessageType int

const (
	MessageTypeRequest MessageType = iota
	MessageTypeResponse
	MessageTypeOneway
	MessageTypeStream
)

type Request struct {
	Id     int64         `key:"id"`
	Action string        `key:"action"`
	Data   []interface{} `key:"data"`
}

type Message struct {
	Id          int64         `key:"id"`
	MessageType MessageType   `key:"messageType"`
	Data        []interface{} `key:"data"`
}

type Response struct {
	Message
}

type Handler struct {
	Session        session.Session
	context        context.Context
	decoder        decoders.Decoder
	encoder        encoders.Encoder
	messageType    websocket.MessageType
	requestChannel chan *Request
}

var runningRequests sync.WaitGroup

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
			Message{
				Id:          request.Id,
				MessageType: MessageTypeResponse,
				Data:        []interface{}{"hi!"},
			},
		})
		if err != nil {
			log.Println(err)
		}
	}

	for {
		var msgType websocket.MessageType
		msgType, byt, err = h.Session.Conn.Read(h.context)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("EOF")
				break
			}
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				break
			}
			log.Println("websocket connection read error:", err)
			break
		}
		if len(byt) == 0 {
			continue
		}
		if msgType != h.messageType {
			log.Println("message type mismatch")
			continue
		}

		request := Request{}
		if err = h.decoder.Unmarshal(byt, &request); err != nil {
			log.Println("decoder unmarshal error:", err)
			break
		}

		runningRequests.Add(1)
		h.requestChannel <- &request
	}

	close(requestChannel)

	runningRequests.Wait()
	ctx, cancelCtx := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	<-h.Session.Conn.CloseRead(ctx).Done()
	cancelCtx()
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
			defer runningRequests.Done()
			result, err := reflectutil.CallPathAutoConvert(actionHandler, request.Action, request.Data...)
			if err != nil {
				log.Println(err)
				return
			}
			lastResult := result[len(result)-1]
			if resultErr, isErr := lastResult.(error); isErr {
				log.Printf("Error during call to %s: %v", request.Action, resultErr)
			}

			encoders.ConvertErrorsToString(&result)
			err = h.reply(&Response{
				Message{
					Id:          request.Id,
					MessageType: MessageTypeResponse,
					Data:        result,
				},
			})
			if err != nil {
				log.Println(err)
			}
		}()
	}
}
