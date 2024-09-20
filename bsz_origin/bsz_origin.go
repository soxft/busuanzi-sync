package bsz_origin

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"regexp"
	"time"
)

// Get data from origin busuanzi
func Get(url string) (BszOriginData, error) {
	client := resty.New().SetTimeout(5 * time.Second).SetRetryCount(5).SetHeaders(map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Referer":    url,
	})

	resp, err := client.R().Get("https://busuanzi.ibruce.info/busuanzi?jsonpCallback=callback")
	if err != nil {
		return BszOriginData{}, err
	}

	if resp.RawResponse.StatusCode != 200 {
		return BszOriginData{}, fmt.Errorf("http request Failed, status code: %d", resp.RawResponse.StatusCode)
	}

	// 使用正则表达式提取 JSON 部分
	re := regexp.MustCompile(`callback\((.*?)\);`)
	matches := re.FindStringSubmatch(resp.String())

	if len(matches) < 2 {
		return BszOriginData{}, fmt.Errorf("no JSON data found, raw: %s", resp.String())
	}

	var originData BszOriginData
	err = json.Unmarshal([]byte(matches[1]), &originData)
	if err != nil {
		return BszOriginData{}, fmt.Errorf("json unmarshal failed: %v", err)
	}

	// 去掉本次请求的统计
	originData.SitePv -= 1
	originData.SiteUv -= 1
	originData.PagePv -= 1

	return originData, nil
}
