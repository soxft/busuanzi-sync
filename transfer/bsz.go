package transfer

import (
	"github.com/spf13/viper"
	"strings"
)

// getKeys 获取 keys
// https://github.com/soxft/busuanzi/blob/main/core/count.go
func getKeys(host string, path string) RKeys {
	var siteUnique = host
	var pathUnique = path

	// 兼容旧版本
	if viper.GetBool("BSZ_PATH_STYLE") == false {
		pathUnique = host + "&" + path
	}

	// encrypt
	switch viper.GetString("BSZ_ENCRYPT") {
	case "MD516":
		siteUnique = MD5(siteUnique)[8:24]
		pathUnique = MD5(pathUnique)[8:24]
	case "MD532":
		siteUnique = MD5(siteUnique)
		pathUnique = MD5(pathUnique)
	default:
		siteUnique = MD5(siteUnique)
		pathUnique = MD5(pathUnique)
	}

	redisPrefix := viper.GetString("REDIS_PREFIX")

	siteUvKey := strings.Join([]string{redisPrefix, "site_uv", siteUnique}, ":")
	pageUvKey := strings.Join([]string{redisPrefix, "page_uv", siteUnique, pathUnique}, ":")

	sitePvKey := strings.Join([]string{redisPrefix, "site_pv", siteUnique}, ":")
	pagePvKey := strings.Join([]string{redisPrefix, "page_pv", siteUnique}, ":")

	return RKeys{
		SitePvKey:  sitePvKey,
		SiteUvKey:  siteUvKey,
		PagePvKey:  pagePvKey,
		PageUvKey:  pageUvKey,
		SiteUnique: siteUnique,
		PathUnique: pathUnique,
	}
}
