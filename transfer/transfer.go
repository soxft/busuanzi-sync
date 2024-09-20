package transfer

import (
	"context"
	"fmt"
	"github.com/soxft/busuanzi-sync/bsz_origin"
	"github.com/soxft/busuanzi-sync/redisutil"
	"github.com/spf13/viper"
	url2 "net/url"
)

var Prefix string

func Init() {
	// 初始化
	Prefix = viper.GetString("REDIS_PREFIX")
}

// IsSynced Page 是否已经同步
// 对于 siteKeys 和 pageKeys  使用 sitemap 的 URL 充当标识
func IsSynced(url string) bool {
	return redisutil.RDB.SIsMember(context.Background(), GetSyncKey(), MD5(url)).Val()
}

// SetSynced 设置 Page 已经同步
// 对于 siteKeys 和 pageKeys  使用 sitemap 的 URL 充当标识
func SetSynced(url string) {
	redisutil.RDB.SAdd(context.Background(), GetSyncKey(), MD5(url))
}

// SyncPage 同步 pageUv 和 PagePV,   << sitePv 和 siteUv 数据需要单独同步 >>
func SyncPage(url string, counts bsz_origin.BszOriginData) {
	u, _ := url2.Parse(url)
	if u.Host == "" {
		return
	}

	ctx := context.Background()
	keys := getKeys(u.Host, u.Path)

	counts.PagePv -= 1

	// pagePv zSet
	redisutil.RDB.ZIncrBy(ctx, keys.PagePvKey, float64(counts.PagePv), keys.PathUnique)

	// 原版没有 pageUv 顾不需要同步
}

// SyncSite 同步 siteUv 和 sitePv
func SyncSite(url string, counts bsz_origin.BszOriginData) {
	u, _ := url2.Parse(url)
	if u.Host == "" {
		return
	}

	ctx := context.Background()
	keys := getKeys(u.Host, u.Path)

	counts.SitePv -= 1
	counts.SiteUv -= 1

	for i := 0; i < int(counts.SiteUv); i++ {
		redisutil.RDB.PFAdd(ctx, keys.SiteUvKey, MD5(fmt.Sprintf("%s:%d", u.Host, i)))
	}
	redisutil.RDB.IncrBy(ctx, keys.SitePvKey, counts.SitePv)
}
