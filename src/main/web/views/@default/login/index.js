Tea.context(function () {
	this.username = "";
	this.password = "";
	this.passwordMd5 = "";

	this.isSubmitting = false;

	if (this.teaDemoEnabled) {
		this.username = "admin";
		this.password = "123456";
	}

	this.$delay(function () {
		this.$find("form input[name='username']").focus();
		this.changePassword();
	});

	this.changePassword = function () {
		this.passwordMd5 = md5(this.password.trim());
	};

	// 更多选项
	this.moreOptionsVisible = false;
	this.showMoreOptions = function () {
		this.moreOptionsVisible = !this.moreOptionsVisible;
	};

	this.beforeSubmit = function () {
		this.isSubmitting = true;
	};

	this.doneSubmit = function () {
		this.isSubmitting = false;
	};
});