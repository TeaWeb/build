# 代理

# 代理功能
* 上报监控信息
* 日志收集

# 目录结构
~~~
TeaAgent/
  agent
  configs/
  logs/
~~~

# 启动Agent
~~~
./agent start
~~~

# 配置
~~~
master: "192.168.1.101:1234"

log.mydomain.com: ...
log.mysql.slowlog: ...
...
~~~

# 