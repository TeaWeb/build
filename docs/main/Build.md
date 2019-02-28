# 从源码启动或编译
## 安装Golang
如果你还没有安装Golang运行环境，请先安装，Go的版本需要在 v1.10.0以上，国内可以从 [GO语言中文网](https://studygolang.com/dl) 下载。

## 从源码运行
1. 从 `https://github.com/TeaWeb/build` 中下载项目源码，放到本地磁盘上；
2. 在开发工具中设置全局变量`GOPATH`为项目目录路径；
3. `cd` 到 *src/main* 目录
4. 执行 `init.sh` 初始化项目，如果下载中出现网络错误，可以尝试多次运行此脚本；
5. 执行 `run.sh` 启动项目；
6. 在浏览器中访问 `http://127.0.0.1:7777` 。

## 从源码编译
1. 从 *https://github.com/TeaWeb/build* 中下载项目源码，放到本地磁盘上；
2. 在开发工具中设置全局变量`GOPATH`为项目目录路径；
3. `cd` 到 *src/main* 目录
4. 执行 `init.sh` 初始化项目，如果下载中出现网络错误，可以尝试多次运行此脚本；
5. 运行 `build-[系统版本].sh` 构建可执行文件；
6. 构建后的文件在 `项目根目录/dist/` 目录下。

## 使用Git下载源码
如果你想使用Git下载源码，可以使用下面命令： 
~~~bash
git clone https://github.com/TeaWeb/build.git
~~~
然后再运行`init.sh`：
~~~bash
./init.sh
~~~

## GoLand
如果你正在使用GoLand开发工具，则可以在下面的界面中设置GOPATH:
![goland.png](goland.png)
