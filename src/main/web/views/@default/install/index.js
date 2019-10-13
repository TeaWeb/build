Tea.context(function () {
	this.skip = function () {
		if (!window.confirm("确定要先跳过这一步吗？可以在后期再设置。")) {
			return;
		}
		this.$post(".skip")
			.success(function () {
				window.location = "/dashboard";
			});
	};

	this.goStep = function (step1) {
		this.step = step1;
	};

	this.step = "dbType";
	this.dbType = "mongo";

	/**
	 * 选择数据库类型
	 */
	this.saveDBType = function () {
		switch (this.dbType) {
			case "mysql":
				this.dbAddr = "127.0.0.1:3306";
				this.dbUsername = "root";
				break;
			case "postgres":
				this.dbAddr = "127.0.0.1:5432";
				this.dbUsername = "postgres";
				break;
		}
		this.goStep("config");
	};

	/**
	 * MongoDB配置
	 */
	this.mongoAddr = "127.0.0.1:27017";
	this.mongoAuthEnabled = false;
	this.mongoUsername = "";
	this.mongoPassword = "";
	this.mongoAuthMechanism = "SCRAM-SHA-1";
	this.mongoAuthMechanismProperties = "";
	this.mongoDBName = "teaweb";

	/**
	 * 数据库设置
	 */
	this.dbAddr = "";
	this.dbUsername = "";
	this.dbPassword = "";
	this.dbName = "teaweb";
	this.autoCreate = true;
	this.dbTestResult = {
		"ok": false,
		"isRunning": false,
		"message": ""
	};

	this.testDB = function () {
		this.dbTestResult.ok = false;
		this.dbTestResult.message = "";
		this.dbTestResult.isRunning = true;

		var params = {
			"dbType": this.dbType,
			"addr": this.dbAddr,
			"username": this.dbUsername,
			"password": this.dbPassword,
			"dbName": this.dbName,
			"autoCreate": this.autoCreate ? 1 : 0
		};
		if (this.dbType == "mongo") {
			params = {
				"dbType": this.dbType,
				"addr": this.mongoAddr,
				"authEnabled": this.mongoAuthEnabled ? 1 : 0,
				"username": this.mongoUsername,
				"password": this.mongoPassword,
				"authMechanism": this.mongoAuthMechanism,
				"authMechanismProperties": this.mongoAuthMechanismProperties,
				"dbName": this.mongoDBName
			};
		}
		this.$post(".test")
			.params(params)
			.timeout(10)
			.success(function () {
				this.dbTestResult.ok = true;
			})
			.fail(function (resp) {
				this.dbTestResult.message = resp.message;
				this.dbTestResult.ok = false;
			})
			.error(function () {
				this.dbTestResult.message = "数据库连接超时";
				this.dbTestResult.ok = false;
			})
			.done(function () {
				this.dbTestResult.isRunning = false;
			});
	};

	/**
	 * 保存
	 */
	this.saveSuccess = false;
	this.saveFailed = false;
	this.saveMessage = "";

	this.saveDB = function () {
		this.saveSuccess = false;
		this.saveFailed = false;
		this.saveMessage = "";

		var params = {
			"dbType": this.dbType,
			"addr": this.dbAddr,
			"username": this.dbUsername,
			"password": this.dbPassword,
			"dbName": this.dbName,
			"autoCreate": this.autoCreate ? 1 : 0
		};
		if (this.dbType == "mongo") {
			params = {
				"dbType": this.dbType,
				"addr": this.mongoAddr,
				"authEnabled": this.mongoAuthEnabled ? 1 : 0,
				"username": this.mongoUsername,
				"password": this.mongoPassword,
				"authMechanism": this.mongoAuthMechanism,
				"authMechanismProperties": this.mongoAuthMechanismProperties,
				"dbName": this.mongoDBName
			};
		}

		this.$post(".save")
			.params(params)
			.success(function () {
				this.saveSuccess = true;
				this.goStep("finish");
			})
			.fail(function (resp) {
				this.saveFailed = true;
				this.saveMessage = resp.message;
			});
	};
});