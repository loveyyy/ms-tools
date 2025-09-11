package dy

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"testing"
)

func TestHome(t *testing.T) {
	fmt.Println(gfile.MainPkgPath())
	//dy := NewDy()
	//var params = g.Map{
	//	"count":         "20",
	//	"refresh_index": "7",
	//	"refer_id":      "",
	//	"refer_type":    "10",
	//	"pull_type":     "2",
	//}
	//homeList, err := dy.HomeList(params)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(homeList)
}

func TestHot(t *testing.T) {
	dy := NewDy()
	var params = g.Map{
		"detail_list":          "1",
		"source":               "6",
		"main_billboard_count": "5",
	}
	hotList, err := dy.HotList(params)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hotList)
}

func TestSearch(t *testing.T) {
	dy := NewDy()
	var params = g.Map{
		"search_channel":       "aweme_video_web",
		"enable_history":       "1",
		"keyword":              "苹果2025秋季发布会",
		"search_source":        "switch_tab",
		"query_correct_type":   "1",
		"is_filter_search":     "0",
		"from_group_id":        "",
		"disable_rs":           "1",
		"offset":               "0",
		"count":                "10",
		"need_filter_settings": "1",
		"list_type":            "single",
	}
	hotList, err := dy.Search("video", params)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hotList)
}
