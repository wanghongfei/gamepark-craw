package crawl

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var HttpClient *http.Client

func init() {
	c := &http.Client{
		Timeout:       time.Second * 20,
	}

	HttpClient = c
}

func GetWithRetry(url string, maxRetry int) (io.Reader, error) {
	if maxRetry < 0 {
		maxRetry = 0
	}

	for tot := 0; tot <= maxRetry; tot++ {
		res, err := Get(url)
		if nil != err {
			log.Printf("failed to request %s, %v\n", url, err)
			continue
		}

		return res, nil
	}

	return nil, fmt.Errorf("failed to request %s after %d retries", url, maxRetry)
}

func Get(url string) (io.Reader, error) {
	log.Printf("send request to %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	addHeader(req)

	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	buf, _ := ioutil.ReadAll(resp.Body)
	return bytes.NewReader(buf), nil
}

func addHeader(req *http.Request) {
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	//req.Header.Add("Accept-Language", "en-US;q=0.8,en;q=0.7")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")
}
