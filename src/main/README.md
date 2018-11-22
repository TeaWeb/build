## 脚本说明
~~~
./init.sh - 初始化项目，在自己开发或者构建之前需要先运行此文件，以下载依赖程序包

./build-all.sh - 构建所有平台上可以运行的项目，会产生多个压缩文件
./build-darwin.sh - 构建darwin平台上可以运行的项目
./build-linux-32.sh - 构建linux 32位平台上可以运行的项目
./build-linux-64.sh - 构建linux 64位平台上可以运行的项目
./build-windows-32.sh - 构建windows 32位平台上可以运行的项目
./build-windows-64.sh - 构建windows 64位平台上可以运行的项目

./run.sh - 直接从源码运行服务

./utils.sh - 一些公用的shell函数
~~~

### 补充说明
构建后的压缩文件会放在`$项目目录/dist/`目录下。
