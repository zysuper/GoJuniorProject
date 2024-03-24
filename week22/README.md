## 作业：部署并压测 WebSocket

部署一个 WebSocket 服务器，利用 wrk 或者 k6 来压测。

而后，你需要：

* 调整 WebSocket 的 read buffer 和 write buffer 设置，超时设置。
* 调整请求和响应的大小。
* 调整并发数。

（可选）如果你有 Linux 设备，那么在 Linux 设备上调整 TCP 有关的参数
运行压测，体会一下这些参数对 WebSocket 性能的影响。

你只需要提交自己实验过程中测出来的数据，试着总结一下里面的规律。
