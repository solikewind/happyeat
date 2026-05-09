package spyun

import (
	"crypto/md5"
	"encoding/hex"
	"sort"
	"strings"
)

// BuildSign 按商鹏开放平台规则生成 sign：参数名 ASCII 排序，空值不参与，sign 本身不参与；
// stringSignTemp = stringA + "&appsecret=" + appsecret，再 MD5 转大写十六进制。
// 参见 https://www.spyun.net.cn/open/index.html
func BuildSign(params map[string]string, appSecret string) string {
	keys := make([]string, 0, len(params))
	for k, v := range params {
		if k == "sign" || v == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+params[k])
	}
	stringA := strings.Join(parts, "&")
	stringSignTemp := stringA + "&appsecret=" + appSecret
	sum := md5.Sum([]byte(stringSignTemp))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}
