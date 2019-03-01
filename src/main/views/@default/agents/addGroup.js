Tea.context(function () {
	this.$delay(function () {
		this.$find("form input[name='name']").focus();
	});

	this.submitSuccess = function () {
		alert("添加成功");
		if (this.from.length > 0) {
			window.location = this.from;
		} else {
			window.location = "/agents/groups";
		}
	};
});