package utils

import (
	"testing"
)

func TestZeroFill(t *testing.T) {
	t.Log(ZeroFill(1, 100/10))
}

func TestHidePhone(t *testing.T) {
	t.Log(HidePhone("15055461510"))
}
