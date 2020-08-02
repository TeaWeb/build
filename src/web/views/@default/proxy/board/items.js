Tea.context(function () {
	/**
	 * 添加指标
	 */
	this.addItemFormVisible = false;
	this.selectedItemCode = "";
	this.selectedItem = null;

	this.showAddItemForm = function () {
		this.addItemFormVisible = !this.addItemFormVisible;
	};

	this.changeItem = function () {
		var that = this;
		this.selectedItem = this.items.$find(function (k, v) {
			return v.code == that.selectedItemCode;
		});
	};

	this.confirmAddItem = function () {
		if (this.selectedItemCode.length == 0) {
			alert("请选择要添加的指标");
			return;
		}
		this.$post(".addItem")
			.params({
				"code": this.selectedItemCode,
				"serverId": this.server.id
			})
			.refresh();
	};

	this.showRunningItemDetail = function (index) {
		this.runningItems.$each(function (k, item) {
			if (k == index) {
				return;
			}
			item.detailVisible = false;
		});

		var item = this.runningItems[index];
		if (item.detailVisible == null) {
			item.detailVisible = false;
		}
		item.detailVisible = !item.detailVisible;
		this.runningItems[index] = item;
		this.$set(this.runningItems, index, item);
	};

	this.deleteRunningItem = function (code) {
		if (!window.confirm("确定要删除指标吗？")) {
			return;
		}
		this.$post(".deleteItem")
			.params({
				"code": code,
				"serverId": this.server.id
			})
			.refresh();
	};

	/**
	 * 图表引用的指标
	 */
	this.showChartItemDetail = function (index) {
		this.chartItems.$each(function (k, item) {
			if (k == index) {
				return;
			}
			item.detailVisible = false;
		});

		var item = this.chartItems[index];
		if (item.detailVisible == null) {
			item.detailVisible = false;
		}
		item.detailVisible = !item.detailVisible;
		this.chartItems[index] = item;
		this.$set(this.chartItems, index, item);
	};
});