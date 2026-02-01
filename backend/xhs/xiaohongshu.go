package xhs

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gogf/gf/v2/os/gfile"
	"log"
	"math/big"
	mrand "math/rand"
	"ms-tools/backend/http"
	"ms-tools/backend/utils"
	"net/url"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
)

type Xhs struct {
	client       *http.HttpClient
	name         string
	Xsecappid    string
	WebBuild     string
	A1           string
	signName     int
	loginSession string
}

type XhsType int

const (
	ACTIVE       XhsType = 0
	HOME         XhsType = 1
	HOME_FEED    XhsType = 2
	HOME_FEED_V1 XhsType = 21
	FEDD_INFO    XhsType = 4
	SEARCH       XhsType = 5
	USER_INFO    XhsType = 6
	USER_FEED    XhsType = 7
)

type BaseResponse struct {
	Code    int    `json:"code" example:"0" description:"状态码"`
	Data    g.Map  `json:"data,omitempty" description:"数据集"`
	Msg     string `json:"msg,omitempty" example:"操作成功" description:"提示消息"`
	Success bool   `json:"success" example:"true" description:"是否成功"`
}

func NewXhs() (xhs *Xhs) {
	xhs = &Xhs{client: http.NewHttp("https://www.xiaohongshu.com/", "xhs"), name: "xhs", signName: 0,
		loginSession: ""}
	return
}

func (x *Xhs) log() *glog.Logger {
	return g.Log(x.name)
}

// 假设的getPlatformCode函数
func (x *Xhs) getPlatformCode(e interface{}) int {
	// 这里实现获取平台代码的逻辑
	switch e {
	case "Android":
		return 2
	case "iOS":
		return 1
	case "Mac OS":
		return 3
	case "Linux":
		return 4
	default:
		return 5
	}
}

// 假设的genRandomString函数
func (x *Xhs) genRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz1234567890"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		ret[i] = letters[num.Int64()]
	}
	return string(ret)
}

// CRC32加密函数
func (x *Xhs) crc32(data string) uint32 {
	var (
		table [256]uint32
		crc   = ^uint32(0)
	)

	// Generate the CRC32 table
	for i := 0; i < 256; i++ {
		r := uint32(i)
		for s := 0; s < 8; s++ {
			if r&1 == 1 {
				r = 0xedb88320 ^ (r >> 1)
			} else {
				r >>= 1
			}
		}
		table[i] = r
	}

	// Calculate the CRC32 checksum
	for _, b := range data {
		crc = table[byte(crc)^byte(b)] ^ (crc >> 8)
	}

	return ^crc
}

func (x *Xhs) createSearchIdWithBigInt() string {
	// 获取当前时间戳（毫秒）
	t := time.Now().UnixNano() / int64(time.Millisecond)
	// 生成随机数
	n := mrand.Intn(2147483646) + 1
	// 使用big.Int
	tBig := big.NewInt(t)
	nBig := big.NewInt(int64(n))

	// 左移64位：tBig <<= 64
	power := big.NewInt(1)
	power.Lsh(power, 64) // 1 << 64
	tBig.Mul(tBig, power)

	// 相加：tBig += nBig
	tBig.Add(tBig, nBig)

	// 转换为36进制字符串
	return tBig.Text(36)
}

func (x *Xhs) getloadId(e string) string {
	r := x.getPlatformCode(e)
	i := fmt.Sprintf("%x%s%d%s%s", time.Now().Unix(), x.genRandomString(30), r, "0", "000")
	a := x.crc32(i)
	result := fmt.Sprintf("%s%d", i, a)
	// 如果长度不足52位，则在CRC32部分前面补0
	if len(result) < 52 {
		// 计算需要补多少个0
		neededZeros := 52 - len(result)
		if neededZeros > 0 {
			// 在i和a之间插入需要的0
			result = i + strings.Repeat("0", neededZeros) + fmt.Sprintf("%d", a)
		}
	}
	// 最终确保长度为52位
	if len(result) > 52 {
		result = result[:52]
	}
	return result
}

type XSTruct struct {
	XS      string `json:"X-s"`
	XT      int64  `json:"X-t"`
	XCommon string `json:"X-common"`
}

func (x *Xhs) xsts(path string, params string) *XSTruct {
	// 调用xhsSign函数
	xs := x.xTS(path, params)
	if xs == nil {
		x.log().Errorf(context.Background(), "xsts XTS 异常")
		return nil
	}
	xsc := x.xsCommon(map[string]interface{}{
		"platform":  "Windows",
		"xsecappid": x.Xsecappid,
		"webBuild":  x.WebBuild,
		"a1":        x.A1,
		"X-S":       xs.XS,
		"X-T":       xs.XT,
	})
	if len(xsc) == 0 {
		x.log().Errorf(context.Background(), "xsts xsCommon 异常")
		return nil
	}
	xs.XCommon = xsc
	return xs
}

func (x *Xhs) xsts1(path string, params string) *XSTruct {
	// 调用xhsSign函数
	xs := x.xTS1(path, params)
	if xs == nil {
		x.log().Errorf(context.Background(), "xsts XTS 异常")
		return nil
	}
	xsc := x.xsCommon(map[string]interface{}{
		"platform":  "Windows",
		"xsecappid": x.Xsecappid,
		"webBuild":  x.WebBuild,
		"a1":        x.A1,
		"X-S":       xs.XS,
		"X-T":       xs.XT,
	})
	if len(xsc) == 0 {
		x.log().Errorf(context.Background(), "xsts xsCommon 异常")
		return nil
	}
	xs.XCommon = xsc
	return xs
}

func (x *Xhs) profileData() string {
	jsFilePath := utils.GetRootPath() + "/js/xhs/gid.js"
	if !gfile.Exists(jsFilePath) {
		x.log().Errorf(context.Background(), "xhs gid.js 文件不存在")
		return ""
	}
	cmd := exec.Command("node", jsFilePath, x.client.GetCookiesStr())
	// 设置SysProcAttr属性，实现静默执行
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	output, err := cmd.Output()
	if err != nil {
		x.log().Errorf(context.Background(), "Command err %v", err.Error())
		return ""
	}
	x.log().Infof(context.Background(), "profileData 返回的值为:%v", string(output))
	return string(output)
}

func (x *Xhs) xsCommon(e map[string]interface{}) string {
	var (
		_, _ interface{}
		s    = e["platform"].(string)
	)
	c := e["X-T"].(int64)
	d := e["X-S"].(string)

	h := "I38rHdgsjopgIvesdVwgIC+oIELmBZ5e3VwXLgFTIxS3bqwErFeexd0ekncAzMFYnqthIhJeSfMDKutRI3KsYorWHPtGrbV0P9WfIi/eWc6eYqtyQApPI37ekmR6QL+5Ii6sdneeSfqYHqwl2qt5B0DBIx+PGDi/sVtkIxdsxuwr4qtiIhuaIE3e3LV0I3VTIC7e0utl2ADmsLveDSKsSPw5IEvsiVtJOqw8BuwfPpdeTFWOIx4TIiu6ZPwrPut5IvlaLbgs3qtxIxes1VwHIkumIkIyejgsY/WTge7eSqte/D7sDcpipedeYrDtIC6eDVw2IENsSqtlnlSuNjVtIvoekqt3cZ7sVo4gIESyIhE4QfquIxhnqz8gIkIfoqwkICZWG73sdlOeVPw3IvAe0fgedfR4Ii5s3Ib02utAIiKsidvekZNeTPt4nAOeWPwEIvS8zeoedPwvp9gsSVwrI3SrIxE5Luwwaqw+rekhZANe1MNe0Pw9ICNsVLoeSbIFIkosSr7sVnFiIkgsVVwHIvln/PtCcPwpIEos0uwQIv5e6uw5Ih3sjuw5NqwGoVwuICzTIvRtQeVAGl/siqtFIhPtIieeYuwoeWccpUOsDskuIhRytPwwzqwAIkesWqtuqIAsVF6s1IbLIE0s6edsiPtccPwrICJeWqwvIiNeT/VOICdeS9geSPtlIxvsVqwT2dT+IEGKIv5s6roex97ejD7sSqtoIkWMIxzhmqtoIvkjIk5sxbc3H//sjqtBIEbwZutWJZesSroeYuwseuw/IvEAIiuZ/Vw8GVt2IC7edM6sVPterVtqIiZwLPwkIv8oadF+IE/ed/7eY0OekVwwHPtzIiPcI3Ns6utCIE3edFTfICYyZlM2IvWGICgeYPw0/qt+QPtjIkrTIEAs1j/eSuwxpMDVICWrJPtLbutwI3rXIxQxIEMMIk/sTlNsTDSeJuwwI3S3rqtreuw2Ii6eTAlXIEbTeg0eWutZPqtvB/qtsDZHI3TknuwBI3gsV/YzIkE4sqwtIkEPwqt1HVtoIhu6ICzNIh89sqwNIiGhIvJe1PtvwcH6IEVOIv0sfnKsxuwAICTMIhIU/chFIv4AGYH7IvOsWS/e0qwwIk7s1W=="
	g1 := "1"
	listMap := gmap.NewListMap(true)
	listMap.Set("s0", x.getPlatformCode(s))
	listMap.Set("s1", "")
	listMap.Set("x0", g1)
	listMap.Set("x1", "4.3.1")
	listMap.Set("x2", s)
	listMap.Set("x3", e["xsecappid"].(string))
	listMap.Set("x4", e["webBuild"].(string))
	listMap.Set("x5", e["a1"].(string))
	listMap.Set("x6", c)
	listMap.Set("x7", d)
	listMap.Set("x8", h)
	jsonDataBytes, _ := json.Marshal(listMap)
	log.Print(string(jsonDataBytes))
	jsFilePath := utils.GetRootPath() + "/js/xhs/xscommon.js"
	if !gfile.Exists(jsFilePath) {
		x.log().Errorf(context.Background(), "xhs xscommon.js 文件不存在")
		return ""
	}
	cmd := exec.Command("node", jsFilePath, x.client.GetCookiesStr(), string(jsonDataBytes), "0")
	// 设置SysProcAttr属性，实现静默执行
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	output, err := cmd.Output()
	if err != nil {
		x.log().Errorf(context.Background(), "Command err %v", err.Error())
		return ""
	}

	x.log().Infof(context.Background(), "X-Common 参数:%v 返回的值为:%v", string(jsonDataBytes), string(output))
	return string(output)
}

func (x *Xhs) xTS(u string, params string) *XSTruct {
	jsFilePath := utils.GetRootPath() + "/js/xhs/xs.js"
	if !gfile.Exists(jsFilePath) {
		x.log().Errorf(context.Background(), "xhs xs.js 文件不存在")
		return nil
	}
	x.log().Infof(context.Background(), "XTS请求的URL为： %v,%v,%v,%v", u, params, jsFilePath, x.client.GetCookiesStr())
	cmd := exec.Command("node", jsFilePath, x.client.GetCookiesStr(), u, params)
	// 设置SysProcAttr属性，实现静默执行
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	output, err := cmd.Output()
	if err != nil {
		x.log().Errorf(context.Background(), "Command err %v", err.Error())
		return nil
	}
	var xs *XSTruct
	err = json.Unmarshal(output, &xs)
	if err != nil {
		x.log().Errorf(context.Background(), "Unmarshal err %v", err.Error())
		return nil
	}
	g.Log().Infof(context.Background(), "L返回的值为： %v,%v", xs.XT, xs.XS)
	return xs
}

func (x *Xhs) xTS1(u string, params string) *XSTruct {
	jsFilePath := utils.GetRootPath() + "/js/xhs/xs2.js"
	if !gfile.Exists(jsFilePath) {
		x.log().Errorf(context.Background(), "xhs xs1.js 文件不存在")
		return nil
	}
	x.log().Infof(context.Background(), "XTS请求的URL为： %v,%v,%v,%v", x.client.GetCookiesStr(), u, params, jsFilePath)
	cmd := exec.Command("node", jsFilePath, x.client.GetCookiesStr(), u, params)
	// 设置SysProcAttr属性，实现静默执行
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	output, err := cmd.Output()
	if err != nil {
		x.log().Errorf(context.Background(), "Command err %v", err.Error())
		return nil
	}
	nString := strings.ReplaceAll(string(output), "{}", "")
	nString = strings.ReplaceAll(nString, "\n", "")
	xs := &XSTruct{
		XS: strings.ReplaceAll(nString, "[Error]", ""),
		XT: gtime.Now().UnixMilli(),
	}
	g.Log().Infof(context.Background(), "L返回的值为： %v,%v", xs.XT, strings.ReplaceAll(xs.XS, "[Error]", ""))
	return xs
}

func (x *Xhs) Init(isInit bool) {
	if isInit {
		x.Xsecappid = "xhs-pc-web"
		x.client.AddCookies("xsecappid", x.Xsecappid)

		x.WebBuild = "5.8.0"
		x.client.AddCookies("webBuild", x.WebBuild)

		e := "WINDOWS"
		x.A1 = x.getloadId(e)
		x.client.AddCookies("a1", x.A1)

		//webId   md5 加密
		hashed := md5.Sum([]byte(x.A1))
		webId := hex.EncodeToString(hashed[:])
		x.client.AddCookies("webId", webId)
		//去登录下获取 web_session
		x.login()
		profileData := x.profileData()
		if len(profileData) == 0 {
			x.log().Errorf(context.Background(), "xhsGId err")
			return
		}
		x.getGID(profileData)
	} else {
		cookiesMap, err := x.client.LoadCookieFromFile(x.name)
		if err == nil {
			x.Xsecappid = cookiesMap["xsecappid"].(string)
			x.WebBuild = cookiesMap["webBuild"].(string)
			x.A1 = cookiesMap["a1"].(string)
		}
	}
}

func (x *Xhs) GetData(xhsType XhsType, params g.Map) (g.Map, error) {
	var result g.Map
	var err error
	switch xhsType {
	case HOME:
		result, err = x.home()
	case HOME_FEED:
		result, err = x.homeFeed(params, false)
	case HOME_FEED_V1:
		result, err = x.homeFeed(params, true)
	case SEARCH:
		result, err = x.search(params)
	case FEDD_INFO:
		result, err = x.feedDetail(params, true)
	case USER_INFO:
		result, err = x.userDetail(params)
	case USER_FEED:
		result, err = x.userFeedList(params)
	default:
		return nil, errors.New("未知的XhsType")
	}
	return result, err
}

func (x *Xhs) login() {
	//登录
	// web_session cookies
	url := "https://edith.xiaohongshu.com/api/sns/web/v1/login/activate"
	params := `{}`

	request, err := x.client.GetNewRequest("POST", url, strings.NewReader(params))
	if err != nil {
		x.log().Errorf(context.Background(), "activate err %v", err.Error())
		return
	}

	// 调用xhsSign函数
	xsts := x.xsts("/api/sns/web/v1/login/activate", params)
	if xsts == nil {
		x.log().Errorf(context.Background(), "Login XTS 异常")
		return
	}
	request.Header.Add("X-S", xsts.XS)
	request.Header.Add("X-S-Common", strings.Trim(xsts.XCommon, "\n"))
	request.Header.Add("X-T", strconv.FormatInt(xsts.XT, 10))
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0")
	request.Header.Add("Accept", "application/json, text/plain, */*")
	request.Header.Add("Content-Type", "application/json;charset=UTF-8")
	request.Header.Add("Referer", "https://www.xiaohongshu.com/")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Cookie", x.client.GetCookiesStr())
	response, err := x.client.SendRequestByReq(request, x.name)
	if err != nil {
		x.log().Errorf(context.Background(), "Login Send 异常:%v", err)
		return
	}
	x.log().Infof(context.Background(), "url:%v,params:%v,resutl:%v", url, params, string(response))
}

func (x *Xhs) getGID(profileData string) {
	//上传设备信息
	// gid
	url := "https://as.xiaohongshu.com/api/sec/v1/shield/webprofile"
	var params = g.Map{}
	params["platform"] = "Windows"
	params["profileData"] = profileData
	params["sdkVersion"] = "4.3.1"
	params["svn"] = 2
	jsonParams, err := json.Marshal(params)
	if err != nil {
		x.log().Errorf(context.Background(), "Marshal err %v", err.Error())
		return
	}
	request, err := x.client.GetNewRequest("POST", url, strings.NewReader(string(jsonParams)))
	if err != nil {
		x.log().Errorf(context.Background(), "activate err %v", err.Error())
		return
	}

	// 调用xhsSign函数
	xsts := x.xsts("/api/sec/v1/shield/webprofile", string(jsonParams))
	if xsts == nil {
		x.log().Errorf(context.Background(), "webprofile XTS 异常")
		return
	}
	request.Header.Add("X-S", xsts.XS)
	request.Header.Add("X-S-Common", strings.Trim(xsts.XCommon, "\n"))
	request.Header.Add("X-T", strconv.FormatInt(xsts.XT, 10))
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0")
	request.Header.Add("Accept", "application/json, text/plain, */*")
	request.Header.Add("Content-Type", "application/json;charset=UTF-8")
	request.Header.Add("Referer", "https://www.xiaohongshu.com/")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Cookie", x.client.GetCookiesStr())
	response, err := x.client.SendRequestByReq(request, x.name)
	if err != nil {
		x.log().Errorf(context.Background(), "webprofile Send 异常:%v", err)
		return
	}
	x.log().Infof(context.Background(), "url:%v,params:%v,resutl:%v", url, params, string(response))
}

// 获取分类和初始化cookies
func (x *Xhs) home() (g.Map, error) {
	request, err := x.client.GetNewRequest("GET", "https://www.xiaohongshu.com/explore", nil)
	if err != nil {
		x.log().Errorf(context.Background(), "userDetail GetNewRequest er %v", err.Error())
		return nil, err
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0")
	request.Header.Add("Referer", "https://www.xiaohongshu.com/explore")
	request.Header.Add("Cookie", x.client.GetCookiesStr())
	result, err := x.client.SendRequestByReq(request, x.name)
	if err != nil {
		x.log().Errorf(context.Background(), "userDetail Send 异常:%v", err)
		return nil, err
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0")
	request.Header.Add("Referer", "https://www.xiaohongshu.com/explore")
	request.Header.Add("Cookie", x.client.GetCookiesStr())
	x.log().Infof(context.Background(), "https://www.xiaohongshu.com/explore   %v %v", nil, string(result))
	reader := strings.NewReader(string(result))
	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		x.log().Errorf(context.Background(), "parse html error:%v", err)
		return nil, err
	}
	homeInfoMap := g.Map{}
	var categoryList []g.Map
	document.Find(".content-container div").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("id")
		x.log().Infof(context.Background(), "id:%v,text:%v", name, s.Text())
		categoryList = append(categoryList, g.Map{
			"name": strings.TrimSpace(s.Text()),
			"id":   name,
		})
	})
	homeInfoMap["category"] = categoryList
	x.log().Infof(context.Background(), "homeInfoMap:%v", homeInfoMap)
	return homeInfoMap, err
}

func (x *Xhs) homeFeed(params1 g.Map, isV2 bool) (g.Map, error) {
	u := "https://edith.xiaohongshu.com/api/sns/web/v1/homefeed"
	jsonParams, err := json.Marshal(params1)
	if err != nil {
		x.log().Errorf(context.Background(), "Marshal err %v", err.Error())
		return nil, err
	}

	request, err := x.client.GetNewRequest("POST", u, strings.NewReader(string(jsonParams)))
	if err != nil {
		x.log().Errorf(context.Background(), "GetNewRequest err %v", err.Error())
		return nil, err
	}
	// 调用xhsSign函数
	var xsts *XSTruct
	if isV2 {
		xsts = x.xsts1("/api/sns/web/v1/homefeed", string(jsonParams))
	} else {
		xsts = x.xsts("/api/sns/web/v1/homefeed", string(jsonParams))
	}
	if xsts == nil {
		x.log().Errorf(context.Background(), "HomeFeed XTS 异常")
		return nil, err
	}
	request.Header.Add("X-s", xsts.XS)
	request.Header.Add("X-S-Common", strings.Trim(xsts.XCommon, "\n"))
	request.Header.Add("X-t", strconv.FormatInt(xsts.XT, 10))
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0")
	request.Header.Add("Accept", "application/json, text/plain, */*")
	request.Header.Add("Content-Type", "application/json;charset=UTF-8")
	request.Header.Add("Referer", "https://www.xiaohongshu.com/")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Cookie", x.client.GetCookiesStr())
	result, err := x.client.SendRequestByReq(request, x.name)
	if err != nil {
		x.log().Errorf(context.Background(), "HomeFeed Send 异常:%v", err)
		return nil, err
	}
	x.log().Infof(context.Background(), "url:%v,params:%v,resutl:%v", u, string(jsonParams), string(result))
	var resultMap *BaseResponse
	err = json.Unmarshal(result, &resultMap)
	if err != nil {
		return nil, err
	}
	return resultMap.Data, nil
}

func (x *Xhs) feedDetail(params1 g.Map, isV2 bool) (g.Map, error) {
	u := "https://edith.xiaohongshu.com/api/sns/web/v1/feed"
	jsonParams, err := json.Marshal(params1)
	if err != nil {
		x.log().Errorf(context.Background(), "Marshal err %v", err.Error())
		return nil, err
	}
	request, err := x.client.GetNewRequest("POST", u, strings.NewReader(string(jsonParams)))
	if err != nil {
		x.log().Errorf(context.Background(), "GetNewRequest err %v", err.Error())
		return nil, err
	}
	// 调用xhsSign函数
	var xsts *XSTruct
	if isV2 {
		xsts = x.xsts1("/api/sns/web/v1/feed", string(jsonParams))
	} else {
		xsts = x.xsts("/api/sns/web/v1/feed", string(jsonParams))
	}
	if xsts == nil {
		x.log().Errorf(context.Background(), "feedDetail XTS 异常")
		return nil, err
	}
	request.Header.Add("X-S", xsts.XS)
	request.Header.Add("X-S-Common", strings.Trim(xsts.XCommon, "\n"))
	request.Header.Add("X-T", strconv.FormatInt(xsts.XT, 10))
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0")
	request.Header.Add("Accept", "application/json, text/plain, */*")
	request.Header.Add("Content-Type", "application/json;charset=UTF-8")
	request.Header.Add("Referer", "https://www.xiaohongshu.com/")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Cookie", x.client.GetCookiesStr())
	result, err := x.client.SendRequestByReq(request, x.name)
	if err != nil {
		x.log().Errorf(context.Background(), "feedDetail Send 异常:%v", err)
		return nil, err
	}
	x.log().Infof(context.Background(), "url:%v,params:%v,resutl:%v", u, string(jsonParams), string(result))
	var resultMap *BaseResponse
	err = json.Unmarshal(result, &resultMap)
	if err != nil {
		return nil, err
	}
	return resultMap.Data, nil
}

func (x *Xhs) SetLoginSession(webSession string) {
	x.loginSession = webSession
	x.client.AddCookies("web_session", webSession)
}

func (x *Xhs) IsLoginSession() string {
	return x.loginSession
}

func (x *Xhs) search(params1 g.Map) (g.Map, error) {
	u := "https://edith.xiaohongshu.com/api/sns/web/v1/search/notes"
	params1["search_id"] = x.createSearchIdWithBigInt()
	x.SetLoginSession(x.loginSession)
	jsonParams, err := json.Marshal(params1)
	if err != nil {
		x.log().Errorf(context.Background(), "Marshal err %v", err.Error())
		return nil, err
	}
	request, err := x.client.GetNewRequest("POST", u, strings.NewReader(string(jsonParams)))
	if err != nil {
		x.log().Errorf(context.Background(), "GetNewRequest er %v", err.Error())
		return nil, err
	}
	// 调用xhsSign函数
	xsts := x.xsts("/api/sns/web/v1/search/notes", string(jsonParams))
	if xsts == nil {
		x.log().Errorf(context.Background(), "HomeFeed XTS 异常")
		return nil, err
	}
	request.Header.Add("X-S", xsts.XS)
	request.Header.Add("X-S-Common", strings.Trim(xsts.XCommon, "\n"))
	request.Header.Add("X-T", strconv.FormatInt(xsts.XT, 10))
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0")
	request.Header.Add("Accept", "application/json, text/plain, */*")
	request.Header.Add("Content-Type", "application/json;charset=UTF-8")
	request.Header.Add("Referer", "https://www.xiaohongshu.com/")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Cookie", x.client.GetCookiesStr())
	result, err := x.client.SendRequestByReq(request, x.name)
	if err != nil {
		x.log().Errorf(context.Background(), "HomeFeed Send 异常:%v", err)
		return nil, err
	}
	x.log().Infof(context.Background(), "url:%v,params:%v,resutl:%v", u, string(jsonParams), string(result))
	var resultMap *BaseResponse
	err = json.Unmarshal(result, &resultMap)
	if err != nil {
		return nil, err
	}
	return resultMap.Data, nil
}

func (x *Xhs) userDetail(params g.Map) (g.Map, error) {
	u := "https://www.xiaohongshu.com/user/profile/" + params["user_id"].(string)
	delete(params, "user_id")
	delete(params, "cursor")
	delete(params, "num")
	delete(params, "image_formats")
	var queryParams = make(url.Values)
	if len(params) > 0 {
		for s, v := range params {
			queryParams.Add(s, v.(string))
		}
	}
	x.client.AddCookies("web_session", x.loginSession)
	request, err := x.client.GetNewRequest("GET", u+"?"+queryParams.Encode(), nil)
	if err != nil {
		x.log().Errorf(context.Background(), "userDetail GetNewRequest er %v", err.Error())
		return nil, err
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0")
	request.Header.Add("Referer", "https://www.xiaohongshu.com/explore")
	request.Header.Add("Cookie", x.client.GetCookiesStr())
	result, err := x.client.SendRequestByReq(request, x.name)
	if err != nil {
		x.log().Errorf(context.Background(), "userDetail Send 异常:%v", err)
		return nil, err
	}
	x.log().Infof(context.Background(), "url:%v,params:%v,resutl:%v", u, queryParams.Encode(), string(result))
	reader := strings.NewReader(string(result))
	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		x.log().Errorf(context.Background(), "parse html error:%v", err)
		return nil, err
	}
	var userDetail = g.Map{}
	document.Find("#userPageContainer > div.user > div").Each(func(i int, s *goquery.Selection) {
		imgUrl, _ := s.Find("div.avatar > div > img").Attr("src")
		userDetail["avatar_url"] = imgUrl
		userName := s.Find("div.info-part > div.info > div.basic-info > div.user-basic > div.user-nickname > div").Text()
		userDetail["nick_name"] = userName
		id := s.Find("div.info-part > div.info > div.basic-info > div.user-basic > div.user-content > span.user-redId").Text()
		userDetail["id"] = id
		ipAddress := s.Find("div.info-part > div.info > div.basic-info > div.user-basic > div.user-content > span.user-IP").Text()
		userDetail["ip_address"] = ipAddress
		desc := s.Find("div.info-part > div.info > div.user-desc").Text()
		userDetail["desc"] = desc
		var tags []string
		s.Find("div.info-part > div.info > div.user-tags > div.tag-item > div").Each(func(i int, s *goquery.Selection) {
			if len(s.Text()) > 0 {
				tags = append(tags, s.Text())
			}
		})
		userDetail["tags"] = tags
		followCount := s.Find("div.info-part > div.info > div.data-info > div > div:nth-child(1) > span.count").Text()
		userDetail["follow_count"] = followCount
		fansCount := s.Find("div.info-part > div.info > div.data-info > div > div:nth-child(2) > span.count").Text()
		userDetail["fans_count"] = fansCount
		likeCount := s.Find("div.info-part > div.info > div.data-info > div > div:nth-child(3) > span.count").Text()
		userDetail["like_count"] = likeCount
	})

	//提取字符串中"channels":至,"isResourceDisplay"之间的字符串
	var reg = regexp.MustCompile(`"notes":(.*?),"isFetchingNotes"`)
	match := reg.FindStringSubmatch(string(result))
	if len(match) > 1 {
		noteList := match[1]
		var feedList []interface{}
		err = json.Unmarshal([]byte(noteList), &feedList)
		if err != nil {
			x.log().Errorf(context.Background(), "json.Unmarshal err %v", err.Error())
			return nil, err
		}
		userDetail["feed_list"] = feedList[0]
	}

	//提取字符串中"channels":至,"isResourceDisplay"之间的字符串
	reg = regexp.MustCompile(`"noteQueries":(.*?),"pageScrolled"`)
	match = reg.FindStringSubmatch(string(result))
	if len(match) > 1 {
		noteList := match[1]
		var feedList []g.Map
		err = json.Unmarshal([]byte(noteList), &feedList)
		if err != nil {
			x.log().Errorf(context.Background(), "json.Unmarshal err %v", err.Error())
			return nil, err
		}
		var feedJson = feedList[0]
		userDetail["hasMore"] = feedJson["hasMore"]
		userDetail["num"] = feedJson["num"]
		userDetail["cursor"] = feedJson["cursor"]
	}
	return userDetail, nil
}

func (x *Xhs) userFeedList(params1 g.Map) (g.Map, error) {
	u := "https://edith.xiaohongshu.com/api/sns/web/v1/user_posted"
	delete(params1, "channel_type")
	delete(params1, "parent_page_channel_type")
	var queryParams = make(url.Values)
	if len(params1) > 0 {
		for s, v := range params1 {
			queryParams.Add(s, v.(string))
		}
	}
	x.client.AddCookies("web_session", x.loginSession)
	request, err := x.client.GetNewRequest("GET", u+"?"+queryParams.Encode(), nil)
	if err != nil {
		x.log().Errorf(context.Background(), "userDetail GetNewRequest er %v", err.Error())
		return nil, err
	}
	// 调用xhsSign函数
	xsts := x.xsts("/api/sns/web/v1/user_posted"+"?"+queryParams.Encode(), "")
	if xsts == nil {
		x.log().Errorf(context.Background(), "HomeFeed XTS 异常")
		return nil, err
	}
	request.Header.Add("X-S", xsts.XS)
	request.Header.Add("X-S-Common", xsts.XCommon)
	request.Header.Add("X-T", strconv.FormatInt(xsts.XT, 10))
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0")
	request.Header.Add("Accept", "application/json, text/plain, */*")
	request.Header.Add("Content-Type", "application/json;charset=UTF-8")
	request.Header.Add("Referer", "https://www.xiaohongshu.com/")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Cookie", x.client.GetCookiesStr())
	result, err := x.client.SendRequestByReq(request, x.name)
	if err != nil {
		x.log().Errorf(context.Background(), "userFeedList Send 异常:%v", err)
		return nil, err
	}
	x.log().Infof(context.Background(), "url:%v,params:%v,resutl:%v", u, queryParams.Encode(), string(result))
	var resultMap *BaseResponse
	err = json.Unmarshal(result, &resultMap)
	if err != nil {
		return nil, err
	}
	return resultMap.Data, nil
}
