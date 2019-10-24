Tea.context(function () {
	this.addresses = [];
	if (this.server != null && this.server.https != null && this.server.https.listen != null) {
		this.addresses = this.server.https.listen;
	}

	/**
	 * 证书
	 */
	this.certTab = "shared";
	if (this.sharedCerts.length == 0 || this.server == null || this.server.certId.length == 0) {
		this.certTab = "upload";
	}

	this.switchCertTab = function (tab) {
		this.certTab = tab;
	};
});