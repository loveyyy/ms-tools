package dy

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"ms-tools/backend/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/tjfoc/gmsm/sm3"
)

type Dy struct {
	client *http.HttpClient
	name   string
}

func NewDy() (dy *Dy) {
	dy = &Dy{client: http.NewHttp("https://www.douyin.com/", "dy"), name: "dy"}
	dy.init()
	return
}

func (dy *Dy) log() *glog.Logger {
	return g.Log(dy.name)
}

const aid = 6383
const pageId = 6241

var baseQuery = g.Map{
	"device_platform":     "webapp",
	"aid":                 "6383",
	"channel":             "channel_pc_web",
	"module_id":           "3003101",
	"pc_client_type":      "1",
	"pc_libra_divert":     "Windows",
	"update_version_code": "170400",
	"support_h265":        "1",
	"support_dash":        "1",
	"version_code":        "170400",
	"version_name":        "17.4.0",
	"cookie_enabled":      "true",
	"screen_width":        "1680",
	"screen_height":       "1120",
	"browser_language":    "zh-CN",
	"browser_platform":    "Win32",
	"browser_name":        "Edge",
	"browser_version":     "140.0.0.0",
	"browser_online":      "true",
	"engine_name":         "Blink",
	"engine_version":      "140.0.0.0",
	"os_name":             "Windows",
	"os_version":          "10",
	"cpu_core_num":        "16",
	"device_memory":       "8",
	"platform":            "PC",
	"downlink":            "10",
	"effective_type":      "4g",
	"round_trip_time":     "50",
}

const base64Chars = "ckdp1h4ZKsUB80/Mfvw36XIgR25+WQAlEi7NLboqYTOPuzmFjJnryx9HVGDaStCe"

const base64Chars1 = "Dkdpgh2ZmsQB80/MfvV36XI1R45-WUAlEixNLwoqYTOPuzKFjJnry79HbGcaStCe"

func (dy *Dy) enSm3(data []byte) []byte {
	h := sm3.New()
	h.Write(data)
	return h.Sum(nil)
}

func (dy *Dy) rc4ChangeEncrypt(plaintext []byte, key []byte) []byte {

	// 初始化S盒
	s := make([]int, 256)
	for i := 0; i < 256; i++ {
		s[i] = 255 - i
	}

	// 密钥调度算法
	j := 0
	for i := 0; i < 256; i++ {
		j = (j*s[i] + j + int(key[i%len(key)])) % 256
		// 交换s[i]和s[j]
		s[i], s[j] = s[j], s[i]
	}

	// 生成密钥流并加密
	i := 0
	j = 0
	var cipher []byte
	for k := 0; k < len(plaintext); k++ {
		i = (i + 1) % 256
		j = (j + s[i]) % 256

		// 交换s[i]和s[j]
		s[i], s[j] = s[j], s[i]
		t := (s[i] + s[j]) % 256
		// 异或运算加密
		cipher = append(cipher, plaintext[k]^byte(s[t]))
	}
	return cipher
}

func (dy *Dy) encryptionUa(input []byte, key string) []byte {
	var result []byte
	padding := 0

	// 处理每3个字节
	for i := 0; i < len(input); i += 3 {
		var b1, b2, b3 byte
		b1 = input[i]

		if i+1 < len(input) {
			b2 = input[i+1]
		} else {
			padding++
		}

		if i+2 < len(input) {
			b3 = input[i+2]
		} else {
			padding++
		}

		// 计算Base64编码的4个字符
		c1 := b1 >> 2
		c2 := (b1&3)<<4 | b2>>4
		c3 := (b2&15)<<2 | b3>>6
		c4 := b3 & 63

		result = append(result, key[c1])
		result = append(result, key[c2])

		if padding < 2 {
			result = append(result, key[c3])
		} else {
			result = append(result, '=')
		}

		if padding < 1 {
			result = append(result, key[c4])
		} else {
			result = append(result, '=')
		}
	}

	return result
}

// generateRandomGarbledArrays 生成一个8位随机数组
func (dy *Dy) generateRandomGarbledArrays(num1 byte, num2 byte) []byte {
	var result = []byte{}
	result = append(result, (num1&170)|(1&85))
	result = append(result, (num1&85)|(1&170))
	result = append(result, (num2&170)|(1&85))
	result = append(result, (num2&85)|(1&170))
	return result
}

// generateRandomGarbledArrays 生成一个8位随机数组
func (dy *Dy) generateRandomGarbledArrays1(num1 byte, num2 byte) []byte {
	var result = []byte{}
	result = append(result, (num1&170)|(3&85))
	result = append(result, (num1&85)|(3&170))
	result = append(result, (num2&170)|(82&85))
	result = append(result, (num2&85)|(82&170))
	return result
}

// 日期时间赋值到数组
func (dy *Dy) assignDateTimeValues(arr []byte, ts int64, index int) {
	arr[index] = byte(ts & 255)
	arr[index+1] = byte(ts >> 8 & 255)
	arr[index+2] = byte(ts >> 16 & 255)
	arr[index+3] = byte(ts >> 24 & 255)
	arr[index+4] = byte(ts / 256 / 256 / 256 / 256 & 255)
	arr[index+5] = byte(ts / 256 / 256 / 256 / 256 / 256 & 255)
}

// buildSecondaryArray 根据 indices 提取 arr_55 中的元素
func (dy *Dy) buildSecondaryArray(arr_55 []byte) []byte {
	indices := []int{9, 18, 28, 32, 45, 4, 44, 19, 10, 23, 12, 37, 24, 39, 3, 22, 35, 21, 5, 42, 1, 27, 6, 40, 30, 14, 33, 34, 2, 43, 15, 45, 29, 25, 16, 13, 8, 38, 26, 17, 36, 20, 11, 0, 31, 7, 46, 47, 48, 49}
	result := make([]byte, len(indices))

	for i, index := range indices {
		if index < len(arr_55) {
			result[i] = arr_55[index]
		}
	}

	return result
}

// getLastNum 根据 dateTime1 计算时间戳数组
func (dy *Dy) getLastNum(dateTime1 int64) []byte {
	dateTime1 += 3

	// 计算时间戳
	timestamp1 := dateTime1 & 255

	// 将时间戳转换为字符串并获取其字符的 ASCII 码
	timestampStr := strconv.Itoa(int(timestamp1)) + ","
	num := make([]byte, len(timestampStr))

	for i, char := range timestampStr {
		num[i] = byte(char)
	}

	return num
}

// GenArrObj 生成最终的数组对象
func (dy *Dy) genArrObj(ts1 int64, ts2 int64, e1 []byte, e2 []byte, ua2 []byte, shebei []byte) []byte {
	// 生成第一个数组
	r1 := rand.Float64()
	arr := dy.generateRandomGarbledArrays(byte(int(r1*65535)&255), byte(int(r1*65535)>>8&255))
	// 生成第二个数组
	arr2 := dy.generateRandomGarbledArrays(byte(int(rand.Float64()*240)>>0), byte(int(rand.Float64()*255)>>0&77|2|16|32|128))

	arr8 := append(arr, arr2...)

	last3 := dy.getLastNum(ts1)

	arr50 := make([]byte, 50)
	arr50[0] = 41
	arr50[1] = byte((time.Now().UnixMilli() - 1721836800000) / 1000 / 60 / 60 / 24 / 14 >> 0)
	arr50[2] = 6
	arr50[3] = byte((ts1 - ts2) & 255)
	dy.assignDateTimeValues(arr50, ts1, 4)
	arr50[10] = 1
	arr50[11] = 0
	arr50[12] = 1
	arr50[13] = 0

	arr50[14] = 28
	arr50[15] = 3
	arr50[16] = 3
	arr50[17] = 0

	// 0.00390625, 1, 0
	arr50[18] = 0
	arr50[19] = 0
	arr50[20] = 0
	arr50[21] = 0

	arr50[22] = e1[9]
	arr50[23] = e1[18]
	arr50[24] = e1[3]

	arr50[25] = e2[10]
	arr50[26] = e2[19]
	arr50[27] = e2[4]
	arr50[28] = ua2[11]
	arr50[29] = ua2[21]
	arr50[30] = ua2[5]
	dy.assignDateTimeValues(arr50, ts2-1, 31)
	arr50[37] = 3
	arr50[38] = pageId & 255
	arr50[39] = pageId >> 8 & 255
	arr50[40] = pageId >> 16 & 255
	arr50[41] = pageId >> 24 & 255
	arr50[42] = aid & 255
	arr50[43] = aid >> 8 & 255
	arr50[44] = aid >> 16 & 255
	arr50[45] = aid >> 24 & 255
	arr50[46] = byte(len(shebei) & 255)
	arr50[47] = byte(len(shebei) & 255 >> 8)
	arr50[48] = byte(len(last3) & 255)
	arr50[49] = byte(len(last3) & 255 >> 8)

	arr1 := dy.buildSecondaryArray(arr50)
	arr1 = append(arr1, shebei...)
	arr1 = append(arr1, last3...)

	xorResult := arr8[0] ^ arr8[1] ^ arr8[2] ^ arr8[3] ^ arr8[4] ^ arr8[5] ^ arr8[6] ^ arr8[7] ^ arr1[0] ^ arr1[1] ^ arr1[2] ^ arr1[3] ^ arr1[4] ^ arr1[5] ^ arr1[6] ^ arr1[7] ^ arr1[8] ^ arr1[9] ^ arr1[10] ^ arr1[11] ^ arr1[12] ^ arr1[13] ^ arr1[14] ^ arr1[15] ^ arr1[16] ^ arr1[17] ^ arr1[18] ^ arr1[19] ^ arr1[20] ^ arr1[21] ^ arr1[22] ^ arr1[23] ^ arr1[24] ^ arr1[25] ^ arr1[26] ^ arr1[27] ^ arr1[28] ^ arr1[29] ^ arr1[30] ^ arr1[31] ^ arr1[32] ^ arr1[33] ^ arr1[34] ^ arr1[35] ^ arr1[36] ^ arr1[37] ^ arr1[38] ^ arr1[39] ^ arr1[40] ^ arr1[41] ^ arr1[42] ^ arr1[43] ^ arr1[44] ^ arr1[45] ^ arr1[46] ^ arr1[47] ^ arr1[48] ^ arr1[49]
	arr1 = append(arr1, xorResult)

	numList := []byte{}
	lenArrAr := len(arr1)

	for i := 0; i < lenArrAr; i += 3 {
		num1 := arr1[i]
		has2 := true
		has3 := true
		var num2, num3 byte
		if i+1 < lenArrAr {
			num2 = arr1[i+1]
		} else {
			has2 = false
		}
		if i+2 < lenArrAr {
			num3 = arr1[i+2]
		} else {
			has3 = false
		}

		if has2 && has3 {
			random := byte(int(rand.Float64()*1000) & 255)
			numList = append(numList,
				(random&145)|(num1&110),
				(random&66)|(num2&189),
				(random&44)|(num3&211),
				(num1&145)|(num2&66)|(num3&44),
			)
		} else if has2 {
			numList = append(numList, num1, num2)
		} else {
			numList = append(numList, num1)
		}
	}

	return append(arr8, numList...)
}

func (dy *Dy) getShebeiArray() []byte {
	var data = "660|960|1680|1072|1680|1072|1680|1120|Win32"
	var sum = []byte{}
	for i := range len(data) {
		sum = append(sum, data[i])
	}
	return sum
}

func (dy *Dy) generateBogus() (string, error) {
	ts1 := time.Now().UnixMilli()
	ts2 := time.Now().UnixMilli() - int64(rand.Float64()*10)
	data := "device_platform=webapp&aid=6383&channel=channel_pc_web&sec_user_id=MS4wLjABAAAATsUznK6spH7EAiNGfre4u2yh5Wl8Mkvypp3CO6CPC40&max_cursor=0&locate_query=false&show_live_replay_strategy=1&need_time_list=1&time_list_query=0&whale_cut_token=&cut_version=1&count=18&publish_video_strategy_type=2&from_user_page=1&update_version_code=170400&pc_client_type=1&pc_libra_divert=Windows&support_h265=1&support_dash=1&version_code=290100&version_name=29.1.0&cookie_enabled=true&screen_width=1680&screen_height=1120&browser_language=zh-CN&browser_platform=Win32&browser_name=Edge&browser_version=134.0.0.0&browser_online=true&engine_name=Blink&engine_version=134.0.0.0&os_name=Windows&os_version=10&cpu_core_num=16&device_memory=8&platform=PC&downlink=10&effective_type=4g&round_trip_time=50&webid=7480745848952243763&uifid=164c22db5016193fd69c8bfb0b166ea3a563c2c88054b8eae8759946ea9753ce5c071410576ddf53e072fceef9fc3fe83ce37a827f67b5497098074be511ac36cf9ea681d00f39b4d71609182bb33aba3a4acd7a75dc73816a34de1b871fe51b58970255277df748171bee6a83e7bf6b63066a9d2f58cbf3a2057cfc79199b1a73a06f5c2c822223452675b4385383ab3def42cb1db21f079d311fbca29b4a6f&msToken=RIwH193C_gNF1Ql1Wj5pivuxyq6P66vVs2qX3bc2F1Ibk8MQMI_pqTxgojiL9ul1xi7yvdoTdCCpwImauf0F9AR1nAUTgmzQA1g3fy8WuWIWt8ZnuIGv-59bd5KAWUD3oFM9NmyDmpj-GfkKPdZWFvK-2LMaD0FBV2LRfPSw2Xvd6IZpvtkAkw%3D%3Ddhzx"
	e1 := dy.enSm3([]byte(data))
	e1 = dy.enSm3(e1)

	key := "dhzx"
	k1 := dy.enSm3([]byte(key))
	k1 = dy.enSm3(k1)

	useAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0"
	unicodeValues := []float64{0.00390625, 1, 0}
	var characters []byte
	for _, value := range unicodeValues {
		characters = append(characters, byte(value))
	}
	ua := dy.rc4ChangeEncrypt([]byte(useAgent), characters)
	log.Println("useAgent value is:", useAgent, len(useAgent))

	ua1 := dy.encryptionUa(ua, base64Chars)
	log.Println("ua1 value is:", string(ua1), len(ua1))
	//
	ua2 := dy.enSm3(ua1)
	log.Println("ua2 value is: ", ua2, len(ua2))

	arr1 := dy.getShebeiArray()
	log.Println("arr1 value is:", arr1, len(arr1))

	result := dy.genArrObj(ts1, ts2, e1, k1, ua2, arr1)
	log.Println("Generated Array:", result, len(result))

	result1 := dy.rc4ChangeEncrypt(result, []byte{211})
	log.Println("result1 value is:", result1, len(result1))

	r1 := rand.Float64()
	r2 := rand.Float64()
	arr := dy.generateRandomGarbledArrays1(byte(int(r1*65535)&255), byte(int(r2*40)>>0))
	var result3 = []byte{}
	result3 = append(result3, arr...)
	result3 = append(result3, result1...)
	log.Println("result3 value is:", r1, r2, result3, len(result3))

	abougs := dy.encryptionUa(result3, base64Chars1)
	log.Println("result2", string(abougs), len(abougs))
	return string(abougs), nil
}

func (dy *Dy) init() {
	//获取初始化ttwid
	url1 := "https://www.douyin.com/jingxuan"
	req, err := dy.client.GetNewRequest("GET", url1, nil)
	if err != nil {
		dy.log().Errorf(context.Background(), "init ttwid err %v", err.Error())
		return
	}
	req.Header.Add("Referer", "https://www.douyin.com/")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0")
	_, err = dy.client.SendRequestByReq(req, dy.name)
	if err != nil {
		dy.log().Errorf(context.Background(), "init ttwid  err %v", err.Error())
		return
	}
	dy.log().Infof(context.Background(), "init ttwid success %v", dy.client.GetCookiesStr(dy.name))
	//active ttwid
	url1 = "https://www.douyin.com/ttwid/check/"
	data := g.Map{
		"aid":     6383,
		"service": "www.douyin.com",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		dy.log().Errorf(context.Background(), "active ttwid  Marshal err %v", err.Error())
		return
	}
	req, err = dy.client.GetNewRequest("POST", url1, strings.NewReader(string(jsonData)))
	if err != nil {
		dy.log().Errorf(context.Background(), "active ttwid GetNewRequest err %v", err.Error())
		return
	}
	req.Header.Add("Referer", "https://www.douyin.com/")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0")
	req.Header.Add("Cookie", dy.client.GetCookiesStr(dy.name))
	_, err = dy.client.SendRequestByReq(req, dy.name)
	if err != nil {
		dy.log().Errorf(context.Background(), "active ttwid  err %v", err.Error())
		return
	}
	dy.log().Infof(context.Background(), "active ttwid success %v", dy.client.GetCookiesStr(dy.name))
}

func (dy *Dy) HomeList(params g.Map) (g.Map, error) {
	url1 := "https://www.douyin.com/aweme/v2/web/module/feed/"
	req, err := dy.client.GetNewRequest("POST", url1, nil)
	if err != nil {
		dy.log().Errorf(context.Background(), "HomeList create request error:%v", err.Error())
		return nil, err
	}
	req.Header.Add("Referer", "https://www.douyin.com/jingxuan")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0")
	req.Header.Add("Cookie", dy.client.GetCookiesStr(dy.name))
	req.Form = make(url.Values)
	if len(params) > 0 {
		for s, v := range params {
			req.Form.Add(s, v.(string))
		}
	}
	for s, v := range baseQuery {
		req.Form.Add(s, v.(string))
	}

	response, err := dy.client.SendRequestByReq(req, dy.name)
	if err != nil {
		dy.log().Errorf(context.Background(), "HomeList SendRequestByReq error:%v", err.Error())
		return nil, err
	}
	dy.log().Infof(context.Background(), "HomeList result %v", string(response))
	var result g.Map
	err = json.Unmarshal(response, &result)
	if err != nil {
		dy.log().Errorf(context.Background(), "HomeList Unmarshal error:%v", err.Error())
		return nil, err
	}
	return result, nil
}

func (dy *Dy) HotList(params g.Map) (g.Map, error) {
	url1 := "https://www.douyin.com/aweme/v1/web/hot/search/list/"
	var queryParams = make(url.Values)
	if len(params) > 0 {
		for s, v := range params {
			queryParams.Add(s, v.(string))
		}
	}
	for s, v := range baseQuery {
		queryParams.Add(s, v.(string))
	}
	req, err := dy.client.GetNewRequest("GET", url1+"?"+queryParams.Encode(), nil)
	if err != nil {
		dy.log().Errorf(context.Background(), "hotList create request error:%v", err.Error())
		return nil, err
	}
	req.Header.Add("Referer", "https://www.douyin.com/")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0")
	req.Header.Add("Cookie", dy.client.GetCookiesStr(dy.name))

	response, err := dy.client.SendRequestByReq(req, dy.name)
	if err != nil {
		dy.log().Errorf(context.Background(), "hotList SendRequestByReq error:%v", err.Error())
		return nil, err
	}
	dy.log().Infof(context.Background(), "url:%v,params:%v,resutl:%v", url1, queryParams.Encode(), string(response))
	var result g.Map
	err = json.Unmarshal(response, &result)
	if err != nil {
		dy.log().Errorf(context.Background(), "hotList Unmarshal error:%v", err.Error())
		return nil, err
	}
	return result, nil
}

func (dy *Dy) Search(tab string, params g.Map) (g.Map, error) {
	//abogus, err := dy.generateBogus()
	//if err != nil {
	//	dy.log().Errorf(context.Background(), "hotList create request error:%v", err.Error())
	//	return nil, err
	//}

	//综合 https://www.douyin.com/aweme/v1/web/general/search/stream/ normal_search
	//视频 https://www.douyin.com/aweme/v1/web/search/item/ normal_search
	//用户 https://www.douyin.com/aweme/v1/web/discover/search/ normal_search

	url1 := "https://www.douyin.com/aweme/v1/web/general/search/stream/"
	if strings.EqualFold(tab, "video") {
		url1 = "https://www.douyin.com/aweme/v1/web/search/item/"
	}

	if strings.EqualFold(tab, "user") {
		url1 = " https://www.douyin.com/aweme/v1/web/discover/search/"
	}

	var queryParams = make(url.Values)
	if len(params) > 0 {
		for s, v := range params {
			queryParams.Add(s, v.(string))
		}
	}
	for s, v := range baseQuery {
		queryParams.Add(s, v.(string))
	}
	req, err := dy.client.GetNewRequest("GET", url1+"?"+queryParams.Encode(), nil)
	if err != nil {
		dy.log().Errorf(context.Background(), "Search create request error:%v", err.Error())
		return nil, err
	}
	req.Header.Add("Referer", "https://www.douyin.com/jingxuan/search/")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0")
	req.Header.Add("Cookie", dy.client.GetCookiesStr(dy.name))
	response, err := dy.client.SendRequestByReq(req, dy.name)
	if err != nil {
		dy.log().Errorf(context.Background(), "Search SendRequestByReq error:%v", err.Error())
		return nil, err
	}
	dy.log().Infof(context.Background(), "url:%v,params:%v,resutl:%v", url1, queryParams.Encode(), string(response))
	var result g.Map
	if err := json.Unmarshal(response, &result); err != nil {
		dy.log().Errorf(context.Background(), "Search Unmarshal error:%v", err.Error())
		return nil, err
	}
	return result, nil
}
