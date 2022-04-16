# BLOGIN: GO-BLOG API
以博客为主题的 go 后端练手项目，适合初学快速上手基本的 go api 开发

## 使用的框架
|功能|框架|
|---|---|
|web框架|[gin](https://github.com/gin-gonic/gin)|
|数据库访问|[gorm](https://github.com/jinzhu/gorm)|
|缓存|[redigo](https://github.com/gomodule/redigo)|
|配置读取|[go-ini](https://github.com/go-ini/ini)|
|JWT|[jwt-go](https://github.com/dgrijalva/jwt-go)|
|参数验证|[go-playground/validator](https://github.com/go-playground/validator)|

## 项目结构

```
.
├─conf // 配置文件
├─logs // 存放生成出来的日志文件
├─middleware // 自定义的中间件
│  └─jwt // JWT 验证包
├─models // 访问数据库的包
├─pkg // 自己写的公用包
│  ├─app
│  ├─e
│  ├─file
│  ├─gredis
│  ├─logging
│  ├─setting
│  ├─upload
│  └─util
├─routers // 路由
│  └─api // api 层，目前 v1 版本则写在 v1 包中
│      └─v1
├─runtime // 运行时产生的文件，如用户上传的图片
│  └─upload
│      └─images
├─service // service 层，处理业务逻辑
│  ├─article_service
│  ├─cache_service
│  └─tag_service
└─sql // sql 文件，用于导入数据库    
```

## 功能
- 业务增删改查
    - Article 相关
    - Tag 相关
- 使用 Redis 做缓存，提高访问速度
- JWT 身份校验（未完善）
- 自动写日志文件

## 运行
1. 修改配置文件
2. 分别启动 mysql 和 redis
3. 项目文件夹下执行 `go run main.go`
4. 测试 api（[api 示例文档](https://documenter.getpostman.com/view/12524171/Uyr4JzQ7)）

## 感谢这些资料
- [跟煎鱼学 Go](https://eddycjy.com/go-categories/)
- [gorm hook使用中的问题及核心源码解读](https://cloud.tencent.com/developer/article/1830811)