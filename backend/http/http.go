package http

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type HttpClient struct {
	client     *http.Client
	cookie     *cookiejar.Jar
	url        *url.URL
	name       string
	cookiesMap map[string]string
}

func NewHttp(host, name string) *HttpClient {
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(host)

	return &HttpClient{
		client:     &http.Client{},
		cookie:     jar,
		url:        u,
		name:       name,
		cookiesMap: make(map[string]string),
	}
}

func (t *HttpClient) SendRequestByReq(req *http.Request, from string) ([]byte, error) {
	res, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}
	g.Log("http").Info(context.Background(), "SendRequestByReq response status", res.StatusCode)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	for _, c := range res.Cookies() {
		t.AddCookies(c.Name, c.Value)
	}
	t.SaveCookieToFile(from)
	return body, nil
}

func (t *HttpClient) GetNewRequest(method string, httpUrl string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, httpUrl, body)
	if err != nil {
		return req, err
	}
	return req, err
}

func (t *HttpClient) SendRequest(method string, httpUrl string, body url.Values, form string) ([]byte, error) {
	req, err := http.NewRequest(method, httpUrl, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}
	return t.SendRequestByReq(req, form)
}

func (t *HttpClient) GetCookies() []*http.Cookie {
	return t.cookie.Cookies(t.url)
}

func (t *HttpClient) AddCookies(key, value string) {
	t.cookie.SetCookies(t.url, []*http.Cookie{{Name: key, Value: value}})
}

func (t *HttpClient) GetCookiesStr() string {
	var cookiesStr1 = ""
	for _, cookies := range t.GetCookies() {
		cookiesStr1 += cookies.Name + "=" + cookies.Value + "; "
	}
	return cookiesStr1
}

func (t *HttpClient) LoadCookieFromFile(form string) (map[string]interface{}, error) {
	var cookiesMap = make(map[string]interface{})
	var value string

	var cookeFileName = form + "_cookies" + ".txt"
	if gfile.Exists(cookeFileName) {
		err := gfile.ReadLines(cookeFileName, func(line string) error {
			value = value + line
			return nil
		})
		if err != nil {
			g.Log("http").Errorf(context.Background(), "LoadCookieFromFile from: %v err:%v", form+"_cookies"+".txt", err.Error())
			return cookiesMap, err
		}
	} else {
		return nil, errors.New("no found cookies")
	}
	if len(value) > 0 {
		cookiesList := strings.Split(value, "; ")
		for _, cookie := range cookiesList {
			if strings.Contains(cookie, "=") {
				kv := strings.Split(cookie, "=")
				cookiesMap[kv[0]] = kv[1]
				t.AddCookies(kv[0], kv[1])
			}
		}
		return cookiesMap, nil
	}
	return nil, errors.New("no found cookies")
}

func (t *HttpClient) SaveCookieToFile(form string) {
	file, err := gfile.Create(form + "_cookies" + ".txt")
	if err != nil {
		g.Log("http").Error(context.Background(), "Create file tool_cookies.json err", err)
		return
	}
	defer file.Close()
	_, _ = file.WriteString(t.GetCookiesStr())
	g.Log("http").Info(context.Background(), "SaveCookieToFile to ", file.Name())
}
