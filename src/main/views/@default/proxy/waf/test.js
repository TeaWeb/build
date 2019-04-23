Tea.context(function () {
	this.step = "data";
	this.result = [];
	this.matched = false;
	this.action = "";
	this.setName = "";

	this.$delay(function () {
		var textareaList = this.$find("form textarea");
		if (textareaList.length > 0) {
			textareaList[0].focus();
		}
	});

	this.goStep = function (step) {
		if (step == "result") {
			this.$find("#data-form button[type='submit']")[0].click();
			return;
		}

		this.step = step;

		if (step == "data") {
			this.$delay(function () {
				var textareaList = this.$find("form textarea");
				if (textareaList.length > 0) {
					textareaList[0].focus();
				}
			});
		}
	};

	this.submitSuccess = function (resp) {
		this.result = resp.data.result;
		this.matched = resp.data.matched;
		this.action = resp.data.action.toUpperCase();
		this.setName = resp.data.setName;
		this.step = "result";
	};
});