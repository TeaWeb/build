* 支持静态文件分发
* 修改server结构：
~~~~
ip address:
    server name1:
        server1, server2, ...
    server name2
        server3, server4, ...
    ...
~~~~
* backends支持https
* teaweb之间可以互联
* 设置自动启动项（任务）
* 支持websocket
* file session增加IP等验证，加强sid的安全性
* 添加配置后自动重启
* 支持rewrite
* widget支持moreUrl
* 代理服务显示连接状态
* 智能检测环境、检测各个Web应用配置、检测各个软件（MySQL、redis等）
* 实现一个安装程序
* 实现自动focus/esc/回车：
~~~
<div v-if="showEditing">
    <input type="text" data-tea-focus=""/>
</div>
~~~
* 添加代理域名，添加关闭按钮（"-"）/支持回车保存/添加取消编辑按钮
* SSL证书可以点击查看内容，或以内容方式添加
* 在代理管理中增加"日志" Tab，可以查看代理的操作后的日志
* 代理增加关闭、启动、重启等功能
* （多端操控：命令行、App）
* 日志设置有效期：比如1天，30天，当月等
* 日志设置条数：比如10000条/代理服务，总200000条
* MongoDB配置界面自动刷新状态
* 将`tt`变成编译工具：
~~~
# 编译
tt build
tt build [DIR]

# 启动
tt 
tt [DIR]

# 测试
tt test
tt test [DIR]
~~~
好处是可以自动从源文件中提取一些信息
* 添加|修改|查看location时，增加URL测试功能
* 配置导入导出
* 预置php python ror、各框架等配置，而且开放第三方
* location支持root
* server支持root
* 支持限流 限速 限时 限并发
* 每个代理、location、rewrite、fastcgi都增加统计
* 参考fiddler
* 加入crontab、task manager功能
* 每个后端可以单独看日志
* 支持websocket
* 重写规则增加测试功能：/proxy/rewrite/test
* 鼠标放到重写规则中`proxy://lb001`上，则常出现弹窗，显示代理的详情
* 删除代理的时候，提示有重写规则或者其他代理在引用
* 实现server的root、fastcgi、rewrite修改功能
* 实现location的root
* 实现server、location、rewrite的拷贝功能
* fastcgi支持PATHINFO
* 限流矩阵：
         总体控制 | 单IP | 单用户 | 单请求 | 日 | 天 | 周 | 月 | 年 | 时间段
QPS 
流量
请求数总量
...
* 每个代理都有单独的操作日志、错误日志、访问日志
* 简单配置不需要重启即可动态加载！！！！！
* 支持守护进程，来守护各个服务，比如redis，nginx之类的
* 增加fastcgi配置：server, location, rewrite
* 代理支持设置变量
* listen port ssl 每个listen单独有ssl配置/每个backend也有单独ssl配置
* include:配置可以导出，可以导入，也可以把导出的配置文件路径填进去
* request.callRoot()智能提示文件名大小写错误
* 可以禁止某个目录的访问
* 将所有配置分为：基本配置、更多配置
* 日志筛选增加：成功/失败/status code
* helper支持BeforeAction(action actions.Action)或BeforeAction(action *actions.Action)
* 实现nginx导入流程：你现在可以关掉nginx，并打开teaweb了
* 使用地图展示地区和省份
* 统计增加代理服务筛选
* 管理启动项/管理定时任务（自动读取crontab）
* 可以把其他图表放到桌面上
* widget可拖动排序
* widget可以配置，并保存
* widget、chart可以单独刷新
* CPU和内存使用情况，应增加Top 5应用
* app:
    * cmd line: 命令行
    * site:官网
    * docs:文档地址
    * operations:操作
    * ports: 端口
    * connections：连接数
    * alerts：报警信息
    * statistics：统计信息
    * cpu：CPU
    * memory：内存
    * disk usage: 文件系统占用空间
    * ...
* chart分概要信息、详细信息，可以展开收起    
* 优化内存使用/profile
* 在访问日志中 > 显示teaweb日志
* 管理界面可以设置只使用https:
* mongodb提示更友好
* fastcgi使用status role来ping server
* 代理换界面，换成竖向：
~~~
Proxy1 | Names | ... Chart |
Proxy2 | Names | ... Chart |
~~~
* !!!!!!!! 加速商业化 !!!!!!!!!!! 成立新公司 !!!!!!!!!!!!!!!!!!
  * 北京茶维科技有限公司
  * 商业版使用插件实现
  * 细化目标人群：测试人员、运维人员、开发人员，每个人群有针对性
* 访问区域报警：如果非允许的地区，比如山东，则会在顶部提示报警信息，或者阻挡
* 实现SCGI: http://python.ca/scgi/protocol.txt
* 实现CGI: http://python.ca/scgi/protocol.txt
* 在日志详情中增加region信息
* Parse和Stat太慢
* 优化添加、管理代理流程
* *成功的关键：可视化+统计+稳定*
* 升级提醒和在线升级
* 在日志界面中添加IP到禁止IP中
* 增加团队管理（Team）
* !!!!!!!!!每次升级必须有50%的功能放在代理、日志改进上!!!!!!!!!!!!!!!!!!
* !!!!!!!!!多用美美的图表，比如散点图！！！！！！！！！！！！！！
* 分解日志，让日志展示每项都很漂亮
* 各种排行换成散点图、地图（圆点有大小）
* 日志的TAB仿照chrome inspect
* 日志可以选择访问列
* 日志可以设置转发
* Dark模式theme
* 参考Kibana：https://www.oschina.net/news/100782/kibana-6-4-2-released
* 地区列表加上logo
* 默认提供Web访问，启动后，直接可以访问 www/ 目录下的文件
* 日志中记录文档的标题（Title）
* 日志筛选框显示预览
* 智能识别vuejs, angular.js, bootstrap.js, ...，标记图标、简介、官方文档。。。
* 智能识别PHP文件（根据扩展名、Header）、Ruby、Python、Java。。。。。
* 智能识别各种框架生成的网页
* 实现日志的请求数据截取、预览功能、响应内容
* 智能识别微信（MicroMessenger）等终端
* 日志终端信息中大量使用图标
* 可以预览更多资源文件
* 日志中可以显示JS错误、404错误以及其他异常，单独列为一个TAB方便操作
* 一键更新IP库、浏览器库
* 日志可以按照资源分类筛选：图片、CSS、Javascript
* 识别127.x.x.x, 192.168...., 10.10.xx.xx 为本机
* 自行实现终端库(os, browser, ip, ...)
* 日志导出到其他API、工具中
* 具体server、proxy、location、rewrite的具体校验错误，在代理界面中要显示出来
* 日志界面中可以添加一组要浏览的日志文件，比如/var/log/mysql/slow_log
* 代理服务器列表可以拖动排序
* 代理服务自动检测后端服务器的可用性，并在界面上显示
* 代理界面展示代理是否启动成功、启动的错误提示
* 自动显示代理后端服务的服务器信息（nginx、tomcat、版本之类的）
* 加速插件化开发，开发几个有价值的插件：
   * HTTP请求过滤：golang、python、PHP、Javascript、Ruby、Lua
   * 敏感词过滤
   * 监控
* QPS突增报警，阈值   
* plugins自动检测，提示升级
* 插件自动加载
* 插件页面加入如何做自己的插件
* 增加网页来源统计
* 自动记录TeaWeb日志
* Fastcgi错误
* 代理设置重新组织界面展示形式
* 实现一个TeaWeb安装工具，直接执行就可以安装：
  ~~~
  ./teaweb-install
  ............................... 30%
  ~~~
* 增加权限控制：https://segmentfault.com/a/1190000016698217
* 大文件上传60s超时
* fastcgi可以直接write给特定的writer，不需要使用[]byte
* 当前正在上传的请求列表
* 最大上传限制  
* 插件：系统错误日志
* 日志分代理Tab
* 打造一款激动人心的产品
* POST数据抓取、文件抓取
* 代理设置增加搜索功能
* Headers增加排序功能
* rewrite target支持完整的HTTP/HTTPS URL
* widget chart可以浮动到各个页面
* value变更时加一个放大特效
* Rewrite遇到重复则停止
* 自动检测本机网络服务加入到代理中、监控中
* 变量文档，在可以使用变量的地方添加帮助文档
* TeaWeb限制登录的IP地址