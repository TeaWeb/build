# 内置变量
## 变量形式
在使用变量的地方，使用`${varName}`来表示变量，比如：
~~~
/${1}/hello${requestPath}?name=${arg.name}
~~~

## 匹配变量
匹配变量值的是有正则表达式的地方，使用匹配结果，通常为一个从0开始的数字，比如在重写规则中：
~~~
/(\w+)/(\w+)
~~~
那么
* `${0}` - 表示整体匹配的内容
* `${1}` - 表示第一个括号匹配的内容
* `${2}` - 表示第二个括号匹配的内容 

可以使用`(?i)`来设置不区分大小写：
~~~
(?i)/index.php
~~~

更多可用的正则表达式可以参考 [RE2 Syntax](https://github.com/google/re2/wiki/Syntax)。

### 命名变量
也可以给变量设置一个名称：
~~~
/(?P<myName>\w+)
~~~
然后就可以在待替换字符串中使用 `${myName}` ：
~~~
/hello/${myName}
~~~

## 请求相关变量
* `${teaVersion}` - TeaWeb版本
* `${remoteAddr}` - 客户端地址（IP），会依次根据X-Forwarded-For、X-Real-IP、RemoteAddr获取
* `${rawRemoteAddr}` - 客户端地址（IP），返回直接连接服务的客户端原始IP地址，从v0.1.3版本加入
* `${remotePort}` - 客户端端口
* `${remoteUser}` - 客户端用户名
* `${requestURI}` - 请求URI
* `${requestPath}` - 请求路径（不包括参数）
* `${requestLength}` - 请求内容长度
* `${requestMethod}` - 请求方法
* `${requestFilename}` - 请求文件路径
* `${scheme}` - 请求协议，`http`或`https`
* `${proto}` - 包含版本的HTTP请求协议，类似于`HTTP/1.0`
* `${timeISO8601}` - ISO 8601格式的时间，比如`2018-07-16T23:52:24.839+08:00`
* `${timeLocal}` - 本地时间，比如`17/Jul/2018:09:52:24 +0800`
* `${msec}` - 带有毫秒的时间，比如`1531756823.054`
* `${timestamp}` - unix时间戳，单位为秒
* `${host}` - 主机名
* `${serverName}` - 接收请求的服务器名
* `${serverPort}` - 接收请求的服务器端口
* `${referer}` - 请求来源URL
* `${userAgent}` - 客户端信息
* `${contentType}` - 请求头部的Content-Type
* `${cookies}` - 所有cookie组合字符串
* `${cookie.NAME}` - 单个cookie值
* `${args}` - 所有参数组合字符串
* `${arg.NAME}` - 单个参数值
* `${headers}` - 所有Header信息组合字符串
* `${header.NAME}` - 单个Header值

## 响应相关变量
* `${requestTime}` - 请求花费时间
* `${bytesSent}` - 发送的内容长度，包括Header（字节）
* `${bodyBytesSent}` - 发送的内容长度，不包括Header（字节）
* `${status}` - 状态码，比如`200`
* `${statusMessage}` - 状态消息，比如`200 OK`
* `${backend.id}` - 后端服务器ID，v0.0.9开始支持
* `${backend.code}` - 后端服务器代号，v0.0.9开始支持
* `${backend.address}` - 后端服务器地址，v0.0.9开始支持
* `${backend.scheme}` - 后端服务器协议，`http`或`https`，v0.0.9开始支持
