package client

import (
	"encoding/json"
	"testing"
)

func TestCall(t *testing.T) {
	serv := NewService("127.0.0.1", 9500, "handler.product", "11111", "2222", "test")
	res, err := serv.Call("Add", "kovey1", "chelsea1")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var result bool
	json.Unmarshal(res, &result)
	t.Logf("result: %t", result)
}
