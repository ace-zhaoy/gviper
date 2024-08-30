package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type FeishuBotHook struct {
	webhookAddr    string
	payloadBuilder func(configName string, err error) io.Reader
}

func NewFeishuBotHook(webhookAddr string) *FeishuBotHook {
	return &FeishuBotHook{
		webhookAddr: webhookAddr,
	}
}

func (f *FeishuBotHook) SetPayloadBuilder(payloadBuilder func(configName string, err error) io.Reader) {
	f.payloadBuilder = payloadBuilder
}

func (f *FeishuBotHook) buildPayload(configName string, err error) io.Reader {
	if f.payloadBuilder != nil {
		return f.payloadBuilder(configName, err)
	}

	return bytes.NewBuffer([]byte(fmt.Sprintf(
		`{"msg_type":"text","content":{"text":"%s"}}`,
		fmt.Sprintf("Config %s reload failed: %v", configName, err),
	)))
}

func (f *FeishuBotHook) Notify(configName string, err error) {
	resp, err1 := http.Post(f.webhookAddr, "application/json", f.buildPayload(configName, err))
	if err1 != nil {
		log.Printf("feishu bot notify failed: %v", err1)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("feishu bot notify failed: %v", resp.StatusCode)
	}
	type Result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	var result Result
	err1 = json.NewDecoder(resp.Body).Decode(&result)
	if err1 != nil {
		log.Printf("feishu bot notify failed: %v", err1)
		return
	}
	if result.Code != 0 {
		log.Printf("feishu bot notify failed: %v %v", result.Code, result.Msg)
		return
	}
}
