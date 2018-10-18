# 参考 http://nginx.org/en/docs/http/ngx_http_upstream_module.html


*configs/example.server.conf*
~~~yaml
# listen: "127.0.0.1:80" # 监听某个地址


## 以下参考 http://nginx.org/en/docs/http/ngx_http_core_module.html#server
name: [ "example.com", "*.example.com" ] # 域名
serverTokens: "on|off" # 是否在响应头部显示Server信息
types: # 自定义mime type
    html: text/html
    css: text/css
    jpg: image/jpeg
defaultType: application/octet-stream # 默认的mime type     

clientBodyTimeout: 60s # 请求体读取超时时间
clientMaxBodySize: 10M # 请求体最大尺寸

errorPage:
    404 /404.html
    500 /50x.html
~~~
