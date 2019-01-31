Tea.context(function () {
	/**
	 * 测试MongoDB连接
	 */
	this.mongoFailed = false;

	this.testMongo = function () {
		this.$get("/mongo/test")
			.fail(function () {
				this.mongoFailed = true;
			});
	};

	/**
	 * 计算未读消息数
	 */
	this.countNoticesBadge = 0;

	this.$delay(function () {
		this.renewNoticeBadge();
	});

	var documentTitle = document.title;
	this.renewNoticeBadge = function () {
		this.$get("/notices/badge")
			.success(function (resp) {
				this.countNoticesBadge = resp.data.count;
				if (this.countNoticesBadge > 0) {
					document.title = documentTitle + "(" + this.countNoticesBadge + "通知)";
				} else {
					document.title = documentTitle;
				}
			})
			.done(function () {
				this.$delay(function () {
					this.renewNoticeBadge();
				}, 60000);
			});
	};
});