package teaconfigs

import (
	"github.com/iwind/TeaGo/Tea"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var acmeLocalUserListFile = "acme_userlist.conf"

// 本地用户列表
type ACMELocalUserList struct {
	Users []*ACMELocalUser `yaml:"users" json:"users"`
}

// 取得共享的本地用户列表
func SharedACMELocalUserList() *ACMELocalUserList {
	data, err := ioutil.ReadFile(Tea.ConfigFile(acmeLocalUserListFile))
	if err != nil {
		return &ACMELocalUserList{}
	}

	userList := &ACMELocalUserList{}
	err = yaml.Unmarshal(data, userList)
	if err != nil {
		return &ACMELocalUserList{}
	}
	return userList
}

// 添加用户
func (this *ACMELocalUserList) AddUser(user *ACMELocalUser) {
	this.Users = append(this.Users, user)
}

// 删除用户
func (this *ACMELocalUserList) RemoveUser(userId string) {
	result := []*ACMELocalUser{}
	for _, user := range this.Users {
		if user.Id == userId {
			continue
		}
		result = append(result, user)
	}
	this.Users = result
}

// 查找用户
func (this *ACMELocalUserList) FindUser(userId string) *ACMELocalUser {
	for _, user := range this.Users {
		if user.Id == userId {
			return user
		}
	}
	return nil
}

// 保存
func (this *ACMELocalUserList) Save() error {
	data, err := yaml.Marshal(this)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(Tea.ConfigFile(acmeLocalUserListFile), data, 0666)
}
