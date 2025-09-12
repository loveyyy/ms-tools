/*
 * @Author: Morningsun 1049935747@qq.com
 * @Date: 2025-09-08 22:28:23
 * @LastEditors: Morningsun 1049935747@qq.com
 * @LastEditTime: 2025-09-12 22:16:03
 * @Description:
 */
package http

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"
	"net/http/httputil"
	stdurl "net/url"
)

func MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func StartProxy() {
	s := g.Server("proxy")
	// Create a new reverse proxy instance
	proxy := httputil.NewSingleHostReverseProxy(nil)
	// Configure error handling for proxy failures
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
		writer.WriteHeader(http.StatusBadGateway)
	}
	s.Use(MiddlewareCORS)
	s.Group("/proxy", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareCORS)
		group.ALL("/{name}", func(r *ghttp.Request) {
			var (
				originalPath = r.Request.URL.Path
				proxyToPath  = r.Get("u").String()
				name         = r.GetRouter("name").String()
			)
			// Parse the upstream URL
			u, _ := stdurl.Parse(proxyToPath)
			proxy.Director = func(req *http.Request) {
				req.URL.RawQuery = u.RawQuery
				req.URL.Scheme = u.Scheme
				req.URL.Host = u.Host
				req.URL.Path = u.Path
				req.URL.RawPath = u.RawPath
				req.Host = u.Host
			}
			// Log the proxy operation
			g.Log().Infof(r.Context(), `proxy:"%s" -> backend:"%s"`, originalPath, proxyToPath)
			// 根据名称设置请求头
			switch name {
			case "xhs":
				r.Request.Header.Set("Referer", "https://www.xiaohongshu.com/")
			case "dy":
				r.Request.Header.Set("Referer", "https://www.douyin.com/")
			default:
				r.Request.Header.Set("Accept", "*/*")
			}
			// Ensure request body can be read multiple times if needed
			r.MakeBodyRepeatableRead(false)
			// Forward the request to the backend server
			proxy.ServeHTTP(r.Response.Writer, r.Request)
		})
	})
	s.SetPort(17892)
	s.Run()
}
