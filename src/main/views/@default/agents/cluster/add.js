Tea.context(function () {
	this.currentStep = "search";

	/**
	 * 搜索主机
	 */
	this.subStep = "form";
	this.hosts = []; // { addr, ip, name, cost, canConnect, state, result, isInstalled }
	this.port = 22;
	this.searchId = 0;

	this.$delay(function () {
		this.$find("#search-form textarea").focus();
	});

	this.initSearching = function () {
		this.countTestedHosts = 0;
		this.countFoundHosts = 0;
		this.countSelectedHosts = 0;
		this.connectTimer = null;
		this.isSearching = true;
		this.selectedHosts = [];
		this.searchId++;
	};

	this.initSearching();

	this.searchDone = function (resp) {
		this.hosts = resp.data.hosts;
		this.subStep = "hosts";
		this.initSearching();

		this.$delay(function () {
			this.connect(this.searchId);
		}, 1000);
	};

	this.connect = function (searchId) {
		this.isSearching = true;

		if (this.connectTimer != null) {
			clearTimeout(this.connectTimer);
			this.connectTimer = null;
		}

		var max = 20; // 每次最多处理多个host
		var paramHosts = [];
		this.hosts.$each(function (k, v) {
			if (paramHosts.length >= max) {
				return;
			}
			if (v.state == "WAITING") {
				paramHosts.push(v.addr);
				v.result = "扫描中";
			}
		});
		if (paramHosts.length == 0) {
			this.isSearching = false;
			return;
		}

		this.$post(".conn")
			.params({
				"hosts": paramHosts,
				"port": this.port
			})
			.success(function (resp) {
				// 如果是不同批次的搜索则放弃，防止结果重合
				if (this.searchId != searchId) {
					return;
				}

				var that = this;
				resp.data.states.$each(function (k, state) {
					var host = that.hosts.$find(function (k, v) {
						return v.addr == state.addr;
					});
					host.cost = state.cost;
					host.canConnect = state.canConnect;
					host.isChecked = host.canConnect;
					host.ip = state.ip;
					host.name = state.name;
					host.state = "READY";
					host.result = "扫描完毕";
					that.countTestedHosts++;
					if (host.canConnect) {
						that.countFoundHosts++;
						that.countSelectedHosts++;
						host.result = (Math.ceil(state.cost * 1000) / 1000) + "ms";
					}
				});

				// 排序
				var newHosts = [];
				this.hosts.$each(function (k, v) {
					if (v.canConnect) {
						newHosts.push(v);
					}
				});
				this.hosts.$each(function (k, v) {
					if (!v.canConnect) {
						newHosts.push(v);
					}
				});
				this.hosts = newHosts;

				if (this.countTestedHosts >= this.hosts.length) {
					this.isSearching = false;
				}
			})
			.done(function () {
				// 如果是不同批次的搜索则放弃，防止结果重合
				if (this.searchId != searchId) {
					return;
				}

				var that = this;
				this.connectTimer = setTimeout(function () {
					that.connect(that.searchId);
				}, 1000)
			});
	};

	this.goSearchForm = function () {
		this.subStep = "form";
		this.searchId = 0;

		if (this.connectTimer != null) {
			clearTimeout(this.connectTimer);
			this.connectTimer = null;
		}
	};

	this.changeChecked = function (host) {
		host.isChecked = !host.isChecked;
		if (host.isChecked) {
			this.countSelectedHosts++;
		} else {
			this.countSelectedHosts--;
		}
	};

	/**
	 * 认证
	 */
	this.authMasterURL = window.location.protocol + "//" + window.location.host;
	this.authUsername = "root";
	this.authPassword = "";
	this.installDir = "/opt/teaweb";
	this.authGroupId = "";
	this.authType = "password";
	this.authKey = "";

	this.goAuth = function () {
		this.currentStep = "auth";
		this.authKey = "";
		this.selectedHosts = this.hosts.$findAll(function (k, v) {
			return v.isChecked;
		});
	};

	this.selectAuthType = function (authType) {
		this.authType = authType;
	};

	this.goSearchResult = function () {
		this.currentStep = "search";
		this.subStep = "hosts";
	};

	this.authSuccess = function (resp) {
		this.currentStep = "install";
		this.authKey = resp.data.key;
		this.$delay(function () {
			this.startInstall();
		});
	};

	/**
	 * 部署
	 */
	this.countInstalledHosts = 0;
	this.isInstalling = false;

	this.startInstall = function () {
		this.isInstalling = true;
		this.selectedHosts.$each(function (k, v) {
			v.state = "READY";
			v.isInstalled = false;
			v.result = "";
			v.hasError = false;
		});
		this.countInstalledHosts = 0;
		this.$delay(function () {
			this.install();
		});
	};

	this.install = function () {
		var max = 10;
		var hosts = [];
		this.selectedHosts.$each(function (k, v) {
			if (hosts.length >= max) {
				return;
			}
			if (v.state == "READY") {
				v.state = "INSTALLING";
				hosts.push(v.addr);
			}
		});

		// 安装完毕
		if (hosts.length == 0) {
			this.isInstalling = false;
			return;
		}

		this.$post(".install")
			.params({
				"hosts": hosts,
				"master": this.authMasterURL,
				"dir": this.installDir,
				"username": this.authUsername,
				"password": this.authPassword,
				"port": this.port,
				"groupId": this.authGroupId,
				"authType": this.authType,
				"key": this.authKey
			})
			.timeout(300)
			.success(function (resp) {
				var that = this;
				resp.data.states.$each(function (k, state) {
					var host = that.selectedHosts.$find(function (k, v) {
						return state.addr == v.addr;
					});
					host.isInstalled = state.isInstalled;
					if (host.isInstalled) {
						that.countInstalledHosts++;
					}
					host.result = state.result;
					host.hasError = state.hasError;
					if (state.ip != null && state.ip.length > 0) {
						host.ip = state.ip;
					}
					if (state.name != null && state.name.length > 0) {
						host.name = state.name;
					}
				});
			})
			.done(function () {
				this.$delay(function () {
					this.install();
				}, 1000);
			});
	};

	this.goBackAuth = function () {
		this.currentStep = "auth";
	};

	this.finish = function () {
		this.currentStep = "finish";
	};
});