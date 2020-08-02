Tea.context(function () {
	this.advancedOptionsVisible = false;

	if (this.server.requestGroups != null) {
		var selectedRequestGroupIds = (this.backend.requestGroupIds == null) ? [] : this.backend.requestGroupIds;
		this.server.requestGroups.$each(function (k, v) {
			v.isChecked = selectedRequestGroupIds.$contains(v.id);
		});
	}

	this.showAdvancedOptions = function () {
		this.advancedOptionsVisible = !this.advancedOptionsVisible;
	};

	this.updateSuccess = function () {
		alert("保存成功");
		window.location = this.from;
	};

	/**
	 * 证书
	 */
	this.useCert = (this.backend.cert != null);
	if (this.backend.cert != null) {
		this.certId = this.backend.cert.id;
		this.certServerName = this.backend.cert.serverName;
	} else {
		this.certId = "";
		this.certServerName = "";
	}
	this.certVisible = false;

	this.$delay(function () {
		this.$watch("useCert", function () {
			this.certVisible = true;
		});
	});

	this.showCert = function () {
		this.certVisible = !this.certVisible;
	};

	/**
	 * 地址
	 */
	this.changeAddress = function () {
		if (this.backend.address == null) {
			return;
		}
		if (/^(http|https):\/\//i.test(this.backend.address)) {
			var pieces = this.backend.address.split("://");
			this.backend.scheme = pieces[0].toLocaleLowerCase();
			this.backend.address = pieces[1];
		}
		var index = this.backend.address.indexOf("/");
		if (index > -1) {
			this.backend.address = this.backend.address.substring(0, index);
		}
	};

	/**
	 * 主机名
	 */
	this.hostError = "";
	this.changeHost = function () {
		var host = this.backend.host.trim().replace(/[a-zA-Z0-9-\\.]/g, "");
		if (host.length > 0) {
			this.hostError = "主机名中含有特殊字符“" + host + "”，可能会导致后端服务器无法解析。";
		} else {
			this.hostError = "";
		}
	};

	/**
	 * request headers
	 */
	this.requestHeadersAdding = false;
	this.requestHeaders = [];
	if (this.backend.requestHeaders != null) {
		this.requestHeaders = this.backend.requestHeaders.$map(function (k, v) {
			return {
				"name": v.name,
				"value": v.value
			};
		});
	}
	this.requestHeadersAddingName = "";
	this.requestHeadersAddingValue = "";
	this.requestHeadersEditingIndex = -1;

	this.addRequestHeader = function () {
		this.requestHeadersAdding = true;
		this.requestHeadersAddingName = "";
		this.requestHeadersAddingValue = "";
		this.$delay(function () {
			this.$find("form input[name='requestHeaderName']").focus();
		});
	};

	this.cancelRequestHeadersAdding = function () {
		this.requestHeadersAdding = false;
		this.requestHeadersEditingIndex = -1;
	};

	this.confirmRequestHeadersAdding = function () {
		if (this.requestHeadersEditingIndex > -1) {
			this.requestHeaders[this.requestHeadersEditingIndex] = {
				"name": this.requestHeadersAddingName,
				"value": this.requestHeadersAddingValue
			};
		} else {
			this.requestHeaders.push({
				"name": this.requestHeadersAddingName,
				"value": this.requestHeadersAddingValue
			});
		}
		this.requestHeadersAddingName = "";
		this.requestHeadersAddingValue = "";
		this.cancelRequestHeadersAdding()
	};

	this.removeRequestHeader = function (index) {
		this.requestHeaders.$remove(index);
		this.cancelRequestHeadersAdding()
	};

	this.editRequestHeader = function (index) {
		this.requestHeadersEditingIndex = index;
		this.requestHeadersAdding = true;
		this.requestHeadersAddingName = this.requestHeaders[index].name;
		this.requestHeadersAddingValue = this.requestHeaders[index].value;
	};

	/**
	 * response headers
	 */
	this.responseHeadersAdding = false;
	this.responseHeaders = [];
	if (this.backend.responseHeaders != null) {
		this.responseHeaders = this.backend.responseHeaders.$map(function (k, v) {
			return {
				"name": v.name,
				"value": v.value
			};
		});
	}
	this.responseHeadersAddingName = "";
	this.responseHeadersAddingValue = "";
	this.responseHeadersEditingIndex = -1;

	this.addResponseHeader = function () {
		this.responseHeadersAdding = true;
		this.responseHeadersAddingName = "";
		this.responseHeadersAddingValue = "";
		this.$delay(function () {
			this.$find("form input[name='responseHeaderName']").focus();
		});
	};

	this.cancelResponseHeadersAdding = function () {
		this.responseHeadersAdding = false;
		this.responseHeadersEditingIndex = -1;
	};

	this.confirmResponseHeadersAdding = function () {
		if (this.responseHeadersEditingIndex > -1) {
			this.responseHeaders[this.responseHeadersEditingIndex] = {
				"name": this.responseHeadersAddingName,
				"value": this.responseHeadersAddingValue
			};
		} else {
			this.responseHeaders.push({
				"name": this.responseHeadersAddingName,
				"value": this.responseHeadersAddingValue
			});
		}
		this.responseHeadersAddingName = "";
		this.responseHeadersAddingValue = "";
		this.cancelResponseHeadersAdding()
	};

	this.removeResponseHeader = function (index) {
		this.responseHeaders.$remove(index);
		this.cancelResponseHeadersAdding()
	};

	this.editResponseHeader = function (index) {
		this.responseHeadersEditingIndex = index;
		this.responseHeadersAdding = true;
		this.responseHeadersAddingName = this.responseHeaders[index].name;
		this.responseHeadersAddingValue = this.responseHeaders[index].value;
	};

	/**
	 * 健康检查URL
	 */
	this.checkOn = this.backend.checkOn;
});