Tea.context(function () {
	this.wafId = this.server.wafId;

	this.formVisible = false;

	this.showForm = function () {
		this.formVisible = true;
	};

	this.cancelForm = function () {
		this.formVisible = false;
	};

	this.updateWAF = function () {
		this.formVisible = true;
	};

	this.submitSuccess = function () {
		alert("保存成功");
		window.location.reload();
	};
});