Tea.context(function () {
	this.authUpdating = this.config.username != null && this.config.username.length > 0;
	this.isTesting = false;
	this.testingError = "";
	this.testingSuccess = "";

	if (this.config.authMechanism.length == 0) {
		this.config.authMechanism = "SCRAM-SHA-1";
	}

	this.updateAuth = function () {
		this.authUpdating = !this.authUpdating;
	};

	this.testConnection = function () {
		var params = {
			host: this.config.host,
			port: this.config.port
		};
		if (this.authUpdating) {
			if (this.config.username.length == 0) {
				alert("请输入用户名");
				this.$find("form input[name='username']").focus();
				return;
			}
			params["username"] = this.config.username;
			params["password"] = this.config.password;
			params["authMechanism"] = this.config.authMechanism;
			params["authMechanismProperties"] = this.config.authMechanismProperties;
		}

		this.isTesting = true;
		this.testingError = "";
		this.testingSuccess = "";

		this.$get(".test")
			.params(params)
			.success(function () {
				this.testingError = "";
				this.testingSuccess = "连接成功！";
			})
			.fail(function (resp) {
				if (resp) {
					this.testingError = resp.message;
				}
			})
			.done(function () {
				this.isTesting = false;
			});
	};
});