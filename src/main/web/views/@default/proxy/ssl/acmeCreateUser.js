Tea.context(function () {
	this.$delay(function () {
		this.$find("form input[name='email']").focus();
	});

	this.submitSuccess = function () {
		alert("创建成功");
		window.location = "/proxy/ssl/acmeUsers?serverId=" + this.server.id;
	};
});