~~~
admin.conf - 管理员配置
server.conf - 管理界面访问地址配置

db.conf - 数据库配置
mongo.conf - MongoDB配置
mysql.conf - MySQL配置
postgres.conf - PostgreSQL配置

notice.conf - 通知全局配置

cache.conf - 缓存全局配置
cache.policy.*.conf - 缓存策略

waflist.conf - WAF全局配置
waf.*.conf - WAF策略

serverlist.conf - 代理服务全局配置
server.*.conf - 单个代理服务配置

ssl.* - HTTPS证书相关文件

node.conf - 当前节点配置
cluster.sum - 主节点的sum信息
node.sum - 当前节点的sum信息
cluster.sum - 集群中心的sum信息

agents/* - Agent相关配置
~~~

`*.conf`配置都使用 [YAML](http://yaml.org/) 格式，通常可以通过管理界面进行修改。