Tea.context(function () {
	this.tables = [];
	this.isLoaded = false;
	this.totalSize = 0;
	this.formatedTotalSize = "";
	this.countCalculatedTables = 0;

	this.$delay(function () {
		this.loadTables();
		this.loadStats();
	});

	this.loadTables = function () {
		this.$post("/settings/database/tables")
			.success(function (resp) {
				this.tables = resp.data.tables;
				this.tables.$each(function (k, v) {
					v.count = "-";
					v.size = "-";
				});
			})
			.done(function () {
				this.$delay(function () {
					this.isLoaded = true;
				}, 100);
			});
	};

	this.loadStats = function () {
		var tableNames = [];
		var max = 10;
		this.tables.$each(function (k, v) {
			if (tableNames.length >= max) {
				return;
			}
			if (v.size == null || v.size == "-") {
				tableNames.push(v.name);
			}
		});
		if (tableNames.length == 0) {
			this.$delay(function () {
				this.loadStats();
			}, 2000);
			return;
		}
		this.$post("/settings/database/tableStat")
			.params({
				"tables": tableNames
			})
			.success(function (resp) {
				this.tables = this.tables.$map(function (k, v) {
					if (typeof (resp.data.result[v.name]) != "undefined") {
						v["count"] = resp.data.result[v.name].count;
						v["size"] = resp.data.result[v.name].size;

						v["formattedSize"] = resp.data.result[v.name].formattedSize;
					}
					return v;
				});

				this.reloadTotalSize();
			})
			.done(function () {
				this.$delay(function () {
					this.loadStats();
				}, 1000);
			});
	};

	this.deleteTable = function (name) {
		if (!window.confirm("确定要删除" + name + "中的所有数据吗？")) {
			return;
		}
		this.$post("/settings/database/deleteTable")
			.params({
				"table": name
			})
			.success(function () {
				this.tables = this.tables.$filter(function (k, v) {
					return v.name != name;
				});
				this.reloadTotalSize();
			});
	};

	this.reloadTotalSize = function () {
		var totalSize = 0;
		var count = 0;
		this.tables.$each(function (k, v) {
			if (v.size != "-") {
				totalSize += v.size;
				count++;
			}
		});
		this.totalSize = totalSize;
		this.countCalculatedTables = count;
		if (totalSize < 1024 * 1024) {
			this.formatedTotalSize = (Math.ceil(totalSize / 1024 * 100) / 100) + "KB";
		} else if (totalSize < 1024 * 1024 * 1024) {
			this.formatedTotalSize = (Math.ceil(totalSize / 1024 / 1024 * 100) / 100) + "MB";
		} else {
			this.formatedTotalSize = (Math.ceil(totalSize / 1024 / 1024 / 1024 * 100) / 100) + "GB";
		}
	};
});