Tea.context(function () {
	if (typeof (teaweb) != "undefined") {
		this.teaweb = teaweb;
	}

	this.$delay(function () {
		var focusInput = this.$refs.focusInput;
		if (focusInput != null) {
			focusInput.focus();
		}
	});

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
		this.loadGlobalEvents();
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

	/**
	 * 菜单辅助函数
	 */
	this.urlPrefix = function () {
		for (var i = 0; i < arguments.length; i++) {
			var b = window.location.pathname.startsWith(arguments[i]);
			if (b) {
				return true;
			}
		}
		return false;
	};

	/**
	 * 当前URL
	 */
	this.globalURL = encodeURIComponent(window.location.toString());

	/**
	 * ss进入搜索
	 */
	var lastSSTime = null;
	this.loadGlobalEvents = function () {
		var that = this;
		document.addEventListener("keyup", function (e) {
			if (e.key == null || e.target == null) {
				return;
			}
			if (["INPUT", "SELECT", "TEXTAREA", "BUTTON"].$contains(e.target.tagName)) {
				return;
			}
			if (e.key.toString() == "s") {
				if (lastSSTime == null) {
					lastSSTime = new Date();
					return;
				}
				var delta = new Date().getTime() - lastSSTime.getTime();
				if (delta < 500) {
					window.location = "/search?from=" + encodeURIComponent(window.location.toString());
					return;
				}
				lastSSTime = new Date();
			}

			if (e.key.toString() == "Escape") {
				that.closeModal();
			}
		});
	};

	/**
	 * 关闭Modal
	 */
	this.closeModal = function () {
		this.$find(".modal").each(function (k, v) {
			v.className = "modal";
		});
	};

	this.showModal = function (modalId) {
		var modal = document.getElementById("chart-setting-modal");
		modal.className = "modal visible";
	};
});

window.NotifyPopup = function (resp) {
	window.parent.teaweb.popupFinish(resp);
};