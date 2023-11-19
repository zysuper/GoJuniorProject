# 第六周作业 

在 Web 的 Handler 部分，有很多 if-else 分支，基本上都是在判定 err !=nil。如下图，每一个 if 里面都要打印日志。

<img src="https://static001.infoq.cn/resource/image/8d/28/8d04b1dea1f6a5d9df2aaaf6905e2c28.png" alt="img" style="zoom:50%;" />

现在要求你优化这些打印日志的逻辑，避免每一处 err !=nil 的时候，都得手动打一个日志。



