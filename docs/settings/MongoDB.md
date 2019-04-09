# MongoDB
**注意：目前TeaWeb只支持MongoDB v3.0以上版本，如果你正在使用老的版本，请更换。**

在这里可以修改MongoDB的连接设置：
![mongodb.png](mongodb.png)

如果没有安装或者连接失败，则会提示连接失败信息：
![mongodb.png](mongodb2.png)

如上图所示，如果在项目目录下没有找到已安装的 *mongodb* ，则提示下载链接和自动安装链接。

## 安装和启动MongoDB
如果还没有安装，可以使用TeaWeb帮你安装（"设置">"MongoDB"界面的底部），目前支持Linux和Darwin（Mac OS X）。也可以从 [https://www.mongodb.com/download-center/community](https://www.mongodb.com/download-center/community)下载并安装符合你的系统的MongoDB，或者从TeaOS镜像下载地址中下载：
* [Linux版本](http://dl.teaos.cn/mongodb-linux-x86_64-4.0.3.tgz)
* [Darwin版本](http://dl.teaos.cn/mongodb-osx-ssl-x86_64-4.0.3.tgz)
* [Windows版本](http://dl.teaos.cn/mongodb-win32-x86_64-2008plus-ssl-4.0.8-signed.msi)

在Linux和MacOS上，解压MongoDB安装包后，建议的启动命令为：
~~~bash
cd MongoDB安装目录
bin/mongod --dbpath=./data/ --fork --logpath=./data/fork.log
~~~

启动后，试着用 `ps` 命令查看MongoDB是否已启动：
~~~bash
ps ax|grep mongo
~~~
命令执行结果应该类似于：
~~~
[root@localhost ~]# ps ax|grep mongo
21040 ?        Sl   632:19bin/mongod --dbpath=./data/ --fork --logpath=./data/fork.log
~~~
