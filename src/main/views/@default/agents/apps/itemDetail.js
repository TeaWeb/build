Tea.context(function () {
	var that = this;

	this.from = encodeURIComponent(window.location.toString());

	if (this.item.thresholds != null) {
		this.item.thresholds.$each(function (k, v) {
			v.levelName = that.noticeLevels.$find(function (k, v1) {
				return v.noticeLevel == v1.code;
			}).name;
		});
	}
});