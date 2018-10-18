# 代理
参考：
* http://nginx.org/en/docs/http/ngx_http_upstream_module.html
* http://nginx.org/en/docs/http/ngx_http_proxy_module.html

# 文件
每个代理服务一个单独的文件，文件名以`.proxy.conf`结尾，放在`configs/`目录下：
~~~
$ROOT/
    configs/
~~~

# 单个代理文件配置
~~~yaml

listen: "127.0.0.1:80" # 多个代理可以共用同一个地址
name: [ "example.com", "*.example.com" ]
backends:  
  server1:
    host: ...
    weight: 5
    backup: true
    maxFails: 3
  server2:
    ...
~~~

# 内置变量
