package hunyuan

import (
	"bufio"
	"context"
	utils "github.com/Lok-Lu/go-tllm/internal"
)

type ChatCompletionStream struct {
	*streamReader[*ChatCompletionResponse]
}

// CreateChatCompletionStream — API call to create a chat completion w/ streaming
// support. It sets whether to stream back partial progress. If set, tokens will be
// sent as data-only server-sent events as they become available, with the
// stream terminated by a data: [DONE] message.
func (c *Client) CreateChatCompletionStream(
	ctx context.Context,
	request ChatCompletionRequest,
) (stream *ChatCompletionStream, err error) {
	c.SetRequestDefault(&request)
	urlSuffix := "/hyllm/v1/chat/completions"

	req, err := c.newStreamRequest(ctx, "POST", urlSuffix, c.GenerateAuthToken(request, urlSuffix), request)
	if err != nil {
		return
	}

	resp, err := c.config.HTTPClient.Do(req) //nolint:bodyclose // body is closed in stream.Close()
	if err != nil {
		return
	}
	//if isFailureStatusCode(resp) {
	//	return nil, c.handleErrorResp(resp)
	//}

	stream = &ChatCompletionStream{
		streamReader: &streamReader[*ChatCompletionResponse]{
			emptyMessagesLimit: c.config.EmptyMessagesLimit,
			reader:             bufio.NewReader(resp.Body),
			response:           resp,
			errAccumulator:     utils.NewErrorAccumulator(),
			unmarshaler:        &utils.JSONUnmarshaler{},
		},
	}
	return
}
