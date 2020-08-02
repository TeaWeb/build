Tea.context(function () {
	this.advancedOptionsVisible = false;

	this.$delay(function () {
		this.$find("form input[name='address']").focus();
	});

	this.showAdvancedOptions = function () {
		this.advancedOptionsVisible = !this.advancedOptionsVisible;
	};

	this.submitSuccess = function () {
		alert("保存成功");
		window.location = this.from;
	};

	/**
	 * 证书
	 */
	this.useCert = false;

	/**
	 * 地址
	 */
	this.address = "";
	this.scheme = "http";
	if (this.isTCP) {
		this.scheme = "tcp";
	}

	this.changeAddress = function () {
		if (this.address == null) {
			return;
		}
		if (/^(http|https):\/\//i.test(this.address)) {
			var pieces = this.address.split("://");
			this.$find("select[name='scheme']").val(pieces[0].toLocaleLowerCase());
			this.address = pieces[1];
		}
		var index = this.address.indexOf("/");
		if (index > -1) {
			this.address = this.address.substring(0, index);
		}
	};

	/**
	 * 主机名
	 */
	this.host = "";
	this.hostError = "";
	this.changeHost = function () {
		var host = this.host.trim().replace(/[a-zA-Z0-9-\\.]/g, "");
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
	 * 健康检查
	 */
	this.checkOn = false;
});