Tea.context(function () {
	this.exportWaf = function () {
		var groupIds = [];
		this.$find(".groups-box input[type='checkbox']").each(function (k, v) {
			if (v.checked) {
				groupIds.push(v.value);
			}
		});
		if (groupIds.length == 0) {
			alert("至少要选择一个规则分组");
			return;
		}
		window.location = "/proxy/waf/export?wafId=" + this.config.id + "&export=1&groupIds=" + groupIds.join(",");
	};
});