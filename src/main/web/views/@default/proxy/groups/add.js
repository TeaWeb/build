Tea.context(function () {
	this.$delay(function () {
		this.$find("form input[name='name']").focus();
	});

	this.submitSuccess = function () {
		alert("保存成功");
		window.location = "/proxy/groups?serverId=" + this.server.id + "&locationId=" + this.locationId + "&websocket=" + this.websocket;
	};

	/**
	 * IP Range
	 */
	this.ipRangeType = "range";
	this.ipRangeAdding = false;
	this.ipRangeFrom = "";
	this.ipRangeTo = "";
	this.ipRangeCIDRIP = "";
	this.ipRangeCIDRBits = "";
	this.ipRanges = [];
	this.ipRangeEditingIndex = -1;
	this.ipRangeVar = "${remoteAddr}";

	this.selectIPRangeType = function (rangeType) {
		this.ipRangeType = rangeType;
		this.$delay(function () {
			if (this.ipRangeType == "range") {
				this.$find("form input[name='ipRangeFrom']").focus();
			} else if (this.ipRangeType == "cidr") {
				this.$find("form input[name='ipRangeCIDRIP']").focus();
			}
		});
	};

	this.addIPRange = function () {
		this.ipRangeAdding = true;
		this.ipRangeType = "range";
		this.ipRangeFrom = "";
		this.ipRangeTo = "";
		this.ipRangeCIDRIP = "";
		this.ipRangeCIDRBits = "";

		this.$delay(function () {
			if (this.ipRangeType == "range") {
				this.$find("form input[name='ipRangeFrom']").focus();
			} else if (this.ipRangeType == "cidr") {
				this.$find("form input[name='ipRangeCIDRIP']").focus();
			}
		});
	};

	this.cancelIPRangeAdding = function () {
		this.ipRangeAdding = false;
		this.ipRangeEditingIndex = -1;
	};

	this.confirmAddIPRange = function () {
		if (this.ipRangeType == "range") {
			if (!this.validateIP(this.ipRangeFrom)) {
				alert("请输入正确的开始IP");
				this.$find("form input[name='ipRangeFrom']").focus();
				return;
			}
			if (!this.validateIP(this.ipRangeTo)) {
				alert("请输入正确的结束IP");
				this.$find("form input[name='ipRangeTo']").focus();
				return;
			}
			if (this.ipRangeEditingIndex > -1) {
				this.ipRanges[this.ipRangeEditingIndex] = {
					"type": "range",
					"from": this.ipRangeFrom,
					"to": this.ipRangeTo,
					"var": this.ipRangeVar
				};
			} else {
				this.ipRanges.push({
					"type": "range",
					"from": this.ipRangeFrom,
					"to": this.ipRangeTo,
					"var": this.ipRangeVar
				});
			}
		} else if (this.ipRangeType == "cidr") {
			if (!this.validateIP(this.ipRangeCIDRIP)) {
				alert("请输入正确的IP地址");
				this.$find("form input[name='ipRangeCIDRIP']").focus();
				return;
			}
			if (!/^\d+$/.test(this.ipRangeCIDRBits)) {
				alert("请输入正确的位数");
				this.$find("form input[name='ipRangeCIDRBits']").focus();
				return;
			}
			if (this.ipRangeEditingIndex > -1) {
				this.ipRanges[this.ipRangeEditingIndex] = {
					"type": "cidr",
					"ip": this.ipRangeCIDRIP,
					"bits": this.ipRangeCIDRBits,
					"var": this.ipRangeVar
				};
			} else {
				this.ipRanges.push({
					"type": "cidr",
					"ip": this.ipRangeCIDRIP,
					"bits": this.ipRangeCIDRBits,
					"var": this.ipRangeVar
				});
			}
		}

		this.ipRangeAdding = false;
		this.ipRangeEditingIndex = -1;
	};

	this.removeIPRange = function (index) {
		this.ipRanges.$remove(index);
	};

	this.editIPRange = function (index) {
		var ipRange = this.ipRanges[index];
		this.ipRangeEditingIndex = index;
		this.ipRangeAdding = true;
		this.ipRangeType = ipRange.type;
		if (ipRange.type == "range") {
			this.ipRangeFrom = ipRange.from;
			this.ipRangeTo = ipRange.to;
			this.ipRangeVar = ipRange.var;
		} else if (ipRange.type == "cidr") {
			this.ipRangeCIDRIP = ipRange.ip;
			this.ipRangeCIDRBits = ipRange.bits;
			this.ipRangeVar = ipRange.var;
		}
	};

	this.validateIP = function (ip) {
		if (ip.length == 0) {
			return false;
		}
		var pieces = ip.split(".");
		if (pieces.length != 4) {
			return false;
		}
		var b = true;
		pieces.$each(function (k, v) {
			if (v.length > 3) {
				b = false;
				return;
			}
			if (!/^\d+$/.test(v)) {
				b = false;
				return;
			}
			v = parseInt(v, 10);
			if (v > 255) {
				b = false;
				return;
			}
		});
		return b;
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
	 * 更多选项
	 */
	this.advancedOptionsVisible = false;

	this.showAdvancedOptions = function () {
		this.advancedOptionsVisible = !this.advancedOptionsVisible;
	};
});