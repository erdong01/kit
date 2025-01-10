package httpClient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type HttpClient struct {
	Url          string
	RequestBody  []byte
	ResponseBody []byte
	Header       http.Header
	Timeout      int64 //default time.Second * 60
	Transport    *http.Transport
	Status       string // e.g. "200 OK"
	StatusCode   int    // e.g. 200
	Err          error
	Method       string
}

func New(url string) *HttpClient {
	return &HttpClient{
		Url: url,
	}
}

func (that *HttpClient) POST(data ...[]byte) *HttpClient {
	return that.SetMethod("POST").Do(data...)
}

func (that *HttpClient) HeaderAdd(key string, value string) *HttpClient {
	if that.Header == nil {
		that.Header = make(http.Header)
	}
	that.Header.Add(key, value)
	return that
}

func (that *HttpClient) SetTransport(transport *http.Transport) *HttpClient {
	that.Transport = transport
	return that
}

func (that *HttpClient) SetTimeout(timeout int64) *HttpClient {
	that.Timeout = timeout
	return that
}
func (that *HttpClient) SetMethod(method string) *HttpClient {
	that.Method = method
	return that
}

func (that *HttpClient) Get(data ...[]byte) *HttpClient {
	return that.SetMethod("GET").Do(data...)
}

func (that *HttpClient) Scan(v any) (err error) {
	return json.Unmarshal(that.ResponseBody, v)
}

func (that *HttpClient) Do(data ...[]byte) *HttpClient {
	client := &http.Client{
		Timeout:   time.Second * 60,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	if len(data) > 0 {
		that.RequestBody = data[0]
	}

	req, err := http.NewRequest(that.Method, that.Url, bytes.NewBuffer(that.RequestBody))
	if err != nil {
		that.Err = err
		return that
	}
	if len(that.Header) > 0 {
		req.Header = that.Header
	} else {
		req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1")
		req.Header.Set("Content-Type", "application/json")
	}
	if that.Timeout > 0 {
		client.Timeout = time.Duration(that.Timeout)
	}
	if that.Transport != nil {
		client.Transport = that.Transport
	}
	resp, err := client.Do(req)
	if err != nil {
		that.Err = err
		return that
	}
	defer resp.Body.Close()
	that.ResponseBody, err = io.ReadAll(resp.Body)
	that.StatusCode = resp.StatusCode
	that.Status = resp.Status
	that.Err = err
	return that
}
