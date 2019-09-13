Tea.context(function () {
	this.isLoaded = false;
	this.values = [];
	this.lastId = "";

	this.$delay(function () {
		this.testMongo();
		this.loadValues();
	});

	this.loadValues = function () {
		this.$post("$")
			.params({
				"agentId": this.agentId,
				"appId": this.app.id,
				"itemId": this.item.id,
				"lastId": this.lastId,
				"level": this.selectedLevel
			})
			.success(function (resp) {
				if (resp.data.values.length == 0) {
					return;
				}
				this.lastId = resp.data.values.$first().id;
				this.values = resp.data.values.$map(function (k, v) {
					v.value = JSON.stringify(v.value, 1, "  ");
					v.costMs = Math.ceil(v.costMs * 1000) / 1000;
					return v;
				}).concat(this.values);
			})
			.fail(function (resp) {
				if (resp != null) {
					console.log(resp.message);
				}
			})
			.done(function () {
				this.$delay(function () {
					this.isLoaded = true;
				});

				this.$delay(function () {
					this.loadValues();
				}, 3000);
			});
	};

	this.clearValues = function () {
		if (!window.confirm("确定要清除所有数值记录吗？")) {
			return;
		}
		this.$post("/agents/apps/clearItemValues")
			.params({
				"agentId": this.agentId,
				"appId": this.app.id,
				"itemId": this.item.id,
				"level": this.selectedLevel
			})
			.refresh();
	};

	this.showValueTab = function (valueIndex, tab) {
		this.values[valueIndex].tab = tab;
		this.$set(this.values, valueIndex, this.values[valueIndex]);
	};
});