Tea.context(function () {
	this.colls = [];
	this.isLoaded = false;
	this.totalSize = 0;
	this.countCalculatedColls = 0;

	this.$delay(function () {
		this.loadColls();
		this.loadStats();
	});

	this.loadColls = function () {
		this.$post(".colls")
			.success(function (resp) {
				this.colls = resp.data.colls;
				this.colls.$each(function (k, v) {
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
		var collNames = [];
		var max = 10;
		this.colls.$each(function (k, v) {
			if (collNames.length >= max) {
				return;
			}
			if (v.size == null || v.size == "-") {
				collNames.push(v.name);
			}
		});
		if (collNames.length == 0) {
			this.$delay(function () {
				this.loadStats();
			}, 2000);
			return;
		}
		this.$post(".collStat")
			.params({
				"collNames": collNames
			})
			.success(function (resp) {
				this.colls = this.colls.$map(function (k, v) {
					if (typeof (resp.data.result[v.name]) != "undefined") {
						v["count"] = resp.data.result[v.name].count;
						v["size"] = resp.data.result[v.name].size;
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

	this.deleteColl = function (name) {
		if (!window.confirm("确定要删除" + name + "中的所有数据吗？")) {
			return;
		}
		this.$post(".deleteColl")
			.params({
				"collName": name
			})
			.success(function () {
				this.colls = this.colls.$filter(function (k, v) {
					return v.name != name;
				});
				this.reloadTotalSize();
			});
	};

	this.reloadTotalSize = function () {
		var totalSize = 0;
		var count = 0;
		this.colls.$each(function (k, v) {
			if (v.size != "-") {
				totalSize += parseFloat(v.size.replace("M", ""));
				count ++;
			}
		});
		this.totalSize = Math.ceil(totalSize * 100) / 100;
		this.countCalculatedColls = count;
	};
});