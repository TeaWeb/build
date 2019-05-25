REM initialize project
set GOPATH=%cd%\..\..\

REM mkdir
if not exist %GOPATH%\src\main\web\tmp mkdir %GOPATH%\src\main\web\tmp

REM go_get function
REM TeaGO
call go_get "github.com/iwind/TeaGo"
call go_get "github.com/pquerna/ffjson"

REM fastcgi
call go_get "github.com/iwind/gofcgi"

REM geo ip
call go_get "github.com/oschwald/maxminddb-golang"
call go_get "github.com/oschwald/geoip2-golang"

REM system cpu, memory, disk ...
call go_get "github.com/shirou/gopsutil"
call go_get "github.com/shirou/w32"
call go_get "github.com/StackExchange/wmi"
call go_get "github.com/go-ole/go-ole"

REM javascript
call go_get "github.com/robertkrimen/otto"

REM msg pack
call go_get "github.com/vmihailenco/msgpack"

REM redis
call go_get "github.com/go-redis/redis"

REM markdown
call go_get "github.com/russross/blackfriday"

REM fsnotify
call go_get "github.com/fsnotify/fsnotify"

REM websocket
call go_get "github.com/gorilla/websocket"

REM leveldb
call go_get "github.com/syndtr/goleveldb/leveldb"

REM go winio
call go_get "github.com/Microsoft/go-winio"

REM ping
call go_get "github.com/tatsushid/go-fastping"

REM ssh
call go_get "github.com/pkg/sftp"

REM mysql
call go_get "github.com/go-sql-driver/mysql"

REM pqsql
call go_get "github.com/lib/pq"

REM siphash
github.com/dchest/siphash

REM aliyun
call go_get "github.com/aliyun/alibaba-cloud-sdk-go"

REM TeaWeb
call go_get "github.com/TeaWeb/code"
echo "   [TeaWeb]Don't worry, you can ignore 'no Go files' warning in TeaWeb code"

call go_get "github.com/TeaWeb/plugin"
call go_get "github.com/TeaWeb/uaparser"
call go_get "github.com/TeaWeb/agent"
call go_get "github.com/TeaWeb/agentinstaller"

echo "[done]"