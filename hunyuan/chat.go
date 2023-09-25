package hunyuan

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/lithammer/shortuuid/v3"
	"net/http"
	"strings"
	"time"
)

// Chat message role defined by the Sensa API.

type ModelName string

var (
	ErrChatCompletionStreamNotSupported = errors.New("streaming is not supported with this method, please use CreateChatCompletionStream") //nolint:lll
)

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest represents a request structure for chat completion API.
type ChatCompletionRequest struct {
	AppID       int                     `json:"app_id"`
	SecretID    string                  `json:"secret_id"`
	Timestamp   int                     `json:"timestamp"`
	Expired     int                     `json:"expired"`
	QueryID     string                  `json:"query_id"`
	Temperature *float32                `json:"temperature,omitempty"`
	TopP        *float32                `json:"top_p,omitempty"`
	Stream      *int                    `json:"stream,omitempty"`
	Messages    []ChatCompletionMessage `json:"messages"`
}

type Delta struct {
	Content string `json:"content"`
}

type ChatCompletionChoice struct {
	Messages     ChatCompletionMessage `json:"messages,omitempty"`
	FinishReason string                `json:"finish_reason"`
	Delta        Delta                 `json:"delta,omitempty"`
}

type ChatCompletionResponseErr struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// ChatCompletionResponse represents a response structure for chat completion API.
type ChatCompletionResponse struct {
	ID      string                    `json:"id"`
	Created string                    `json:"created"`
	Choices []ChatCompletionChoice    `json:"choices"`
	Usage   Usage                     `json:"usage"`
	Error   ChatCompletionResponseErr `json:"error"`
	Note    string                    `json:"note"`
	ReqID   string                    `json:"req_id"`
}

// CreateChatCompletion — API call to Create a completion for the chat message.
func (c *Client) CreateChatCompletion(
	ctx context.Context,
	request ChatCompletionRequest,
) (response *ChatCompletionResponse, err error) {
	//if request.Stream == 1 {
	//	err = ErrChatCompletionStreamNotSupported
	//	return
	//}

	c.SetRequestDefault(&request)
	urlSuffix := "/hyllm/v1/chat/completions"
	req, err := c.requestBuilder.Build(ctx, http.MethodPost, c.fullURL(urlSuffix), request)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response, c.GenerateAuthToken(request, urlSuffix))
	return
}

// GenerateAuthToken  Generate Auth Token
func (c *Client) GenerateAuthToken(req ChatCompletionRequest, urlSuffix string) string {
	var signStr strings.Builder

	signStr.WriteString(c.config.BaseURL)
	signStr.WriteString(urlSuffix)

	var msgStr strings.Builder
	for _, m := range req.Messages {
		msgStr.WriteString(fmt.Sprintf(`{"role":"%s","content":"%s"},`, m.Role, m.Content))
	}
	msg := strings.TrimSuffix(msgStr.String(), ",")

	signStr.WriteString(fmt.Sprintf("?app_id=%d&expired=%d&messages=%s&query_id=%s&secret_id=%s&stream=%d&",
		req.AppID, req.Expired, fmt.Sprintf("[%s]", msg), req.QueryID, req.SecretID, *req.Stream))

	if req.Temperature != nil {
		signStr.WriteString(fmt.Sprintf("temperature=%g&", *req.Temperature))
	}

	signStr.WriteString(fmt.Sprintf("timestamp=%d&", req.Timestamp))

	if req.TopP != nil {
		signStr.WriteString(fmt.Sprintf("top_p=%g", *req.TopP))
	}

	h := hmac.New(sha1.New, []byte(c.config.secretKey))
	h.Write([]byte(strings.TrimSuffix(signStr.String(), "&")))
	encryptedStr := h.Sum([]byte(nil))
	var signature = base64.StdEncoding.EncodeToString(encryptedStr)
	return signature
}

func (c *Client) SetRequestDefault(req *ChatCompletionRequest) {
	req.AppID = c.config.appid
	req.SecretID = c.config.accessKey
	req.Timestamp = int(time.Now().Unix() + 1000)
	req.Expired = int(time.Now().Unix() + 24*60*60)
	req.QueryID = shortuuid.NewWithNamespace("")
}
