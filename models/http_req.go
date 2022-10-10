package models

import (
	"crypto/tls"
	"io"
	"net/http"

	oo "github.com/Port3-Network/liboo"
	"github.com/pkg/errors"
)

func HttpGet(url string, timeout int64) (data []byte, err error) {
	if len(url) <= 0 {
		return nil, errors.Wrap(err, "empty uri error")
	}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	transport := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
	}

	client := &http.Client{
		Transport: transport,
		// Timeout:   time.Duration(timeout) * time.Second,
	}

	resp, err := client.Do(request)
	if err != nil {
		oo.LogD("url:[%s] client.Do err, msg: %v", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	if http.StatusOK != resp.StatusCode {
		return nil, errors.Wrap(err, "StatusCode error")
	}

	buf, err := io.ReadAll(resp.Body)
	if nil != err {
		oo.LogD("url:[%s] ioutil.ReadAll err, msg: %v", url, err)
		return nil, errors.Wrap(err, "ioutil.ReadAll error")
	}
	return buf, nil
}
