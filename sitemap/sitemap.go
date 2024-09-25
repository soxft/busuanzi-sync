package sitemap

import (
	"github.com/gocolly/colly/v2"
	"log"
	"strconv"
)

// Get 函数返回解析后的 URL 列表
func Get(sitemapURL string) ([]URL, error) {
	// 创建 Colly 收集器
	c := colly.NewCollector()

	var urls []URL // 用于存储解析后的 URL 信息

	c.OnXML("//url", func(e *colly.XMLElement) {
		loc := e.ChildText("loc")
		//priority := e.ChildText("priority")
		//changeFreq := e.ChildText("changefreq")
		//lastChange := e.ChildText("lastmod")

		// 将解析的信息添加到切片中
		urls = append(urls, URL{
			Loc: loc,
			// Priority:        parsePriority(priority), // 解析为 float64
			// ChangeFrequency: changeFreq,
			// LastChange:      lastChange,
		})
	})

	// 开始抓取
	err := c.Visit(sitemapURL)
	if err != nil {
		return nil, err // 返回错误
	}

	return urls, nil // 返回解析后的 URL 列表
}

// 辅助函数：将字符串解析为 float64
func parsePriority(priority string) float64 {
	p, err := strconv.ParseFloat(priority, 64)
	if err != nil {
		log.Printf("Error parsing priority: %v, using default value 0.0", err)
		return 0.0 // 如果解析失败，返回默认值 0.0
	}
	return p
}
