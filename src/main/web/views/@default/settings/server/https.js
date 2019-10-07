Tea.context(function () {
	this.addresses = (this.server != null) ? this.server.https.listen.join("\n") : [];

	/**
	 * 证书
	 */
	this.certTab = "shared";
	if (this.sharedCerts.length == 0 || this.server.certId.length == 0) {
		this.certTab = "upload";
	}

	this.switchCertTab = function (tab) {
		this.certTab = tab;
	};
});