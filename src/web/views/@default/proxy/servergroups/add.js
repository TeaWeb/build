Tea.context(function () {
	this.success = function () {
		teaweb.success("保存成功", function () {
			window.location = "/proxy/servergroups";
		});
	};
});