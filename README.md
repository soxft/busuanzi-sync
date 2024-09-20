# busuanzi-sync

> 这是一个用于将 [原版不蒜子](https://busuanzi.ibruce.info/) 的统计数据同步到 [soxft/busuanzi](https://github.com/soxft/busuanzi) 的工具。


## 使用方法

1. 按照提示填写 .env 文件
2. 运行脚本 `go build && ./busuanzi-sync`


## 工作流程

1. 读取博客的 sitemap.xml 获取到 博客相关的所有页面
2. 顺序读取原版不蒜子的统计数据
3. 同步到 soxft/busuanzi 使用的 redis 中


## 注意

如果您使用 docker-compose 部署 soxft/busuanzi, 那么可能需要在 docker-compose.yml 中添加如下配置：

```
......
  redis:
    image: redis:alpine
    restart: always
    + ports:
    +  - "6379:6379"
......
```
这样才能让 busuanzi-sync 在宿主机上连接到 busuanzi 的 redis 服务。