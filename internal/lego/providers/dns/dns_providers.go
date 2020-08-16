package dns

import (
	"fmt"

	"github.com/TeaWeb/build/internal/lego/challenge"
	"github.com/TeaWeb/build/internal/lego/challenge/dns01"
	"github.com/TeaWeb/build/internal/lego/providers/dns/acmedns"
	"github.com/TeaWeb/build/internal/lego/providers/dns/alidns"
	"github.com/TeaWeb/build/internal/lego/providers/dns/auroradns"
	"github.com/TeaWeb/build/internal/lego/providers/dns/azure"
	"github.com/TeaWeb/build/internal/lego/providers/dns/bindman"
	"github.com/TeaWeb/build/internal/lego/providers/dns/bluecat"
	"github.com/TeaWeb/build/internal/lego/providers/dns/cloudflare"
	"github.com/TeaWeb/build/internal/lego/providers/dns/cloudns"
	"github.com/TeaWeb/build/internal/lego/providers/dns/cloudxns"
	"github.com/TeaWeb/build/internal/lego/providers/dns/conoha"
	"github.com/TeaWeb/build/internal/lego/providers/dns/designate"
	"github.com/TeaWeb/build/internal/lego/providers/dns/digitalocean"
	"github.com/TeaWeb/build/internal/lego/providers/dns/dnsimple"
	"github.com/TeaWeb/build/internal/lego/providers/dns/dnsmadeeasy"
	"github.com/TeaWeb/build/internal/lego/providers/dns/dnspod"
	"github.com/TeaWeb/build/internal/lego/providers/dns/dode"
	"github.com/TeaWeb/build/internal/lego/providers/dns/dreamhost"
	"github.com/TeaWeb/build/internal/lego/providers/dns/duckdns"
	"github.com/TeaWeb/build/internal/lego/providers/dns/dyn"
	"github.com/TeaWeb/build/internal/lego/providers/dns/easydns"
	"github.com/TeaWeb/build/internal/lego/providers/dns/exec"
	"github.com/TeaWeb/build/internal/lego/providers/dns/exoscale"
	"github.com/TeaWeb/build/internal/lego/providers/dns/fastdns"
	"github.com/TeaWeb/build/internal/lego/providers/dns/gandi"
	"github.com/TeaWeb/build/internal/lego/providers/dns/gandiv5"
	"github.com/TeaWeb/build/internal/lego/providers/dns/gcloud"
	"github.com/TeaWeb/build/internal/lego/providers/dns/glesys"
	"github.com/TeaWeb/build/internal/lego/providers/dns/godaddy"
	"github.com/TeaWeb/build/internal/lego/providers/dns/hostingde"
	"github.com/TeaWeb/build/internal/lego/providers/dns/httpreq"
	"github.com/TeaWeb/build/internal/lego/providers/dns/iij"
	"github.com/TeaWeb/build/internal/lego/providers/dns/inwx"
	"github.com/TeaWeb/build/internal/lego/providers/dns/joker"
	"github.com/TeaWeb/build/internal/lego/providers/dns/lightsail"
	"github.com/TeaWeb/build/internal/lego/providers/dns/linode"
	"github.com/TeaWeb/build/internal/lego/providers/dns/linodev4"
	"github.com/TeaWeb/build/internal/lego/providers/dns/mydnsjp"
	"github.com/TeaWeb/build/internal/lego/providers/dns/namecheap"
	"github.com/TeaWeb/build/internal/lego/providers/dns/namedotcom"
	"github.com/TeaWeb/build/internal/lego/providers/dns/netcup"
	"github.com/TeaWeb/build/internal/lego/providers/dns/nifcloud"
	"github.com/TeaWeb/build/internal/lego/providers/dns/ns1"
	"github.com/TeaWeb/build/internal/lego/providers/dns/oraclecloud"
	"github.com/TeaWeb/build/internal/lego/providers/dns/otc"
	"github.com/TeaWeb/build/internal/lego/providers/dns/ovh"
	"github.com/TeaWeb/build/internal/lego/providers/dns/pdns"
	"github.com/TeaWeb/build/internal/lego/providers/dns/rackspace"
	"github.com/TeaWeb/build/internal/lego/providers/dns/rfc2136"
	"github.com/TeaWeb/build/internal/lego/providers/dns/route53"
	"github.com/TeaWeb/build/internal/lego/providers/dns/sakuracloud"
	"github.com/TeaWeb/build/internal/lego/providers/dns/selectel"
	"github.com/TeaWeb/build/internal/lego/providers/dns/stackpath"
	"github.com/TeaWeb/build/internal/lego/providers/dns/transip"
	"github.com/TeaWeb/build/internal/lego/providers/dns/vegadns"
	"github.com/TeaWeb/build/internal/lego/providers/dns/vscale"
	"github.com/TeaWeb/build/internal/lego/providers/dns/vultr"
	"github.com/TeaWeb/build/internal/lego/providers/dns/zoneee"
)

// NewDNSChallengeProviderByName Factory for DNS providers
func NewDNSChallengeProviderByName(name string) (challenge.Provider, error) {
	switch name {
	case "acme-dns":
		return acmedns.NewDNSProvider()
	case "alidns":
		return alidns.NewDNSProvider()
	case "azure":
		return azure.NewDNSProvider()
	case "auroradns":
		return auroradns.NewDNSProvider()
	case "bindman":
		return bindman.NewDNSProvider()
	case "bluecat":
		return bluecat.NewDNSProvider()
	case "cloudflare":
		return cloudflare.NewDNSProvider()
	case "cloudns":
		return cloudns.NewDNSProvider()
	case "cloudxns":
		return cloudxns.NewDNSProvider()
	case "conoha":
		return conoha.NewDNSProvider()
	case "designate":
		return designate.NewDNSProvider()
	case "digitalocean":
		return digitalocean.NewDNSProvider()
	case "dnsimple":
		return dnsimple.NewDNSProvider()
	case "dnsmadeeasy":
		return dnsmadeeasy.NewDNSProvider()
	case "dnspod":
		return dnspod.NewDNSProvider()
	case "dode":
		return dode.NewDNSProvider()
	case "dreamhost":
		return dreamhost.NewDNSProvider()
	case "duckdns":
		return duckdns.NewDNSProvider()
	case "dyn":
		return dyn.NewDNSProvider()
	case "fastdns":
		return fastdns.NewDNSProvider()
	case "easydns":
		return easydns.NewDNSProvider()
	case "exec":
		return exec.NewDNSProvider()
	case "exoscale":
		return exoscale.NewDNSProvider()
	case "gandi":
		return gandi.NewDNSProvider()
	case "gandiv5":
		return gandiv5.NewDNSProvider()
	case "glesys":
		return glesys.NewDNSProvider()
	case "gcloud":
		return gcloud.NewDNSProvider()
	case "godaddy":
		return godaddy.NewDNSProvider()
	case "hostingde":
		return hostingde.NewDNSProvider()
	case "httpreq":
		return httpreq.NewDNSProvider()
	case "iij":
		return iij.NewDNSProvider()
	case "inwx":
		return inwx.NewDNSProvider()
	case "joker":
		return joker.NewDNSProvider()
	case "lightsail":
		return lightsail.NewDNSProvider()
	case "linode":
		return linode.NewDNSProvider()
	case "linodev4":
		return linodev4.NewDNSProvider()
	case "manual":
		return dns01.NewDNSProviderManual()
	case "mydnsjp":
		return mydnsjp.NewDNSProvider()
	case "namecheap":
		return namecheap.NewDNSProvider()
	case "namedotcom":
		return namedotcom.NewDNSProvider()
	case "netcup":
		return netcup.NewDNSProvider()
	case "nifcloud":
		return nifcloud.NewDNSProvider()
	case "ns1":
		return ns1.NewDNSProvider()
	case "oraclecloud":
		return oraclecloud.NewDNSProvider()
	case "otc":
		return otc.NewDNSProvider()
	case "ovh":
		return ovh.NewDNSProvider()
	case "pdns":
		return pdns.NewDNSProvider()
	case "rackspace":
		return rackspace.NewDNSProvider()
	case "route53":
		return route53.NewDNSProvider()
	case "rfc2136":
		return rfc2136.NewDNSProvider()
	case "sakuracloud":
		return sakuracloud.NewDNSProvider()
	case "stackpath":
		return stackpath.NewDNSProvider()
	case "selectel":
		return selectel.NewDNSProvider()
	case "transip":
		return transip.NewDNSProvider()
	case "vegadns":
		return vegadns.NewDNSProvider()
	case "vultr":
		return vultr.NewDNSProvider()
	case "vscale":
		return vscale.NewDNSProvider()
	case "zoneee":
		return zoneee.NewDNSProvider()
	default:
		return nil, fmt.Errorf("unrecognised DNS provider: %s", name)
	}
}
