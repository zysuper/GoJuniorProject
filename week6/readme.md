# 第六周作业 

## 作业要求

在 Web 的 Handler 部分，有很多 if-else 分支，基本上都是在判定 err !=nil。如下图，每一个 if 里面都要打印日志。

<img src="https://static001.infoq.cn/resource/image/8d/28/8d04b1dea1f6a5d9df2aaaf6905e2c28.png" alt="img" style="zoom:50%;" />

现在要求你优化这些打印日志的逻辑，避免每一处 err !=nil 的时候，都得手动打一个日志。

## 实现

使用 gin 的 `ERROR MANAGEMENT` 机制，结合自定义中间件实现。

以 user handler 的 `SignUp` 方法为例，异常分支，设置 `cox.Error`

![image-20231119225119910](/Users/zysuper/Library/Application Support/typora-user-images/image-20231119225119910.png)

然后，在 log 中间件，对 Error 进行统一打印：

![image-20231119225254930](/Users/zysuper/Library/Application Support/typora-user-images/image-20231119225254930.png)

最终日志效果如下：

![image-20231119225331017](/Users/zysuper/Library/Application Support/typora-user-images/image-20231119225331017.png)

### 代码快捷跳转

[user.go](./webook/internal/web/user.go)

[log.go](./webook/internal/web/middleware/log.go)

