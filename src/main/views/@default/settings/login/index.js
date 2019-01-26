Tea.context(function () {
	this.generateKey = function () {
		this.$post("/settings/login/generateKey")
			.refresh();
	};
});