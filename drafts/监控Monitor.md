# 监控

## 界面
~~~
监控界面 | 添加+

|cpu|memory|load|disk|io|...|

|-----------------|
| Redis       |x| |
|                 |  ....
| Icon/Chart/State|
|-----------------|

~~~

## 智能检测
~~~
teaservices/
  redis
  mysql
  ...
~~~


## 监控项
~~~
watching.tasks:
    id: "123456",
    name: "Redis001",
    app: "redis",
    info: {
        ...
    },
    api: {
        http:...
        https:...
        socket:...
        port:...
        script:...
        file:...
        pipe:
    },
    timeout:...,
    interval: "1h"
~~~