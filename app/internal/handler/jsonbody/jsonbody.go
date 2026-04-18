// Package jsonbody 提供不依赖 Content-Length 的 JSON 请求体解析。
//
// go-zero 的 httpx.Parse 在 withJsonBody 为 false 时（常见于 chunked、Content-Length 未设置）
// 不会读取 JSON body，可选数字字段（如 sort）会恒为零。凡依赖 JSON body 的 POST/PUT 可改用本包 Decode。
package jsonbody

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const maxBody = 1 << 20 // 1MB，与 go-zero httpx 默认上限同量级

// Decode 读取并 json.Unmarshal 到 v；空 body 不报错。调用方需先 ParsePath 再 Decode 以合并 path 与 body。
func Decode(r *http.Request, v any) error {
	if r.Body == nil {
		return nil
	}
	defer r.Body.Close()
	b, err := io.ReadAll(io.LimitReader(r.Body, maxBody))
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return nil
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	return dec.Decode(v)
}
