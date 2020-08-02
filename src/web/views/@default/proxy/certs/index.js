Tea.context(function () {
	this.countAll = 0;
	this.countActive = 0;
	this.countExpired = 0;
	this.countExpiring7Days = 0;
	this.countExpiring30Days = 0;
	this.countACME = 0;
	this.countCA = 0;

	this.timeTab = "all";

	var that = this;
	var allCerts = this.certs;
	this.certs.$each(function (k, cert) {
		that.countAll++;
		if (cert.isActive) {
			that.countActive++;
		}
		if (cert.isExpired) {
			that.countExpired++;
		}
		if (cert.isExpiring7Days) {
			that.countExpiring7Days++;
		}
		if (cert.isExpiring30Days) {
			that.countExpiring30Days++;
		}
		if (cert.isACME) {
			that.countACME++;
		}
		if (cert.isCA) {
			that.countCA++;
		}
	});

	this.selectTimeTab = function (tab) {
		this.timeTab = tab;
		this.certs = [];

		this.$delay(function () {
			this.certs = allCerts.$filter(function (k, cert) {
				if (tab == "all") {
					return true;
				}
				if (tab == "active") {
					return cert.isActive;
				}
				if (tab == "expired") {
					return cert.isExpired;
				}
				if (tab == "expiring7Days") {
					return cert.isExpiring7Days;
				}
				if (tab == "expiring30Days") {
					return cert.isExpiring30Days;
				}
				if (tab == "acme") {
					return cert.isACME;
				}

				if (tab == "ca") {
					return cert.isCA;
				}
			});
		}, 100);
	};

	this.deleteCert = function (certId) {
		if (!window.confirm("确定要删除此证书吗？")) {
			return;
		}
		this.$post(".delete")
			.params({
				"certId": certId
			})
			.refresh();
	};
});