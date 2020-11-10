module github.com/TeaWeb/build

go 1.14

// 临时解决winio对父子进程支持错误的问题
replace github.com/Microsoft/go-winio v0.4.14 => github.com/bi-zone/go-winio v0.4.15

require (
	github.com/Azure/azure-sdk-for-go v45.1.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.3
	github.com/Azure/go-autorest/autorest/adal v0.9.2
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.1
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/Azure/go-autorest/autorest/validation v0.3.0 // indirect
	github.com/BurntSushi/toml v0.3.1
	github.com/JamesClonk/vultr v2.0.1+incompatible
	github.com/Microsoft/go-winio v0.4.14
	github.com/OpenDNS/vegadns2client v0.0.0-20180418235048-a3fa4a771d87
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/TeaWeb/agentinstaller v0.0.0-20200816121010-ed1b610d1130
	github.com/TeaWeb/plugin v0.0.0-20200816024143-17a5fe926d98
	github.com/TeaWeb/uaparser v0.0.0-20190526084055-a1c9449348d8
	github.com/akamai/AkamaiOPEN-edgegrid-golang v0.9.18
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.411
	github.com/aws/aws-sdk-go v1.29.15
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/cloudflare/cloudflare-go v0.13.0
	github.com/cpu/goacmedns v0.0.3
	github.com/dchest/captcha v0.0.0-20170622155422-6a29415a8364
	github.com/dchest/siphash v1.2.1
	github.com/decker502/dnspod-go v0.2.0
	github.com/dnsimple/dnsimple-go v0.63.0
	github.com/exoscale/egoscale v1.19.0
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-redis/redis/v8 v8.0.0-beta.7
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/gophercloud/gophercloud v0.12.0
	github.com/gorilla/websocket v1.4.2
	github.com/iij/doapi v0.0.0-20190504054126-0bbf12d6d7df
	github.com/iwind/TeaGo v0.0.0-20201110043415-859f4b3b98f3
	github.com/iwind/gofcgi v0.0.0-20181229122301-daea2786cb0d
	github.com/jlaffaye/ftp v0.0.0-20200812143550-39e3779af0db
	github.com/json-iterator/go v1.1.10
	github.com/labbsr0x/bindman-dns-webhook v1.0.2
	github.com/lib/pq v1.8.0
	github.com/linode/linodego v0.20.0
	github.com/mailru/easyjson v0.7.3
	github.com/miekg/dns v1.1.31
	github.com/namedotcom/go v0.0.0-20180403034216-08470befbe04
	github.com/nrdcg/auroradns v1.0.1
	github.com/nrdcg/goinwx v0.8.1
	github.com/oracle/oci-go-sdk v23.0.0+incompatible
	github.com/oschwald/geoip2-golang v1.4.0
	github.com/ovh/go-ovh v1.1.0
	github.com/pkg/sftp v1.11.0
	github.com/pquerna/ffjson v0.0.0-20190930134022-aa0246cd15f7
	github.com/rainycape/memcache v0.0.0-20150622160815-1031fa0ce2f2
	github.com/robertkrimen/otto v0.0.0-20191219234010-c382bd3c16ff
	github.com/sacloud/libsacloud v1.36.2
	github.com/shirou/gopsutil v2.20.7+incompatible
	github.com/stretchr/testify v1.6.1
	github.com/syndtr/goleveldb v1.0.0
	github.com/tatsushid/go-fastping v0.0.0-20160109021039-d7bb493dee3e
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	github.com/timewasted/linode v0.0.0-20160829202747-37e84520dcf7
	github.com/transip/gotransip v5.8.2+incompatible
	github.com/urfave/cli v1.22.4
	github.com/vmihailenco/msgpack v4.0.4+incompatible
	go.mongodb.org/mongo-driver v1.4.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20200803210538-64077c9b5642
	google.golang.org/api v0.30.0
	gopkg.in/ns1/ns1-go.v2 v2.4.2
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	gopkg.in/square/go-jose.v2 v2.5.1
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
)
