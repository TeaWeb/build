Tea.context(function () {
	this.isRunning = false;
	this.isFinished = false;
	this.response = "";
	this.error = "";

	this.submitBefore = function () {
		this.isRunning = true;
		this.isFinished = false;
		this.response = "";
		this.error = "";
	};

	this.submitSuccess = function (resp) {
		this.isRunning = false;
		this.isFinished = true;
		this.response = resp.data.response;
		this.error = resp.data.error;
	};

	this.submitFail = function (resp) {
		this.isRunning = false;
		this.isFinished = true;
		this.response = "";
		this.error = resp.errors[0].messages[0];
	};

	this.submitError = function () {
		this.isRunning = false;
		this.isFinished = true;
		this.response = "";
		this.error = "请求超时";
	};
});