Tea.context(function () {
	this.passwordUpdating = false;

	this.updatePassword = function () {
		this.passwordUpdating = !this.passwordUpdating;
	};
});