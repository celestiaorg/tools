package tools

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type options struct {
	timeout     time.Duration
	retrials    uint
	httpHeaders []map[string]string
}

type ReqOption interface {
	apply(*options)
}

type timeoutOption time.Duration

func (t timeoutOption) apply(opts *options) {
	opts.timeout = time.Duration(t)
}
func HttpWithTimeout(t time.Duration) ReqOption {
	return timeoutOption(t)
}

type httpHeadersOption map[string]string

func (h httpHeadersOption) apply(opts *options) {
	opts.httpHeaders = append(opts.httpHeaders, h)
}
func HttpAddHeader(key, value string) ReqOption {
	return httpHeadersOption(httpHeadersOption{key: value})
}

type retrialsOption uint

func (r retrialsOption) apply(opts *options) {
	opts.retrials = uint(r)
}
func HttpWithRetry(r uint) ReqOption {
	return retrialsOption(r)
}

func HttpGetReqPersist(url string, opts ...ReqOption) ([]byte, error) {
	// set default options
	options := options{
		timeout:  60,
		retrials: 20,
	}

	for _, o := range opts {
		o.apply(&options)
	}

	spaceClient := http.Client{
		Timeout: time.Second * options.timeout,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for _, h := range options.httpHeaders {
		for k, v := range h {
			req.Header.Set(k, v)
		}
	}

	var reqErr error
	var statusCode int

	for retry := uint(0); retry < options.retrials; retry++ {

		res, reqErr := spaceClient.Do(req)
		if res != nil {
			statusCode = res.StatusCode
		}

		if reqErr != nil || statusCode != 200 {

			if e, ok := reqErr.(net.Error); !ok || !e.Timeout() {
				if reqErr == nil {
					resBody := []byte{}
					if res != nil && res.Body != nil {
						resBody, _ = ioutil.ReadAll(res.Body)
					}

					reqErr = NewErrorf(statusCode, "http error: %d %v \n%s", statusCode, e, resBody)
				}
				return nil, reqErr
			}

			fmt.Printf("Retry: %d\n", retry+1)
			// Let's wait for a while and try again
			time.Sleep(250 * time.Millisecond)
			continue

		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		body, reqErr := ioutil.ReadAll(res.Body)
		if reqErr != nil {
			return nil, reqErr
		}
		return body, nil
	}

	if reqErr == nil {
		reqErr = NewErrorf(statusCode, "http error")
	}

	return nil, reqErr
}

func HttpPostReqPersist(url string, payload io.Reader, opts ...ReqOption) ([]byte, error) {
	// set default options
	options := options{
		timeout:  60,
		retrials: 20,
	}

	for _, o := range opts {
		o.apply(&options)
	}

	spaceClient := http.Client{
		Timeout: time.Second * options.timeout,
	}

	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}

	for _, h := range options.httpHeaders {
		for k, v := range h {
			req.Header.Set(k, v)
		}
	}

	var reqErr error
	var statusCode int

	for retry := uint(0); retry < options.retrials; retry++ {

		res, reqErr := spaceClient.Do(req)
		if res != nil {
			statusCode = res.StatusCode
		}

		if reqErr != nil || statusCode != 200 {

			if e, ok := reqErr.(net.Error); !ok || !e.Timeout() {
				if reqErr == nil {

					resBody := []byte{}
					if res != nil && res.Body != nil {
						resBody, _ = ioutil.ReadAll(res.Body)
					}

					reqErr = NewErrorf(statusCode, "http error: %d %v \n%s", statusCode, e, resBody)
				}
				return nil, reqErr
			}

			fmt.Printf("Retry: %d ", retry+1)
			// Let's wait for a while and try again
			time.Sleep(250 * time.Millisecond)
			continue

		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		body, reqErr := ioutil.ReadAll(res.Body)
		if reqErr != nil {
			return nil, reqErr
		}
		return body, nil
	}

	if reqErr == nil {
		reqErr = NewErrorf(statusCode, "http error")
	}
	return nil, reqErr
}

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() error {
	return nil
}

func ReadAll(rc io.ReadCloser) ([]byte, error) {
	defer rc.Close()

	if cb, ok := rc.(*ClosingBuffer); ok {
		return cb.Bytes(), nil
	}

	return ioutil.ReadAll(rc)
}
