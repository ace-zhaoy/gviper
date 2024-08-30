package notifications

import (
	"bytes"
	"fmt"
	"github.com/ace-zhaoy/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestNewFeishuBotHook(t *testing.T) {
	webhookAddr := "https://open.feishu.cn/open-apis/bot/v2/hook/xxxxx"
	f := NewFeishuBotHook(webhookAddr)
	assert.Equal(t, f.webhookAddr, webhookAddr)
}

func TestFeishuBotHook_Notify(t *testing.T) {
	f := NewFeishuBotHook(os.Getenv("FEISHU_WEBHOOK_ADDR"))
	f.Notify("test", errors.New("test notify"))

	// feishu msg:
	// Config test reload failed: test unmarshal error
}

func TestFeishuBotHook_SetPayloadBuilder(t *testing.T) {
	f := NewFeishuBotHook(os.Getenv("FEISHU_WEBHOOK_ADDR"))
	f.SetPayloadBuilder(func(configName string, err error) io.Reader {
		return bytes.NewBuffer([]byte(fmt.Sprintf(
			`{
    "msg_type": "post",
    "content": {
        "post": {
            "zh_cn": {
                "title": "配置发布失败",
                "content": [
                    [{
                        "tag": "text",
                        "text": "%s"
                    }, {
                        "tag": "text",
                        "text": "%v"
                    }, {
                        "tag": "a",
                        "text": "查看详情",
                        "href": "http://www.example.com/"
                    }]
                ]
            }
        }
    }
}`,
			configName,
			err,
		)))
	})
	f.Notify("test", errors.New("test setPayloadBuilder"))

	// feishu msg:
	// 配置发布失败
	// testtest setPayloadBuilder查看详情
}
