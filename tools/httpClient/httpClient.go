package httpClient

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

type HttpClient struct {
	Url          string
	RequestBody  []byte
	ResponseBody []byte

	Header    http.Header
	Timeout   int64 //default time.Second * 60
	Transport *http.Transport
}

func (that *HttpClient) POST() (err error) {
	client := &http.Client{
		Timeout: time.Second * 60,
	}
	req, err := http.NewRequest("POST", that.Url, bytes.NewBuffer(that.RequestBody))
	if err != nil {
		return
	}
	if len(that.Header) > 0 {
		req.Header = that.Header
	}
	if that.Timeout > 0 {
		client.Timeout = time.Duration(that.Timeout)
	} else {
		req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1")
	}
	if that.Transport != nil {
		client.Transport = that.Transport
	} else {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	that.ResponseBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}
