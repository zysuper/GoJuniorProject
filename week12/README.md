# 第十一次作业

## 要求

在前面的例子里面，我直接在 Handler 层面上聚合了 gRPC 服务。

理论上来说，这不符合 DDD 的设计，但是好用。

按照 DDD 的设计来说的话，这边应该是要把 Interactive 的 gRPC 做成一个 Repository，而后在 ~~ArticleRepository~~  里面完成
Interactive 相关的组装。

换言之，将 Interactive 看做是 Article 的一个部分。

因此你的作业就是，用这种形态来集成 gRPC 的 Interactive。并且感受一下课堂风格和这种风格之间的差异。

### 实现

web Handler -> ArticleService -> ArticleRepository -> IntrRepository

### 快速跳转

1. [repository/IntrRepository.go](./webook/internal/repository/IntrRepository.go)
1. [repository/ArticleRepository.go](./webook/internal/repository/article.go#L27)
1. [service/article.go](./webook/internal/service/article.go#26)
1. [web/article.go](./webook/internal/web/article.go)
1. [wire.go](./webook/wire.go#47)


