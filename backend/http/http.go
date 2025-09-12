package http

import (
	"errors"
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

func (t *HttpClient) GetCookiesStr(from string) string {
	var cookiesStr1 = ""
	for _, cookies := range t.GetCookies() {
		if len(cookies.Value) > 0 && len(cookies.Name) > 0 {
			if strings.EqualFold(from, "xhs") && (strings.EqualFold(cookies.Name, "web_session") || strings.EqualFold(cookies.Name, "acw_tc")) {
				break
			}
			if strings.EqualFold(from, "dy") && !strings.EqualFold(cookies.Name, "ttwid") {
				break
			}
			cookiesStr1 += cookies.Name + "=" + cookies.Value + "; "
		}
	}
	return cookiesStr1
}

func (t *HttpClient) LoadCookieFromFile() (map[string]interface{}, error) {
	var cookiesMap = make(map[string]interface{})
	value := t.cookiesMap[t.name]
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

func (t *HttpClient) SaveCookieToFile(from string) {
	t.cookiesMap[t.name] = t.GetCookiesStr(from)
}
