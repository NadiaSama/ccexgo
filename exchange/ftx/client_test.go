package ftx

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestClientRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Deactivate()

	data := `{"a": "b", "c": 1.234444}`

	httpmock.RegisterResponder(http.MethodHead, "https://ftx.com/api/test123?name=a&war=b",
		func(req *http.Request) (*http.Response, error) {
			d, err := ioutil.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("read body fail %s", err.Error())
			}
			var v map[string]interface{}
			if err := json.Unmarshal(d, &v); err != nil {
				t.Fatalf("unmarshal fail %s", err.Error())
			}

			if v["a"].(string) != "b" || v["c"].(float64) != 1.234444 {
				t.Fatalf("bad value %v", v)
			}

			return httpmock.NewBytesResponse(200, []byte(`{"success": false, "error": "test error"}`)), nil
		})

	ctx := context.Background()
	client := NewRestClient("", "")
	param := url.Values{}
	param.Add("name", "a")
	param.Add("war", "b")
	err := client.request(ctx, http.MethodHead, "/test123", param, bytes.NewBufferString(data), true, nil)
	if !strings.Contains(err.Error(), "test error") {
		t.Errorf("bad error %s", err.Error())
	}
}

func TestBalance(t *testing.T) {
	ctx := context.Background()
	client := NewRestClient("IdAd0lvi9sq5C3zcckdI46gD0BG4c0DFHo1qdPKa", "flT9mRCn1KZBuFlg5JU6SpQKbXVtN1yif1JwnDvv")
	var ret []Balance
	err := client.request(ctx, http.MethodGet, "/subaccounts", nil, nil, true, &ret)
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("ret:%+v", ret)
}
