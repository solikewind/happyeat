package menu

import "testing"

func TestSameObjectImageURL(t *testing.T) {
	raw := "https://happyeat.example.com/menu/20260508/1778251268744.jpg"
	signed := raw + "?q-sign-algorithm=sha1&q-ak=AKID"

	if !sameObjectImageURL(signed, raw) {
		t.Fatal("signed url should match raw url")
	}
	if !sameObjectImageURL(raw, raw) {
		t.Fatal("identical urls should match")
	}
	if sameObjectImageURL(raw, "https://other.example.com/a.jpg") {
		t.Fatal("different urls should not match")
	}
}
