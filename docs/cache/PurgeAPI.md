# 清除缓存指令
从v0.1.1开始，可以在请求Header中加入`Tea-Cache-Purge: 1`指令清除单个URL的缓存。

比如`/book/index.html`已经被缓存，则可以再次请求此地址，同时加入指令来清除缓存：
~~~http
GET /book/index.html HTTP/1.1
...
Tea-Cache-Purge: 1
Tea-Key: z8O4MuXixbKH6aiVyZigYTxxovRblR3u
...
~~~

其中`Tea-Cache-Purge`值固定为`1`，`Tea-Key`为TeaWeb登录用户的密钥（可以在"设置" -- "登录设置"中查看）。