package httpClient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"
)

type HttpClient struct {
	Url                string
	RequestBody        []byte
	ResponseBody       []byte
	Header             http.Header
	Timeout            int64 //default time.Second * 60
	Transport          *http.Transport
	Status             string // e.g. "200 OK"
	StatusCode         int    // e.g. 200
	Err                error
	Method             string
	InsecureSkipVerify bool
}

func New(url string) *HttpClient {
	return &HttpClient{
		Url:                url,
		InsecureSkipVerify: true,
	}
}

func (that *HttpClient) POST(data ...[]byte) *HttpClient {
	return that.SetMethod("POST").Do(data...)
}

// Add adds the key, value pair to the header.
// It appends to any existing values associated with key.
// The key is case insensitive; it is canonicalized by
// [CanonicalHeaderKey].
func (that *HttpClient) HeaderAdd(key string, value string) *HttpClient {
	if that.Header == nil {
		that.Header = make(http.Header)
	}
	that.Header.Add(key, value)
	return that
}

// Get gets the first value associated with the given key. If
// there are no values associated with the key, Get returns "".
// It is case insensitive; [textproto.CanonicalMIMEHeaderKey] is
// used to canonicalize the provided key. Get assumes that all
// keys are stored in canonical form. To use non-canonical keys,
// access the map directly.
func (that *HttpClient) HeaderSet(key string, value string) *HttpClient {
	if that.Header == nil {
		that.Header = make(http.Header)
	}
	that.Header.Set(key, value)
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
	if that.Transport == nil {
		that.Transport = &http.Transport{
			DialContext:           (&net.Dialer{Timeout: 30 * time.Second, KeepAlive: 60 * time.Second}).DialContext,
			TLSHandshakeTimeout:   30 * time.Second, // 调整TLS握手超时时间
			ResponseHeaderTimeout: 30 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConns:          64,               // 最大空闲连接数
			MaxIdleConnsPerHost:   32,               // 每个主机的最大空闲连接数
			IdleConnTimeout:       90 * time.Second, // 空闲连接的超时时间
		}
	}

	if that.InsecureSkipVerify {
		that.Transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client := &http.Client{
		Timeout:   time.Second * 60,
		Transport: that.Transport,
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
		if req.Header.Get("User-Agent") == "" {
			req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1")
		}
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}
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
