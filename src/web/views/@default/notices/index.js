Tea.context(function () {
	this.setRead = function (scope, noticeIds, msg) {
		if (msg != null) {
			if (!window.confirm(msg)) {
				return;
			}
		}
		this.$post("/notices/setRead")
			.params({
				"scope": scope,
				"noticeIds": (noticeIds != null) ? noticeIds : this.notices.$map(function (k, v) {
					return v.id;
				})
			})
			.success(function () {
				if (scope == "page") {
					window.location.reload();
				} else {
					window.location = "/notices";
				}
			});
	};

	this.reloadPage = function () {
		window.location.reload();
	};

	/**
	 * 声音
	 */
	this.testNoticeAudio = function () {
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
	};

	this.changeSoundOn = function () {
		this.$delay(function () {
			if (this.soundOn) {
				this.testNoticeAudio();
			}
			this.$post(".sound")
				.params({
					"on": this.soundOn ? 1 : 0
				});
		}, 100);
	};
});