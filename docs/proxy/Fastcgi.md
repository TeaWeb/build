# Fastcgi
可以使用TeaWeb直接将请求分发到后端的Fastcgi。

## 前提
在使用Fastcgi分发之前，请保证"后端服务器"列表为空，否则TeaWeb会直接将请求分发到后端服务器。

## 步骤1 - 添加路径规则
比如我们的Fastcgi分发的是PHP文件（像php-fpm），通常扩展名都是 *.php*，那么可以在"路径规则"中添加一条规则：
1、点击"路径规则"下的"新路径规则"，输入相关信息：
![fastcgi1.png](fastcgi1.png)

2、点击底部的对号图标，保存。

## 步骤2 - 添加Fastcgi配置
3、在跳转后的页面"Fastcgi配置"下点击"+"加号图标，然后填入以下信息：
![fastcgi2.png](fastcgi2.png)
其中 *Fastcgi地址* 是Fastcgi的端口地址，如果你是使用unix socket启动，可以填入unix socket的绝对路径；*SCRIPT_FILENAME* 是接收请求的入口文件，通常是一个脚本，比如PHP文件；*DOCUMENT_ROOT* 是脚本所在本目录。在这里还可以添加更多的Fastcgi参数。

4、点击底部的"对号"图标，保存

5、根据顶部的提示，重启后即可生效。

## 设置首页文件
如果不想在首页输入*index.php*这样的路径，则可以在代理服务的"基本信息"中设置"首页文件"：
~~~
index.html index.php
~~~
多个首页文件使用空格隔开。

## 分发静态内容
如果网站有静态内容需要分发，可以在代理服务的"基本信息"中设置"文档根目录"。