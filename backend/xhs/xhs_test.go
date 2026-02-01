package xhs

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"testing"
)

func TestHome(t *testing.T) {
	x := NewXhs()
	x.Init(false)
	var params = g.Map{}
	params["cursor_score"] = ""
	params["num"] = 18
	params["refresh_type"] = 1
	params["note_index"] = 12
	params["unread_begin_note_id"] = ""
	params["unread_end_note_id"] = ""
	params["category"] = "homefeed.fashion_v3"
	params["search_key"] = ""
	params["need_num"] = 8
	params["image_formats"] = []string{"jpg", "webp", "avif"}
	params["need_filter_image"] = false
	params["unread_note_count"] = 0
	data, err := x.GetData(HOME_FEED_V1, params)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)

	//var items = data["items"].([]interface{})
	//
	//fmt.Println(data["cursor_score"], len(items))
	//
	//params = g.Map{}
	//params["cursor_score"] = data["cursor_score"]
	//params["num"] = 30
	//params["refresh_type"] = 3
	//params["note_index"] = len(items) + 1
	//params["unread_begin_note_id"] = ""
	//params["unread_end_note_id"] = ""
	//params["category"] = "homefeed.fashion_v3"
	//params["search_key"] = ""
	//params["need_num"] = 15
	//params["image_formats"] = []string{"jpg", "webp", "avif"}
	//params["need_filter_image"] = false
	//params["unread_note_count"] = 0
	//data, err = x.GetData(HOME_FEED_V1, params)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(data)

}

func TestSearch(t *testing.T) {
	x := NewXhs()
	var params = g.Map{}
	params["keyword"] = "222222"
	params["page"] = 2
	params["page_size"] = 20
	params["search_id"] = x.createSearchIdWithBigInt()
	params["sort"] = "general"
	params["note_type"] = 0
	params["ext_flags"] = []string{}
	params["geo"] = ""
	params["image_formats"] = []string{"jpg", "webp", "avif"}
	homeInfo, err := x.GetData(SEARCH, params)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(homeInfo)
}

func TestUserInfo(t *testing.T) {
	x := NewXhs()
	var params = g.Map{}
	params["user_id"] = "5b67f8c469d6ce00016515db"
	params["xsec_token"] = "ABo_CisNnINOgb7a1f5llz8GoAkROHdFYg3-oY3b7Rm34="
	params["xsec_source"] = "pc_feed"
	homeInfo, err := x.GetData(USER_INFO, params)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(homeInfo)
}

func TestUserFeed(t *testing.T) {
	x := NewXhs()
	var params = g.Map{}
	params["user_id"] = "65a940d1000000000802035d"
	params["num"] = "30"
	params["cursor"] = "6899c1750000000023037418"
	params["xsec_token"] = "ABlIc4g5Bj4icFSz1rALaKHHfntFZ3mKyjMBkVpyT0IyQ="
	params["xsec_source"] = "pc_feed"
	params["image_formats"] = "jpg,webp,avif"
	homeInfo, err := x.GetData(USER_FEED, params)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(homeInfo)
}
