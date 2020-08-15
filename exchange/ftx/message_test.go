package ftx

import (
	"testing"

	"github.com/NadiaSama/ccexgo/internal/rpc"
)

func TestMessageDecode(t *testing.T) {
	cc := NewCodeC()

	e := []byte(`{"channel": "", "market": "", "type": "error", "code": 1001, "msg": "not login"}`)
	if resp, err := cc.Decode(e); err != nil {
		t.Errorf("parse error fail %s", err.Error())
	} else {
		if r := resp.(*rpc.Result); r.Error == nil {
			t.Errorf("expect error fail %v", *r)
		}
	}

	s := []byte(`{"channel": "ch", "market": "m", "type": "subscribed"}`)
	if resp, err := cc.Decode(s); err != nil {
		t.Errorf("parse error fail %s", err.Error())
	} else {
		if r := resp.(*rpc.Result); r.ID != "chm" {
			t.Errorf("bad id %v", *r)
		}
	}

	i := []byte(`{"type": "info", "code": 20001}`)
	if _, err := cc.Decode(i); err == nil {
		t.Errorf("parse info fail")
	} else {
		if _, ok := err.(*rpc.StreamError); !ok {
			t.Errorf("expect streamerror %v", err)
		}
	}
}
