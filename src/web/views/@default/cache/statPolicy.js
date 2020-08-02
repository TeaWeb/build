Tea.context(function () {
	this.result = null;
	this.isRunning = false;

	this.run = function () {
		this.isRunning = true;
		this.result = null;

		this.$post(".statPolicy")
			.params({
				"filename": this.policy.filename
			})
			.success(function (resp) {
				this.result = resp.data.result;
			})
			.done(function () {
				this.isRunning = false;
			});
	};
});