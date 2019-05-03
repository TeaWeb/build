# 常用命令
## 打印帮助信息
~~~bash
./bin/teaweb -h
~~~

## 打印版本信息
~~~bash
./bin/teaweb -v
~~~

## 启动服务
启动服务并在后台运行：
~~~bash
./bin/teaweb start
~~~

如果要在前端启动服务，并阻塞当前进程，可以使用：
~~~bash
./bin/teaweb
~~~

## 停止服务
~~~bash
./bin/teaweb stop
~~~

## 重启服务
~~~bash
./bin/teaweb restart
~~~

## 重新加载代理配置
v0.1以后支持
~~~bash
./bin/teaweb reload
~~~

## 重置服务状态
v0.0.8以后支持
~~~bash
./bin/teaweb reset
~~~

## 查看服务状态
v0.1以后支持
~~~bash
./bin/teaweb status
~~~
