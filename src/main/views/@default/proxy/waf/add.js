Tea.context(function () {
	this.$delay(function () {
		this.$find("form input[name='name']").focus();
	});

	this.submitSuccess = function () {
		alert("添加成功");
		window.location = "/proxy/waf";
	};
});