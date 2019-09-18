Tea.context(function () {
	this.$delay(function () {
		this.$find("form input[name='name']").focus();

		teaweb.datepicker("day-from-input");
		teaweb.datepicker("day-to-input");
	});

	this.submitSuccess = function () {
		alert("保存成功");
		window.location = "/agents/groups/detail?groupId=" + this.group.id;
	};

	/**
	 * 更多选项
	 */
	this.advancedOptionsVisible = false;

	this.showAdvanced = function () {
		this.advancedOptionsVisible = !this.advancedOptionsVisible;
	};
});