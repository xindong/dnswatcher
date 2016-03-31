监控 dns 服务器的稳定性

### 工作方式

dnswatcher 每隔1秒钟向指定(-watch)的dns服务器查询指定(-domain)的域名记录。
如果超时或没有记录返回会记录失败数(failed)量。
如果成功，则记为一个成功(passed)的查询。
每隔1分钟会打印成功和失败的总数。

### 安装

`go install github.com/xindong/dnswatcher`

### 用法

```
-domain
    domains that need to be watched (default "google.com,baidu.com")
-help
    usages
-timeout
    set query timeout in seconds (default 5)
-watch
    nsserver that need to be watched (default "8.8.8.8,233.5.5.5")
```

例如 `dnswatcher -watch=8.8.8.8.233.5.5.5 -domain=google.com,baidu.com -timeout=10` 代表
监控向 8.8.8.8 和 233.5.5.5 查询 google.com baidu.com 的响应情况，超时时间设定为10秒
