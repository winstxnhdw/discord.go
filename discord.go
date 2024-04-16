package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Discord struct {
	httpClient            *http.Client
	token                 string
	rootEndpoint          string
	rateLimit             time.Duration
	context               context.Context
	contextCancelFunction context.CancelFunc
}

type DiscordMessage struct {
	Content string `json:"content"`
}

func Create(discordToken string, rateLimit time.Duration) *Discord {
	context, contextCancelFunction := context.WithCancel(context.Background())

	return &Discord{
		httpClient:            &http.Client{},
		token:                 discordToken,
		rootEndpoint:          "https://discord.com/api/v9",
		rateLimit:             rateLimit,
		context:               context,
		contextCancelFunction: contextCancelFunction,
	}
}

func (d Discord) request(context context.Context, endpoint string, data io.Reader) ([]byte, error) {
	select {
	case <-time.After(d.rateLimit):
	case <-context.Done():
		return nil, d.context.Err()
	}

	request, err := http.NewRequest("POST", d.rootEndpoint+endpoint, data)

	if err != nil {
		return nil, err
	}

	request.Header = map[string][]string{
		"User-Agent":    {"Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/116.0"},
		"Content-Type":  {"application/json"},
		"Authorization": {d.token},
	}

	response, err := d.httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := response.Body.Close()

		if err != nil {
			panic(err)
		}
	}()

	return io.ReadAll(response.Body)
}

func (d Discord) Message(id string, message string) ([]byte, context.CancelFunc, error) {
	requestContext, requestContextCancelFunction := context.WithCancel(d.context)

	data, err := json.Marshal(&DiscordMessage{Content: message})

	if err != nil {
		return nil, requestContextCancelFunction, err
	}

	body, err := d.request(requestContext, "/channels/"+id+"/messages", bytes.NewBuffer(data))

	return body, requestContextCancelFunction, err
}

func (d Discord) Dispose() {
	d.contextCancelFunction()
	d.httpClient.CloseIdleConnections()
}
