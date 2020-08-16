package teadb

type BaseDAO struct {
	driver DriverInterface
}

func (this *BaseDAO) SetDriver(driver DriverInterface) {
	this.driver = driver
}
