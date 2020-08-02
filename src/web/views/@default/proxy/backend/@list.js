Tea.context(function () {
	this.from = encodeURIComponent(window.location.toString());
	this.currentLocation = window.location.toString();
	this.query = Tea.serialize(this.queryParams);
	this.inWebsocketPage = (this.currentLocation.indexOf("/websocket") > 0) ? "1" : "0";

	this.padZero = function (s) {
		s = s.toString();
		if (s.length == 1) {
			return "0" + s;
		}
		return s;
	};

	var that = this;
	this.allNormalBackends = [];
	this.allBackupBackends = [];
	this.normalBackends = [];
	this.backupBackends = [];
	this.scheduling = null;
	this.isLoaded = false;

	this.groups = [];
	this.selectedGroupId = "default";

	this.$delay(function () {
		this.loadData();
	});

	this.$delay(function () {
		// scroll to bottom
		if (window.location.hash == "#scheduling") {
			window.scrollTo(0, 10000);
		}
	}, 300);

	this.loadData = function () {
		this.$get("/proxy/backend/data")
			.params(this.queryParams)
			.success(function (resp) {
				var that = this;
				this.allNormalBackends = resp.data.normalBackends.$map(function (k, v) {
					if (v.isDown) {
						var date = new Date(v.downTime);
						v.downTime = that.padZero(date.getMonth() + 1) + "-" + that.padZero(date.getDate()) + " " + that.padZero(date.getHours()) + ":" + that.padZero(date.getMinutes()) + ":" + that.padZero(date.getSeconds());
					}
					return v;
				});
				this.allBackupBackends = resp.data.backupBackends.$map(function (k, v) {
					if (v.isDown) {
						var date = new Date(v.downTime);
						v.downTime = that.padZero(date.getMonth() + 1) + "-" + that.padZero(date.getDate()) + " " + that.padZero(date.getHours()) + ":" + that.padZero(date.getMinutes()) + ":" + that.padZero(date.getSeconds());
					}
					return v;
				});

				this.normalBackends = this.allNormalBackends.$filter(function (k, v) {
					if (v.requestGroupIds == null || v.requestGroupIds.length == 0) {
						return that.selectedGroupId == "default";
					}

					return v.requestGroupIds.$contains(that.selectedGroupId);
				});
				this.backupBackends = this.allBackupBackends.$filter(function (k, v) {
					if (v.requestGroupIds == null || v.requestGroupIds.length == 0) {
						return that.selectedGroupId == "default";
					}

					return v.requestGroupIds.$contains(that.selectedGroupId);
				});

				this.scheduling = resp.data.scheduling;
				var that = this;
				this.groups = resp.data.groups.$map(function (_, group) {
					var count = 0;
					that.allNormalBackends.$each(function (_, backend) {
						// 默认分组
						if ((backend.requestGroupIds == null || backend.requestGroupIds.length == 0) && group.id == "default") {
							count++;
						}

						// 正常分组
						if (backend.requestGroupIds != null && backend.requestGroupIds.$contains(group.id)) {
							count++;
						}
					});
					that.allBackupBackends.$each(function (_, backend) {
						// 默认分组
						if ((backend.requestGroupIds == null || backend.requestGroupIds.length == 0) && group.id == "default") {
							count++;
						}

						// 正常分组
						if (backend.requestGroupIds != null && backend.requestGroupIds.$contains(group.id)) {
							count++;
						}
					});
					group.count = count;
					return group;
				});
			})
			.done(function () {
				this.isLoaded = true;
				this.$delay(function () {
					this.loadData();
				}, 5000);
			})
			.fail(function (resp) {
				console.log(resp.message);
			});
	};

	this.selectGroup = function (groupId) {
		this.selectedGroupId = groupId;

		this.normalBackends = this.allNormalBackends.$filter(function (k, v) {
			if (v.requestGroupIds == null || v.requestGroupIds.length == 0) {
				return that.selectedGroupId == "default";
			}

			return v.requestGroupIds.$contains(that.selectedGroupId);
		});
		this.backupBackends = this.allBackupBackends.$filter(function (k, v) {
			if (v.requestGroupIds == null || v.requestGroupIds.length == 0) {
				return that.selectedGroupId == "default";
			}

			return v.requestGroupIds.$contains(that.selectedGroupId);
		});
	};

	this.deleteBackend = function (backendId) {
		if (!window.confirm("确定要删除此服务器吗？")) {
			return;
		}
		var query = this.queryParams;
		query["backendId"] = backendId;
		this.$post("/proxy/backend/delete")
			.params(query);
	};

	this.putOnline = function (backend) {
		if (!window.confirm("确定要上线此服务器吗？")) {
			return;
		}
		var query = this.queryParams;
		query["backendId"] = backend.id;
		this.$post("/proxy/backend/online")
			.params(query)
			.success(function () {
				backend.isDown = false;
				backend.currentFails = 0;
			});
	};

	this.clearFails = function (backend) {
		if (!window.confirm("确定要清除此服务器的失败次数吗？此操作不会改变上线状态")) {
			return;
		}
		var query = this.queryParams;
		query["backendId"] = backend.id;
		this.$post("/proxy/backend/clearFails")
			.params(query)
			.success(function () {
				backend.currentFails = 0;
			});
	};
});