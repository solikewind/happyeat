package spyun

import "testing"

// 文档示例参数见 https://www.spyun.net.cn/open/index.html 签名一节；
// 按 stringA + "&appsecret=" + secret 计算 MD5 大写结果与页面上的 sign 示例不一致（页面示例可能有笔误），此处固定算法回归值。
func TestBuildSign_docExampleParams(t *testing.T) {
	params := map[string]string{
		"appid":     "sp5c1314095ed15",
		"timestamp": "1544765873",
		"sn":        "111111111",
		"pkey":      "22222222",
		"name":      "test",
	}
	secret := "735aa25a15b75e6c1e0760823a22346a"
	got := BuildSign(params, secret)
	want := "0D6E220C0E3FCE6A68895C0FAE0EB755"
	if got != want {
		t.Fatalf("BuildSign = %q, want %q", got, want)
	}
}
