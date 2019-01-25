# FAQ
## TeaWeb是干吗用的？
`TeaWeb` 是一个Web反向代理的服务工具，类似于Nginx、LVS之类，只不过 `TeaWeb` 试图提供一个可视化的界面，让用户操作特别简单，同时也自动实现日志、监控、统计等功能。

## 我能使用TeaWeb代替nginx吗？
`nginx`是一个非常优秀的Web Server，如果你在大规模地在用，不建议轻易更换。如果正在小规模使用，`TeaWeb`也提供了`nginx`具有的基础Web代理功能，既可以分发静态文件，可以分发Fastcgi请求，也实现了分发到后端服务器，所以假如你没有特殊额外的需求，完全可以使用TeaWeb代替nginx。

## TeaWeb怎样配置与php-fpm配合支持PHP呢？
可以在路径规则中使用Fastcgi配置，[请在这里查看相关文档](../proxy/Fastcgi.md)。

## 系统提示server selection timeout
出现这种提示的时候，说明MongoDB连接失败，请检查MongoDb连接，可以在界面右上角"设置" > MongoDB中查看MongoDB错误信息。

