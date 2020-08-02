# 一些有用的脚本
* `teaweb` - TeaWeb控制服务
* `teaweb-agent` - Agent控制服务
* `man/teaweb.1` - Man Page

需要修改脚本里的`INSTALL_DIR`为实际的安装目录，然后拷贝到`/etc/init.d`目录。

## 添加服务
~~~bash
chkconfig --add teaweb
chkconfig --add teaweb-agent
~~~
