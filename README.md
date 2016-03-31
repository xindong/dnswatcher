### 监控 dns 服务器的稳定性

例如 `go run main.go -watch=8.8.8.8.233.5.5.5 -domain=google.com,baidu.com -timeout=10` 监控向 8.8.8.8 和 233.5.5.5 查询 google.com baidu.com 的响应情况，超时时间设定为10秒
