/*
 * @Author: Morningsun 1049935747@qq.com
 * @Date: 2025-03-18 20:55:56
 * @LastEditors: Morningsun 1049935747@qq.com
 * @LastEditTime: 2025-09-08 22:35:25
 * @Description:
 */
package backend

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"ms-tools/backend/dy"
	"ms-tools/backend/http"
	"ms-tools/backend/utils"
	"ms-tools/backend/xhs"
)

// App struct
type App struct {
	ctx  context.Context
	xhs1 *xhs.Xhs
	dy1  *dy.Dy
}

// NewApp creates a new App application struct
func NewApp(xhs1 *xhs.Xhs, dy1 *dy.Dy) *App {
	return &App{
		xhs1: xhs1,
		dy1:  dy1,
	}
}

// startup is called at application startup
func (a *App) Startup(ctx context.Context) {
	g.Log().Infof(ctx, "ms-tools isRunning version is %v , main dir is %v", g.Cfg().MustGet(ctx, "appVersion", "unknow"), utils.GetRootPath())
	// Perform your setup here
	a.ctx = ctx
	// 启动 goframe 反向代理（xhs/dy 等媒体资源）
	http.StartProxy()
	go a.xhs1.Init(true)
	go a.dy1.Init()
}

// domReady is called after front-end resources have been loaded
func (a App) DomReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) Shutdown(ctx context.Context) {
	// Perform your teardown here
}
