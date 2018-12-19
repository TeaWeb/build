# FAQ
## 我能使用TeaWeb代替nginx吗？
`nginx`是一个非常优秀的Web Server，如果你在大规模地在用，不建议轻易更换。如果正在小规模使用，`TeaWeb`也提供了`nginx`具有的基础Web代理功能，既可以分发静态文件，可以分发Fastcgi请求，也实现了分发到后端服务器，所以假如你没有特殊额外的需求，完全可以使用TeaWeb代替nginx。

## TeaWeb怎样配置与php-fpm配合支持PHP呢？
可以在路径规则中使用Fastcgi配置，[请在这里查看相关文档](../proxy/Fastcgi.md)。

