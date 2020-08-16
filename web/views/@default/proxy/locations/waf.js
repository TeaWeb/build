Tea.context(function () {
	this.$delay(function () {
		this.$find("#location-waf-menu").focus();
	});

	this.wafId = this.location.wafId;
	if (this.wafId.length == 0 && this.location.wafOn) {
		this.wafId = "none";
	}

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