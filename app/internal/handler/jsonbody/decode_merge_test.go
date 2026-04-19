package jsonbody

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/solikewind/happyeat/app/internal/types"
)

// 验证：先设 path 字段再 Decode body 时，id 不被清零、sort 能写入（与 UpdateMenuCategory 行为一致）
func TestDecodeMergePathAndBody(t *testing.T) {
	var req types.UpdateMenuCategoryReq
	req.Id = 7
	body := `{"name":"热菜","description":"","sort":3}`
	dec := json.NewDecoder(strings.NewReader(body))
	dec.UseNumber()
	if err := dec.Decode(&req); err != nil {
		t.Fatal(err)
	}
	if req.Id != 7 {
		t.Fatalf("Id was reset: got %d want 7", req.Id)
	}
	if req.Name != "热菜" || req.Description != "" || req.Sort != 3 {
		t.Fatalf("body fields: %+v", req)
	}
}

func TestDecodeUint32FromDecoderUseNumber(t *testing.T) {
	var req types.UpdateMenuCategoryReq
	dec := json.NewDecoder(bytes.NewReader([]byte(`{"sort":42}`)))
	dec.UseNumber()
	if err := dec.Decode(&req); err != nil {
		t.Fatal(err)
	}
	if req.Sort != 42 {
		t.Fatalf("Sort=%d want 42", req.Sort)
	}
}
