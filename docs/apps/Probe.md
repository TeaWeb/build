# 自定义App探针

## 脚本API
* ProcessProbe - 进程探针
  * `id` string - 探针ID，每个探针使用此ID作为区分，通常为系统自动生成，不要轻易修改
  * `author` string - 探针作者
  * `name` string - App名称
  * `site` string - App官方网站
  * `docSite` string - 官方文档网址
  * `developer` string - App开发者公司、团队或者个人名称
  * `commandName` string - App启动的命令名称，通常是启动服务的命令的最后一段，比如对于 `/usr/bin/svnserve -d /home/svn`，可以填入 `svnserve`
  * `commandPatterns` []string - 如果服务启动了多个进程，或者有多个服务有相同的命令文件名的时候，可能会有多个匹配结果，可以使用匹配规则来只匹配我们想要的。匹配规则中支持正则表达式。
  * `commandVersion` string - 获取版本信息的命令，在其中可以使用 `${commandFile}` 表示命令行文件路径，`${commandDir}` 命令行文件所在目录，常见的有：
    * `${commandFile} --version` 
    * `${comamndFile} -version`
    * `${comamndFile} -v`
  * `onProcess()` function - 筛选进程时回调：
      ~~~javascript
      probe.onProcess(function (p) {
          return true;
      });
      ~~~
      如果返回`true`表示是匹配正确的进程，如果是`false`，此进程会被忽略；进程对象`p`可以使用的属性为：
      * `name` string - 名称
      * `pid` int - 进程ID
      * `ppid` int - 父进程ID
      * `cwd` string - 启动时所在目录
      * `user` string - 启动进程的用户名
      * `uid` string - 启动进程的用户ID
      * `gid` string - 启动进程的用户所在组ID
      * `cmdline` string - 完整的命令行
      * `file` string - 命令行文件路径
      * `dir` string - 命令行文件所在目录
      * `isRunning` bool - 是否正在运行 
  * `onParseVersion()` function - 分析版本时回调，如果设置了`commandVersion`，并且运行后获得了App的版本号，会调用此函数来进行二次分析：
      ~~~javascript
      probe.onParseVersion(function (v) {
          return v;
      });
      ~~~~