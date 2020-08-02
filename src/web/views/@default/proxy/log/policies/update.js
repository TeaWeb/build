Tea.context(function () {
	this.$delay(function () {
		this.$find("form input[name='name']").focus();
	});

	this.submitSuccess = function () {
		alert("修改成功");
		window.location = "/proxy/log/policies/policy?policyId=" + this.policy.id;
	};

	/**
	 * 存储格式
	 */
	this.storageFormat = this.policy.options.format;
	this.storageTemplate = this.policy.options.template;

	if (this.storageFormat != "template") {
		this.storageTemplate = "${remoteAddr} - [${timeLocal}] \"${request}\" ${status} ${bodyBytesSent} \"${http.Referer}\" \"${http.UserAgent}\"";
	}

	this.selectedFormat = null;

	this.changeStorageFormat = function () {
		var that = this;
		this.selectedFormat = this.formats.$find(function (k, v) {
			return v.code == that.storageFormat;
		});
	};
	this.changeStorageFormat();

	/**
	 * 存储类型
	 */
	this.storageType = this.policy.type;
	this.selectedStorage = null;

	this.changeStorageType = function () {
		if (this.storageType == "") {
			return;
		}
		var that = this;
		this.selectedStorage = this.storages.$find(function (k, v) {
			return v.type == that.storageType;
		});
	};

	this.changeStorageType();

	/**
	 * syslog
	 */
	this.syslogProtocol = "";
	if (this.policy.type == "syslog") {
		this.syslogProtocol = this.policy.options.protocol;
	}

	/**
	 * 更多设置
	 */
	this.moreOptionsVisible = false;
});