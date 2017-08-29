package jsonrpc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

// Process incoming request with dispatcher
func Process(reader io.Reader, dispatcher Dispatcher, writer io.Writer) {
	defer func() {
		if err := recover(); err != nil {
			writer.Write([]byte(fmt.Sprintf(InternalErrorTemplate, err)))
		}
	}()

	process(reader, dispatcher, writer)
}

func process(reader io.Reader, dispatcher Dispatcher, writer io.Writer) {
	buffer := bufio.NewReader(reader)
	decoder := json.NewDecoder(buffer)
	encoder := json.NewEncoder(writer)

	token, _, err := buffer.ReadRune()
	if err != nil {
		encoder.Encode(NewErrorResponse(NewError(ParseError, err)))
	}

	switch token {
	case '[':
		buffer.UnreadRune()

		ProcessBatchMessage(decoder, dispatcher, encoder)
	case '{':
		buffer.UnreadRune()

		ProcessSingleMessage(decoder, dispatcher, encoder)
	default:
		encoder.Encode(NewErrorResponse(NewError(ParseError, "unexpected token ", token)))
	}
}

// ProcessBatchMessage reads array of Message from JSON decoder and process them by provided dispatcher in parallel
func ProcessBatchMessage(decoder *json.Decoder, dispatcher Dispatcher, encoder *json.Encoder) {
	var messages []Message

	if err := decoder.Decode(&messages); err != nil {
		encoder.Encode(NewErrorResponse(NewError(ParseError, err)))
	}

	if len(messages) == 0 {
		encoder.Encode(NewErrorResponse(NewError(InvalidRequest)))
	}

	proxy := NewBatchProxy(encoder)

	for _, message := range messages {
		if message.IsRequest() {
			request := message.ToRequest()
			if request.IsNotification() {
				dispatcher.DispatchNotification(request)
			} else {
				proxy.processes.Add(1)
				go request.Process(dispatcher, Responder{request, proxy})
			}
		} else {
			message.ToResponse().Process(dispatcher)
		}
	}

	proxy.Wait()
}

// ProcessSingleMessage reads one Message from JSON decoder and proccess it by provided dispatcher
func ProcessSingleMessage(decoder *json.Decoder, dispatcher Dispatcher, encoder *json.Encoder) {
	var message Message

	if err := decoder.Decode(&message); err != nil {
		encoder.Encode(NewErrorResponse(NewError(ParseError, err)))
	}

	if message.IsRequest() {
		request := message.ToRequest()
		if request.IsNotification() {
			dispatcher.DispatchNotification(request)
		} else {
			request.Process(dispatcher, Responder{request, SingleProxy{encoder}})
		}
	} else {
		message.ToResponse().Process(dispatcher)
	}
}
