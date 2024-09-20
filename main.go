package main

import (
	"bufio"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/schollz/progressbar/v3"
	"github.com/soxft/busuanzi-sync/bsz_origin"
	"github.com/soxft/busuanzi-sync/config"
	"github.com/soxft/busuanzi-sync/redisutil"
	"github.com/soxft/busuanzi-sync/sitemap"
	"github.com/soxft/busuanzi-sync/transfer"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type Task struct {
	Url     string
	Retries int
}

func main() {
	// 安全询问
	//if pass := safeQA(); pass == false {
	//	os.Exit(0)
	//}

	config.Init()
	redisutil.Init()
	transfer.Init()

	// 读取 sitemap
	log.Println("[INFO] 开始读取 sitemap.xml")
	urls, err := sitemap.Get(viper.GetString("SITEMAP_URL"))
	if err != nil {
		log.Fatalf("读取 sitemap 失败: %v", err)
	}

	log.Printf("--- URL: %d , 线程数: %d ---", len(urls), viper.GetInt("THREADS"))

	log.Println("[INFO] 尝试同步 SitePV 与 SiteUV") // 使用 sitemap URL 作为 site 的标识
	if transfer.IsSynced(viper.GetString("SITEMAP_URL")) == false {
		originData, err := bsz_origin.Get(viper.GetString("SITEMAP_URL"))
		if err != nil {
			log.Fatalf("[FAILED] 获取 Site 数据失败, 请检查网络连接或重试: %v", err)
		}

		transfer.SyncSite(viper.GetString("SITEMAP_URL"), originData)
		transfer.SetSynced(viper.GetString("SITEMAP_URL"))
	}

	log.Println("[INFO] 尝试同步 PagePV 与 PageUV")
	// 初始化 进度条
	bar := progressbar.NewOptions(len(urls),
		progressbar.OptionSetDescription("Processing..."),
		progressbar.OptionShowCount(),
	)

	var wg sync.WaitGroup
	wg.Add(len(urls))
	// 线程池
	var pool *ants.PoolWithFunc
	pool, _ = ants.NewPoolWithFunc(viper.GetInt("THREADS"), func(data interface{}) {
		url := data.(Task)
		if url.Retries >= viper.GetInt("MAX_RETRY") {
			wg.Done()
			log.Printf("\n[FAIL]: %s > 重试 %d/%d, 超过最大尝试次数", url.Url, url.Retries, viper.GetInt("MAX_RETRY"))
			return
		}

		originData, err := bsz_origin.Get(url.Url)
		if err != nil {
			log.Printf("\n[FAIL]: %s > 重试 %d/%d, %v", url.Url, url.Retries, viper.GetInt("MAX_RETRY"), err)
			// 重新加入 队列
			url.Retries++
			_ = pool.Invoke(url)
			return
		}

		transfer.SyncPage(url.Url, originData)
		transfer.SetSynced(url.Url)
		_ = bar.Add(1)
		wg.Done()
	})

	defer pool.Release()

	// 同步 Page 数据
	for _, url := range urls {
		// 如果已经同步过，直接跳过
		if transfer.IsSynced(url.Loc) {
			_ = bar.Add(1)
			wg.Done()
			continue
		}

		if err := pool.Invoke(Task{
			Url:     url.Loc,
			Retries: 1,
		}); err != nil {
			log.Printf("[FAIL] POOL error: %v", err)
		}
	}

	// 等待所有任务完成
	wg.Wait()

	fmt.Println("")
	log.Println("--- 所有任务已完成 ---")
}

func safeQA() bool {
	reader := bufio.NewReader(os.Stdin)
	log.Println("执行此脚本请务必提前手动备份 REDIS (dump.rdb)")
	log.Println("继续操作将可能造成不可逆的数据丢失")
	log.Print("我确认已经备份数据库了 (y/N): ")

	input, err := reader.ReadString('\n')
	if err != nil {
		log.Println("读取输入时出错:", err)
		return false
	}

	pass := false
	// 去除输入的换行符和空格
	input = strings.TrimSpace(input)
	switch input {
	case "y", "Y":
		log.Println("--- 开始执行脚本 ---")
		time.Sleep(3 * time.Second)
		pass = true
	default:
		log.Println("--- 请先备份数据库 ---")
	}

	return pass
}
