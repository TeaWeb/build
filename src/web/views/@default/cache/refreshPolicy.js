Tea.context(function () {
	this.prefixes = "";
	this.isSubmitting = false;
	this.isSuccess = false;
	this.count = 0;

	this.reset = function () {
		this.prefixes = "";
		this.$refs.focusInput.focus();
	};

	this.submitBefore = function () {
		this.isSubmitting = true;
		this.isSuccess = false;
	};

	this.submitSuccess = function (resp) {
		this.isSubmitting = false;
		this.isSuccess = true;
		this.count = resp.data.count;

		this.$delay(function () {
			this.isSuccess = false;
		}, 3000);
	};
});