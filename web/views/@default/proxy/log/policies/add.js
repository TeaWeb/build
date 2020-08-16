Tea.context(function () {
	this.$delay(function () {
		this.$find("form input[name='name']").focus();
	});

	this.submitSuccess = function () {
		alert("添加成功");
		window.location = "/proxy/log/policies";
	};

	/**
	 * 存储格式
	 */
	this.storageFormat = this.formats[0].code;
	this.storageTemplate = "${remoteAddr} - [${timeLocal}] \"${request}\" ${status} ${bodyBytesSent} \"${http.Referer}\" \"${http.UserAgent}\"";
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
	this.storageType = "";
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

	/**
	 * syslog
	 */
	this.syslogProtocol = "none";

	/**
	 * 更多设置
	 */
	this.moreOptionsVisible = false;
});