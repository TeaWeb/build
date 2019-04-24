Tea.context(function () {
	this.$delay(function () {
		this.$find("form input[name='name']").focus();
	});

	this.submitSuccess = function (resp) {
		alert("保存成功");
		window.location = "/proxy/waf/group?wafId=" + this.config.id + "&groupId=" + resp.data.id + "&inbound=" + (this.inbound ? 1 : 0);
	};
});