# 参数
参数是指TeaWAF提供一组可供检查的对象，参数值从HTTP请求和响应中获取。

在WAF中可以检查的参数列表如下：

## 客户端地址（IP）
* 前缀：`${remoteAddr}`
* 描述：试图通过分析X-Real-IP等Header获取的客户端地址，比如192.168.1.100
* 是否有子参数：NO

## 客户端源地址（IP）
* 前缀：`${rawRemoteAddr}`
* 描述：直接连接的客户端地址，比如192.168.1.100
* 是否有子参数：NO

## 客户端端口
* 前缀：`${remotePort}`
* 描述：直接连接的客户端地址端口
* 是否有子参数：NO

## 客户端用户名
* 前缀：`${remoteUser}`
* 描述：通过BasicAuth登录的客户端用户名
* 是否有子参数：NO

## 请求URI
* 前缀：`${requestURI}`
* 描述：包含URL参数的请求URI，比如/hello/world?lang=go
* 是否有子参数：NO

## 请求路径
* 前缀：`${requestPath}`
* 描述：不包含URL参数的请求路径，比如/hello/world
* 是否有子参数：NO

## 请求内容长度
* 前缀：`${requestLength}`
* 描述：请求Header中的Content-Length
* 是否有子参数：NO

## 请求体内容
* 前缀：`${requestBody}`
* 描述：通常在POST或者PUT等操作时会附带请求体，最大限制32M
* 是否有子参数：NO

## 请求URI和请求体组合
* 前缀：`${requestAll}`
* 描述：${requestURI}和${requestBody}组合
* 是否有子参数：NO

## 请求表单参数
* 前缀：`${requestForm}`
* 描述：获取POST或者其他方法发送的表单参数，最大请求体限制32M
* 是否有子参数：YES

## 上传文件
* 前缀：`${requestUpload}`
* 描述：获取POST上传的文件信息，最大请求体限制32M
* 是否有子参数：YES
* 可选子参数
   * `最小文件尺寸`：值为 `minSize`
   * `最大文件尺寸`：值为 `maxSize`
   * `扩展名(如.txt)`：值为 `ext`
   * `原始文件名`：值为 `name`
   * `表单字段名`：值为 `field`

## 请求JSON参数
* 前缀：`${requestJSON}`
* 描述：获取POST或者其他方法发送的JSON，最大请求体限制32M，使用点（.）符号表示多级数据
* 是否有子参数：YES

## 请求方法
* 前缀：`${requestMethod}`
* 描述：比如GET、POST
* 是否有子参数：NO

## 请求协议
* 前缀：`${scheme}`
* 描述：比如http或https
* 是否有子参数：NO

## HTTP协议版本
* 前缀：`${proto}`
* 描述：比如HTTP/1.1
* 是否有子参数：NO

## 主机名
* 前缀：`${host}`
* 描述：比如teaos.cn
* 是否有子参数：NO

## 请求来源URL
* 前缀：`${referer}`
* 描述：请求Header中的Referer值
* 是否有子参数：NO

## 客户端信息
* 前缀：`${userAgent}`
* 描述：比如Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103
* 是否有子参数：NO

## 内容类型
* 前缀：`${contentType}`
* 描述：请求Header的Content-Type
* 是否有子参数：NO

## 所有cookie组合字符串
* 前缀：`${cookies}`
* 描述：比如sid=IxZVPFhE&city=beijing&uid=18237
* 是否有子参数：NO

## 单个cookie值
* 前缀：`${cookie}`
* 描述：单个cookie值
* 是否有子参数：YES

## 所有URL参数组合
* 前缀：`${args}`
* 描述：比如name=lu&age=20
* 是否有子参数：NO

## 单个URL参数值
* 前缀：`${arg}`
* 描述：单个URL参数值
* 是否有子参数：YES

## 所有Header信息
* 前缀：`${headers}`
* 描述：使用
隔开的Header信息字符串
* 是否有子参数：NO

## 单个Header值
* 前缀：`${header}`
* 描述：单个Header值
* 是否有子参数：YES

## CC统计
* 前缀：`${cc}`
* 描述：统计某段时间段内的请求信息
* 是否有子参数：YES
* 可选子参数
   * `请求数`：值为 `requests`

## 响应状态码
* 前缀：`${status}`
* 描述：响应状态码，比如200、404、500
* 是否有子参数：NO

## 响应Header
* 前缀：`${responseHeader}`
* 描述：响应Header值
* 是否有子参数：YES

## 响应内容
* 前缀：`${responseBody}`
* 描述：响应内容字符串
* 是否有子参数：NO

## 响应内容长度
* 前缀：`${bytesSent}`
* 描述：响应内容长度，通过响应的Header Content-Length获取
* 是否有子参数：NO

## TeaWeb版本
* 前缀：`${teaVersion}`
* 描述：比如0.1.3
* 是否有子参数：NO
