# 安装
1. 可以从 [https://github.com/TeaWeb/build/releases](https://github.com/TeaWeb/build/releases) 下载对应的Release版本，目前支持MacOS(darwin)、Linux、Windows三种版本；
2. 使用unzip解压到某个目录，比如 teaweb-v0.0.1/；
3. 然后执行：
    ~~~shell
    cd teaweb-v0.0.1/
    bin/teaweb start
    ~~~
4. 如果没有出现错误的话，可以在浏览器中访问：
    ~~~
    http://127.0.0.1:7777
    ~~~
    其中`127.0.0.1`可能需要换成你服务器的IP，而且我们默认使用了`7777`端口（可以在`configs/server.conf`中修改），如果访问遇到了问题，请检查防火墙设置；
5. 使用用户名`admin`和密码`123456`登录，可以在`configs/admin.conf`中修改这些信息，也可以在设置界面中修改。

## MongoDB
TeaWeb需要使用MongoDB来记录日志和其他数据，如果已经安装，可以在"设置">"MongoDB"中修改MongoDB的连接参数。如果还没有安装，可以使用TeaWeb帮你安装（"设置">"MongoDB"界面的底部），也可以从 [https://www.mongodb.com/download-center](https://www.mongodb.com/download-center)下载并安装符合你的系统的MongoDB。

## CentOS7
在CentOS7上，如果你需要使用`7777`端口，可能要在firewall中注册一个规则：
~~~
firewall-cmd --zone=public --add-port=7777/tcp --permanent
firewall-cmd —reload
~~~

要使用插件服务，请确保`ps`、`pgrep`和`lsof`命令可用，如果没有安装对应的命令，可以使用以下命令安装：
~~~
yum install procps
yum install lsof
~~~

## Windows
Windows版本的目录下自带有*start.bat*，请解压后，直接双击运行*start.bat*即可。