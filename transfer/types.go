package transfer

type RKeys struct {
	SitePvKey  string
	SiteUvKey  string
	PagePvKey  string
	PageUvKey  string
	SiteUnique string
	PathUnique string
}

type Counts struct {
	SitePv int64
	SiteUv int64
	PagePv int64
	PageUv int64
}
