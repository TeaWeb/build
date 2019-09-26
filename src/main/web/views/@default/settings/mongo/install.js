Tea.context(function () {
	this.downloadProgress = 0;
	this.installProgress = 0;
	this.startProgress = 0;
	this.checkProgress = 0;
	this.timeLeft = "";
	this.timeSeconds = 0;

	this.install = function () {
		this.isInstalling = true;
		this.$post(".install")
			.timeout(3600)
			.success(function () {
				window.location.reload();
			})
			.fail(function (response) {
				alert(response.message);
				this.isInstalling = false;
			});
	};

	this.progress = function () {
		this.$get(".installStatus")
			.success(function (response) {
				var status = response.data.status;
				var percent = response.data.percent;

				if (status == "download") {
					this.downloadProgress = percent;
					this.timeLeft = response.data.timeLeft;
					this.timeSeconds = response.data.timeSeconds;
				} else if (status == "install") {
					this.downloadProgress = 100;
					this.installProgress = percent;
				} else if (status == "start") {
					this.downloadProgress = 100;
					this.installProgress = 100;
					this.startProgress = percent;
				} else if (status == "check") {
					this.downloadProgress = 100;
					this.installProgress = 100;
					this.startProgress = 100;
					this.checkProgress = percent;

					if (percent == 100) {
						window.location.reload();
					}
				}
			})
			.done(function () {
				this.$delay(function () {
					this.progress();
				}, 1000);
			});
	};

	this.$delay(function () {
		this.progress();
	});
});