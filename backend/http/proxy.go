/*
 * @Author: Morningsun 1049935747@qq.com
 * @Date: 2025-09-08 22:28:23
 * @LastEditors: Morningsun 1049935747@qq.com
 * @LastEditTime: 2025-09-08 23:17:23
 * @Description:
 */
package http

import (
	"io"
	"net/http"
	stdurl "net/url"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Constants for server configuration
const (
	// UpStream defines the backend server URL that the proxy will forward requests to
	UpStream = "http://127.0.0.1:17892"
)

// StartProxy 启动一个 goframe 反向代理服务
// 访问示例：/proxy/{name}?u={encodedTarget}
// name: xhs | dy | default，用于设置不同站点需要的请求头
func StartProxy() {
	s := g.Server("mediaProxy")

	s.Group("").GET("/proxy/{name}", func(r *ghttp.Request) {
		name := r.GetRouter("name").String()
		raw := r.GetQuery("u").String()
		if raw == "" {
			r.Response.WriteStatus(http.StatusBadRequest, "missing u")
			return
		}
		// 校验 URL
		if _, err := stdurl.Parse(raw); err != nil {
			r.Response.WriteStatus(http.StatusBadRequest, "bad url")
			return
		}

		// 组装请求
		req, err := http.NewRequest("GET", raw, nil)
		if err != nil {
			r.Response.WriteStatus(http.StatusInternalServerError, err.Error())
			return
		}

		// 根据名称设置请求头
		switch name {
		case "xhs":
			req.Header.Set("Referer", "https://www.xiaohongshu.com/")
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36 Edg/134.0.0.0")
			req.Header.Set("Accept", "*/*")
			req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
			req.Header.Set("Sec-Fetch-Dest", "document")
		case "dy":
			req.Header.Set("Referer", "https://www.douyin.com/")
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36 Edg/134.0.0.0")
			req.Header.Set("Accept", "*/*")
			req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
			req.Header.Set("Sec-Fetch-Dest", "document")
		default:
			req.Header.Set("Accept", "*/*")
			req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
			req.Header.Set("Sec-Fetch-Dest", "document")
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			r.Response.WriteStatus(http.StatusBadGateway, err.Error())
			return
		}
		defer resp.Body.Close()

		// 透传响应头
		if ct := resp.Header.Get("Content-Type"); ct != "" {
			r.Response.Writer.Header().Set("Content-Type", ct)
		}
		if cl := resp.Header.Get("Content-Length"); cl != "" {
			r.Response.Writer.Header().Set("Content-Length", cl)
		}
		r.Response.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(r.Response.Writer, resp.Body)
	})

	// 监听地址
	s.SetAddr("127.0.0.1:17892")
	go s.Start()
}
