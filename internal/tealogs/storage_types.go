package tealogs

// 存储引擎类型
type StorageType = string

const (
	StorageTypeFile    StorageType = "file"
	StorageTypeES      StorageType = "es"
	StorageTypeMySQL   StorageType = "mysql"
	StorageTypeTCP     StorageType = "tcp"
	StorageTypeSyslog  StorageType = "syslog"
	StorageTypeCommand StorageType = "command"
)
