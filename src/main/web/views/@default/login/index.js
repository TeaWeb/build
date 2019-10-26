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

	this.failSubmit = function (resp) {
		if (resp.message != null && resp.message.length > 0) {
			alert(resp.message);
			return;
		}
		if (resp.errors != null && resp.errors.length > 0) {
			var err = resp.errors[0];
			if (err.messages != null && err.messages.length > 0) {
				alert(err.messages[0]);
			}
			if (err.param == "username") {
				this.$refs.usernameRef.focus();
			} else if (err.param == "password") {
				this.$refs.passwordRef.focus();
			} else if (err.param == "refresh") {
				window.location.reload();
			}
		}
	};

	this.doneSubmit = function () {
		this.isSubmitting = false;
	};
});