Tea.context(function () {
	this.connect = function () {
		this.error = "重试中...";
		this.$post(".connect").refresh();
	};
});