Tea.context(function () {
	this.$delay(function () {
		if (this.certs.length == 0) {
			this.addCert();
		}
		this.sortable();
	});

	this.httpsOn = (this.server.ssl != null && this.server.ssl.on);

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
	this.addCert = function () {
		this.certs.push({
			"description": "",
			"certFile": "",
			"keyFile": "",
			"isLocal": false
		});
		this.$delay(function () {
			this.$find("#cert-descriptions-input-" + (this.certs.length - 1)).focus();
		});
	};

	this.removeCert = function (index) {
		if (!window.confirm("确定要移除此证书吗？")) {
			return;
		}

		this.certs.$remove(index);
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
});