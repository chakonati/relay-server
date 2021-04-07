package handlers

import (
	"context"
	"io"
	"log"
	"net/http"
	"server/decoders"
	"server/encoders"
	"server/session"
	"server/subscriptions"
	"sync"
	"sync/atomic"
	"time"

	"github.com/stoewer/go-strcase"

	"github.com/toowoxx/go-structs"

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

type Message interface {
	ID() int64
	SetID(id int64)
}

type Request struct {
	Id          int64         `key:"id"`
	MessageType MessageType   `key:"messageType"`
	Action      string        `key:"action"`
	Parameters  []interface{} `key:"parameters"`
}

func (r Request) ID() int64 {
	return r.Id
}

func (r Request) SetID(id int64) {
	r.Id = id
}

type MessageHeader struct {
	Id          int64       `key:"id"`
	MessageType MessageType `key:"messageType"`
}

var (
	_ Message = (*Request)(nil)
	_ Message = (*Notification)(nil)
)

type Handler struct {
	Session        session.Session
	context        context.Context
	decoder        decoders.Decoder
	encoder        encoders.Encoder
	messageType    websocket.MessageType
	requestChannel chan *Request
}

type Notification struct {
	Id          int64       `key:"id"`
	MessageType MessageType `key:"messageType"`
	Data        interface{} `key:"data"`
}

func (n Notification) ID() int64 {
	return n.Id
}

func (n Notification) SetID(id int64) {
	n.Id = id
}

var runningRequests sync.WaitGroup

var nextMessageId int64 = 0

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
		err := h.replyStruct(&request, struct{ Reply string }{"hi!"})
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

func (h *Handler) send(i interface{}) error {
	bytes, err := h.encoder.Marshal(i)
	if err != nil {
		return err
	}

	err = h.Session.Conn.Write(h.context, h.messageType, bytes)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) sendMessage(message Message) error {
	message.SetID(atomic.AddInt64(&nextMessageId, 1))
	return h.send(message)
}

func (h *Handler) reply(response map[string]interface{}) error {
	return h.send(response)
}

func (h *Handler) replyStruct(request *Request, obj interface{}) error {
	return h.sendStruct(request.Id, MessageTypeResponse, obj)
}

func normalizeMap(responseMap *map[string]interface{}) {
	for name, value := range *responseMap {
		delete(*responseMap, name)
		newName := strcase.LowerCamelCase(name)
		(*responseMap)[newName] = value
		switch value.(type) {
		case map[string]interface{}:
			m := value.(map[string]interface{})
			normalizeMap(&m)
		case error:
			(*responseMap)[newName] = value.(error).Error()
		}
	}
}

func (h *Handler) sendStruct(requestId int64, messageType MessageType, obj interface{}) error {
	header := &MessageHeader{
		Id:          requestId,
		MessageType: messageType,
	}
	var responseMap map[string]interface{}
	if obj != nil {
		responsePayload := obj
		responseMap = structs.Map(responsePayload)
		structs.FillMap(header, responseMap)
	} else {
		responseMap = structs.Map(header)
	}

	normalizeMap(&responseMap)

	return h.reply(responseMap)
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
			result, err := reflectutil.CallPathAutoConvert(actionHandler, request.Action, request.Parameters...)
			if err != nil {
				log.Println(err)
				return
			}
			if len(result) > 1 {
				log.Println("More than one result:", len(result))
			}
			var obj interface{} = nil
			if len(result) == 1 {
				obj = result[0]
			}
			err = h.replyStruct(request, obj)

			if err != nil {
				log.Println(err)
			}
		}()
	}
}

func (h *Handler) Disconnected() {
	subscriptions.Subscriptions.MessagesSubscription.Unsubscribe(h)
}
