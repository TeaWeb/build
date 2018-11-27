# 状态码分析
## 内置函数
* [resp.header()](#respheader)
* [resp.value()](#respvalue)

## resp.header()
获取Header值。

方法调用：
~~~javascript
resp.header("名称");
~~~

示例：
~~~javascript
var status = resp.header("My-Status-Code");
var contentType = resp.header("Content-Type");
~~~

## resp.value()
获取JSON数据中的某个字段。

方法调用：
~~~javascript
resp.value("字段");
~~~

示例：
~~~
var code = resp.value("code"); // => 10000
var age = resp.value("data.age"); // => 20

/**
// 示例数据
{
    "code": 10000,
    "data": {
        "name": "lu",
        "age": 20,
        "books": [ "Golang", "Python", "PHP" ],
        "isVip": true
    }
}
**/
~~~
