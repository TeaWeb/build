Tea.context(function () {
	this.$delay(function () {
		this.$find("form input[name='endpoint']").focus();
	});

	this.submitSuccess = function () {
		window.location = "/proxy/tunnel?serverId=" + this.server.id;
	};

	this.generateSecret = function () {
		this.$post(".generateSecret")
			.success(function (resp) {
				this.tunnel.secret = resp.data.secret;
			});
	};
});