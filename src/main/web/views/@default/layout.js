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
	var firstLoad = true;
	this.renewNoticeBadge = function () {
		this.$get("/notices/badge")
			.success(function (resp) {
				if (!firstLoad && resp.data.soundOn && resp.data.count > this.countNoticesBadge) {
					// play audio
					var audioBox = document.createElement("AUDIO");
					audioBox.setAttribute("control", "");
					audioBox.setAttribute("autoplay", "");
					audioBox.innerHTML = "<source src=\"/audios/notice.ogg\" type=\"audio/ogg\"/>";
					document.body.appendChild(audioBox);
					audioBox.play().then(function () {
						setTimeout(function () {
							document.body.removeChild(audioBox);
						}, 2000);
					}).catch(function (e) {
						console.log(e.message);
					});
				}
				firstLoad = false;
				this.countNoticesBadge = resp.data.count;
				if (this.countNoticesBadge > 0) {
					document.title = "(" + this.countNoticesBadge + "通知)" + documentTitle;
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

	/**
	 * 底部伸展框
	 */
	this.footerOuterVisible = false;

	this.showQQGroupQrcode = function () {
		this.footerOuterVisible = !this.footerOuterVisible;
	};

	/**
	 * 左侧子菜单
	 */
	this.showSubMenu = function (menu) {
		if (menu.alwaysActive) {
			return;
		}
		if (this.teaSubMenus.menus != null && this.teaSubMenus.menus.length > 0) {
			this.teaSubMenus.menus.$each(function (k, v) {
				if (menu.id == v.id) {
					return;
				}
				v.isActive = false;
			});
		}
		menu.isActive = !menu.isActive;
	};

	this.$delay(function () {
		var activeItem = this.$find(".main .sub-menu .item.active");
		if (activeItem.length > 0) {
			activeItem.focus();
		}
	}, 0);
});