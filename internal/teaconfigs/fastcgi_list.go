package teaconfigs

import "math/rand"

// FastcgiList接口定义
type FastcgiListInterface interface {
	// 校验
	ValidateFastcgi() error

	// 取得所有的Fastcgi
	AllFastcgi() []*FastcgiConfig

	// 根据ID查找Fastcgi
	FindFastcgi(fastcgiId string) *FastcgiConfig

	// 添加Fastcgi
	AddFastcgi(fastcgi *FastcgiConfig)

	// 删除Fastcgi
	RemoveFastcgi(fastcgiId string)

	// 查找下一个可用的Fastcgi
	NextFastcgi() *FastcgiConfig
}

// FastcgiList定义
type FastcgiList struct {
	Fastcgi []*FastcgiConfig `yaml:"fastcgi" json:"fastcgi"`
}

// 校验
func (this *FastcgiList) ValidateFastcgi() error {
	for _, f := range this.Fastcgi {
		err := f.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

// 取得所有的Fastcgi
func (this *FastcgiList) AllFastcgi() []*FastcgiConfig {
	if this.Fastcgi == nil {
		return []*FastcgiConfig{}
	}
	return this.Fastcgi
}

// 根据ID查找Fastcgi
func (this *FastcgiList) FindFastcgi(fastcgiId string) *FastcgiConfig {
	for _, f := range this.Fastcgi {
		if f.Id == fastcgiId {
			f.Validate()
			return f
		}
	}
	return nil
}

// 添加Fastcgi
func (this *FastcgiList) AddFastcgi(fastcgi *FastcgiConfig) {
	this.Fastcgi = append(this.Fastcgi, fastcgi)
}

// 删除Fastcgi
func (this *FastcgiList) RemoveFastcgi(fastcgiId string) {
	result := []*FastcgiConfig{}
	for _, f := range this.Fastcgi {
		if f.Id == fastcgiId {
			continue
		}
		result = append(result, f)
	}
	this.Fastcgi = result
}

// 查找下一个可用的Fastcgi
func (this *FastcgiList) NextFastcgi() *FastcgiConfig {
	result := []*FastcgiConfig{}
	for _, f := range this.Fastcgi {
		if !f.On {
			continue
		}
		result = append(result, f)
	}
	count := len(result)
	if count == 0 {
		return nil
	}
	return result[rand.Int()%count]
}
