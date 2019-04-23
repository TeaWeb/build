Tea.context(function () {
	this.$delay(function () {
		this.$find("form input[name='name']").focus();
	});

	this.submitSuccess = function (resp) {
		alert("保存成功");
		window.location = "/proxy/waf/group?wafId=" + this.config.id + "&groupId=" + this.group.id;
	};
});