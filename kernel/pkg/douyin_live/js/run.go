package js

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/dop251/goja"
	"github.com/elliotchance/orderedmap"
	"os"
	"strings"
)

var (
	runtime *goja.Runtime
	getSign func(string) string
)

// js代码来源:https://github.com/hua0512/stream-rec/blob/f5a13e5ccc7df051b4f537321ffa259275aaa1ba/platforms/src/main/resources/douyin-webmssdk.js
// https://github.com/LyzenX/DouyinLiveRecorder/blob/main/dylr/util/webmssdk.js
// 感谢前辈们做出的贡献

func Init(filename string, ua string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	jsStr := string(file)
	runtime = goja.New()
	err = runtime.Set("navigator.userAgent", ua)
	if err != nil {
		return err
	}
	jsdom := "navigator = {" +
		"userAgent: '" + ua + "'" +
		"};" +
		"window=this;" +
		"document ={};" +
		"window.navigator = navigator;" +
		"setTimeout=function(){};"

	_, err = runtime.RunString(jsdom + jsStr)
	if err != nil {
		return err
	}

	err = runtime.ExportTo(runtime.Get("get_sign"), &getSign)
	return err
}

func Signature(params *orderedmap.OrderedMap) string {
	// 使用 strings.Builder 构建签名字符串
	var sigParams strings.Builder
	first := true
	for _, key := range params.Keys() {
		if !first {
			sigParams.WriteString(",")
		} else {
			first = false
		}
		value, _ := params.Get(key)
		sigParams.WriteString(fmt.Sprintf("%s=%s", key, value))
	}
	hash := md5.Sum([]byte(sigParams.String()))
	return getSign(hex.EncodeToString(hash[:]))
}
