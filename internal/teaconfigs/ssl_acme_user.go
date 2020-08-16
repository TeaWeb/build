package teaconfigs

import (
	"crypto"
	"github.com/TeaWeb/build/internal/lego/registration"
)

// ACME用户账号定义
type ACMEUser struct {
	Email        string
	Registration *registration.Resource
	Key          crypto.PrivateKey
}

func (this *ACMEUser) GetEmail() string {
	return this.Email
}
func (this *ACMEUser) GetRegistration() *registration.Resource {
	return this.Registration
}
func (this *ACMEUser) GetPrivateKey() crypto.PrivateKey {
	return this.Key
}
