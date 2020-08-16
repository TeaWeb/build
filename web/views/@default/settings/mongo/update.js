Tea.context(function () {
	this.isTesting = false;
	this.testingError = "";
	this.testingSuccess = "";

	if (this.config.authMechanism.length == 0) {
		this.config.authMechanism = "SCRAM-SHA-1";
	}

	this.updateAuth = function () {
		this.config.authEnabled = !this.config.authEnabled;
	};

	this.testConnection = function () {
		var params = {
			host: this.config.host,
			port: this.config.port,
			dbName: this.config.dbName,
			authEnabled: this.config.authEnabled ? 1 : 0
		};
		if (this.config.authEnabled) {
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
			.timeout(10)
			.success(function () {
				this.testingError = "";
				this.testingSuccess = "连接成功！";
			})
			.fail(function (resp) {
				if (resp) {
					this.testingError = resp.message;
				}
			})
			.error(function () {
				this.testingError = "连接超时";
			})
			.done(function () {
				this.isTesting = false;
			});
	};

	/**
	 * 高级选项
	 */
	this.advancedOptionsVisible = false;

	this.showAdvancedOptions = function () {
		this.advancedOptionsVisible = !this.advancedOptionsVisible;
	};
});