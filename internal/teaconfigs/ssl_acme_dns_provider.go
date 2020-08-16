package teaconfigs

// 自定义ACME DNS解析
type ACMEDNSProvider struct {
	apiAuthToken string
}

func NewACMEDNSProvider(apiAuthToken string) (*ACMEDNSProvider) {
	return &ACMEDNSProvider{
		apiAuthToken: apiAuthToken,
	}
}

func (this *ACMEDNSProvider) Present(domain, token, keyAuth string) error {
	return nil
}

func (this *ACMEDNSProvider) CleanUp(domain, token, keyAuth string) error {
	return nil
}
