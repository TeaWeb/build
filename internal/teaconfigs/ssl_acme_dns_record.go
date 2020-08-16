package teaconfigs

// ACME DNS记录
type ACMEDNSRecord struct {
	FQDN  string `yaml:"fqdn" json:"fqdn"`
	Value string `yaml:"value" json:"value"`
}
