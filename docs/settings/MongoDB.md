# MongoDB
在这里可以修改MongoDB的连接设置：
![mongodb.png](mongodb.png)

如果没有安装或者连接失败，则会提示连接失败信息：
![mongodb.png](mongodb2.png)

如上图所示，如果在项目目录下没有找到已安装的 *mongodb* ，则提示下载链接和自动安装链接。

## 安装和启动MongoDB
如果你还没有安装MongoDB，可以使用TeaWeb帮你安装（"设置">"MongoDB"界面的底部），也可以从 [https://www.mongodb.com/download-center/community](https://www.mongodb.com/download-center/community)下载并安装符合你的系统的MongoDB，目前支持Linux和MacOS。

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
