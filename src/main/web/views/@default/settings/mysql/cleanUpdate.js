Tea.context(function () {
	if (this.accessLog == null) {
		this.accessLog = {
			"cleanHour": "",
			"keepDays": ""
		};
	}

	this.$delay(function () {
		this.$find("form input[name='accessLogCleanHour']").focus();
	});

	this.saveSuccess = function () {
		alert("保存成功");
		window.location = "/settings/mysql/clean";
	};
});