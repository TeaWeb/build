Tea.context(function () {
	this.$delay(function () {
		if (this.certs.length == 0) {
			this.addCert();
		}
		this.sortable();
	});

	this.httpsOn = (this.server.ssl != null && this.server.ssl.on);
	this.http2Enabled = (this.server.ssl == null || !this.server.ssl.http2Disabled);

	if (this.server.ssl == null) {
		this.server.ssl = {
			"certificate": "",
			"certificateKey": "",
			"listen": []
		};
	}
	if (this.server.ssl.listen == null) {
		this.server.ssl.listen = [];
	}
	if (this.server.ssl.listen.length == 0) {
		this.server.ssl.listen = ["0.0.0.0:443"];
	}

	this.submitSuccess = function () {
		alert("修改成功");

		window.location = "/proxy/ssl?serverId=" + this.server.id;
	};

	/**
	 * 证书
	 */
	this.certIndex = 0;

	this.addCert = function () {
		this.certs.push({
			"description": "",
			"certFile": "",
			"keyFile": "",
			"isLocal": false,
			"isShared": true
		});
		this.certIndex = this.certs.length - 1;
		this.$delay(function () {
			this.$find("#cert-descriptions-input-" + (this.certs.length - 1)).focus();
		});
	};

	this.removeCert = function (index) {
		if (!window.confirm("确定要移除此证书吗？")) {
			return;
		}

		this.certs.$remove(index);
		if (this.certs.length > 0) {
			if (index >= 1) {
				this.certIndex = index - 1;
			} else {
				this.certIndex = 0;
			}
		}
	};

	/**
	 * 绑定的网络地址
	 */
	this.listenAdding = false;
	this.addingListenName = "";
	this.editingListenIndex = -1;

	this.addListen = function () {
		this.listenAdding = true;
		this.editingListenIndex = -1;
		this.$delay(function () {
			this.$find("form input[name='addingListenName']").focus();
		});
	};

	this.editListen = function (index) {
		this.listenAdding = true;
		this.editingListenIndex = index;
		this.$delay(function () {
			this.$find("form input[name='addingListenName']").focus();
		});
		this.addingListenName = this.server.ssl.listen[index];
	};

	this.confirmAddListen = function () {
		this.addingListenName = this.addingListenName.trim();
		if (this.addingListenName.length == 0) {
			alert("网络地址不能为空");
			this.$find("form input[name='addingListenName']").focus();
			return;
		}
		if (this.editingListenIndex > -1) {
			this.server.ssl.listen[this.editingListenIndex] = this.addingListenName;
		} else {
			this.server.ssl.listen.push(this.addingListenName);
		}
		this.cancelListenAdding();
	};

	this.cancelListenAdding = function () {
		this.listenAdding = false;
		this.addingListenName = "";
		this.editingListenIndex = -1;
	};

	this.removeListen = function (index) {
		this.server.ssl.listen.$remove(index);
		this.cancelListenAdding();
	};

	/**
	 * 高级设置
	 */
	this.advancedOptionsVisible = false;
	this.cipherSuitesOn = (this.server.ssl.cipherSuites != null && this.server.ssl.cipherSuites.length > 0);
	this.selectedCipherSuites = [];
	if (this.server.ssl.cipherSuites != null && this.server.ssl.cipherSuites.length > 0) {
		this.selectedCipherSuites = this.server.ssl.cipherSuites;
	}
	var allCipherSuites = this.cipherSuites.$copy();
	if (this.selectedCipherSuites.length > 0) {
		var that = this;
		this.cipherSuites = allCipherSuites.$findAll(function (k, v) {
			return !that.selectedCipherSuites.$contains(v);
		});
	}

	this.showAdvancedOptions = function () {
		this.advancedOptionsVisible = !this.advancedOptionsVisible;
	};

	this.addCipherSuite = function (cipherSuite) {
		if (!this.selectedCipherSuites.$contains(cipherSuite)) {
			this.selectedCipherSuites.push(cipherSuite);
		}
		this.cipherSuites.$removeValue(cipherSuite);
	};

	this.removeCipherSuite = function (cipherSuite) {
		if (!window.confirm("确定要删除此套件吗？")) {
			return;
		}
		this.selectedCipherSuites.$removeValue(cipherSuite);

		var that = this;
		this.cipherSuites = allCipherSuites.$findAll(function (k, v) {
			return !that.selectedCipherSuites.$contains(v);
		});
	};

	this.formatCipherSuite = function (cipherSuite) {
		return cipherSuite.replace(/(AES|3DES)/, "<var>$1</var>");
	};

	this.addBatchCipherSuites = function (suites) {
		var that = this;
		suites.$each(function (k, v) {
			if (that.selectedCipherSuites.$contains(v)) {
				return;
			}
			that.selectedCipherSuites.push(v);
		});
	};

	this.clearCipherSuites = function () {
		this.selectedCipherSuites = [];
		var that = this;
		this.cipherSuites = allCipherSuites.$findAll(function (k, v) {
			return !that.selectedCipherSuites.$contains(v);
		});
	};

	/**
	 * 拖动排序
	 */
	this.sortable = function () {
		var box = this.$find(".cipher-suites-box")[0];
		var that = this;
		Sortable.create(box, {
			draggable: ".label",
			handle: ".icon.handle",
			onStart: function () {

			},
			onUpdate: function (event) {

			}
		});
	};

	/**
	 * hsts
	 */
	this.hstsOptionsVisible = false;
	if (this.hsts.domains == null) {
		this.hsts.domains = [];
	}

	this.showMoreHSTS = function () {
		this.hstsOptionsVisible = !this.hstsOptionsVisible;
		if (this.hstsOptionsVisible) {
			this.changeHSTSMaxAge();
			this.$delay(function () {
				window.scroll(0, 10000);
			});
		}
	};

	this.changeHSTSMaxAge = function () {
		var v = this.hsts.maxAge;
		if (isNaN(v)) {
			this.hsts.days = "-";
			return;
		}
		this.hsts.days = parseInt(v / 86400);
		if (isNaN(this.hsts.days)) {
			this.hsts.days = "-";
		} else if (this.hsts.days < 0) {
			this.hsts.days = "-";
		}
	};

	this.setHSTSMaxAge = function (maxAge) {
		this.hsts.maxAge = maxAge;
		this.changeHSTSMaxAge();
	};

	/**
	 * 域名
	 */
	this.hstsDomainAdding = false;
	this.addingHstsDomain = "";
	this.hstsDomainEditingIndex = -1;

	this.addHstsDomain = function () {
		this.hstsDomainAdding = true;
		this.hstsDomainEditingIndex = -1;
		this.$delay(function () {
			this.$find("form input[name='addingHstsDomain']").focus();
		});
	};

	this.editHstsDomain = function (index) {
		this.hstsDomainEditingIndex = index;
		this.addingHstsDomain = this.hsts.domains[index];
		this.hstsDomainAdding = true;
		this.$delay(function () {
			this.$find("form input[name='addingHstsDomain']").focus();
		});
	};

	this.confirmAddHstsDomain = function () {
		this.addingHstsDomain = this.addingHstsDomain.trim();
		if (this.addingHstsDomain.length == 0) {
			alert("域名不能为空");
			this.$find("form input[name='addingHstsDomain']").focus();
			return;
		}
		if (this.hstsDomainEditingIndex > -1) {
			this.hsts.domains[this.hstsDomainEditingIndex] = this.addingHstsDomain;
		} else {
			this.hsts.domains.push(this.addingHstsDomain);
		}
		this.cancelHstsDomainAdding();
	};

	this.cancelHstsDomainAdding = function () {
		this.hstsDomainAdding = false;
		this.addingHstsDomain = "";
		this.hstsDomainEditingIndex = -1;
	};

	this.removeHstsDomain = function (index) {
		this.cancelHstsDomainAdding();
		this.hsts.domains.$remove(index);
	};

	/**
	 * CA证书
	 */
	this.caCertsVisible = false;

	this.showCACerts = function () {
		this.caCertsVisible = !this.caCertsVisible;
	};

	var that = this;
	this.selectedCACerts = this.caCerts.$filter(function (k, v) {
		v.isSelected = that.clientCACertIds.$contains(v.id);
		return v.isSelected;
	});
	this.leftCACerts = this.caCerts.$filter(function (k, v) {
		v.isSelected = that.clientCACertIds.$contains(v.id);
		return !v.isSelected;
	});

	this.selectCACert = function (cert) {
		cert.isSelected = true;
		this.selectedCACerts = this.caCerts.$findAll(function (k, v) {
			return v.isSelected;
		});
		this.leftCACerts = this.caCerts.$findAll(function (k, v) {
			return !v.isSelected;
		});
	};

	this.removeSelectedCACert = function (cert) {
		cert.isSelected = false;
		this.selectedCACerts = this.caCerts.$findAll(function (k, v) {
			return v.isSelected;
		});
		this.leftCACerts = this.caCerts.$findAll(function (k, v) {
			return !v.isSelected;
		});
	};
});