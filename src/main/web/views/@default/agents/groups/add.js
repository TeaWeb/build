Tea.context(function () {
	this.$delay(function () {
		this.$find("form input[name='name']").focus();
		teaweb.datepicker("day-from-input");
		teaweb.datepicker("day-to-input");
	});

	this.submitSuccess = function () {
		alert("添加成功");
		if (this.from.length > 0) {
			window.location = this.from;
		} else {
			window.location = "/agents/groups";
		}
	};

	/**
	 * 更多选项
	 */
	this.advancedOptionsVisible = false;

	this.showAdvanced = function () {
		this.advancedOptionsVisible = !this.advancedOptionsVisible;
	};
});