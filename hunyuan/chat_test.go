package hunyuan

import (
	"context"
	"testing"
)

func TestClient_CreateChatCompletion(t *testing.T) {
	var a float32 = 1.0
	var p float32 = 0.8
	var s int = 0
	ak := ""
	sk := ""
	appid := 0
	client, _ := NewClient(ak, sk, appid)

	req := ChatCompletionRequest{
		Messages: []ChatCompletionMessage{
			{
				Role:    "user",
				Content: "帮我写一首诗歌",
			},
		},
		Stream:      &s,
		Temperature: &a,
		TopP:        &p,
	}
	t.Log(client.CreateChatCompletion(context.Background(), req))
	//r, err := client.CreateChatCompletionStream(context.Background(), req)
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Log(r)
	//for {
	//	fmt.Println(1)
	//	aa, err := r.Recv()
	//	if err != nil {
	//		t.Error(err)
	//		if errors.Is(err, io.EOF) {
	//			fmt.Println(1)
	//		}
	//		break
	//	}
	//	t.Log(aa)
	//}
}
