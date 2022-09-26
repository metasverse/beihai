package jwt

import "testing"

func TestParseToken(t *testing.T) {
	id := 1
	token, err := GenToken(int64(id))
	if err != nil {
		t.Error(err)
	}
	t.Log(token)
	//id2, err := ParseToken(token)
	id2, err := ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiY3JlYXRlX2F0IjoxNjYwNDA4NDU0fQ.BuUJ8W9HVWhnWzLAv_lEGrVMF-4X22JRRsTADmdL2GE")
	if err != nil {
		t.Error(err)
	}
	t.Log(id2)
}
