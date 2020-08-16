Tea.context(function () {
	var that = this;

	this.selectedIndex = -1;

	this.$delay(function () {
		this.$find("form input[name='keyword']").focus();

		var that = this;

		document.addEventListener("keyup", function (e) {
			if (e == null || e.key == null || e.target == null) {
				return;
			}
			var keyString = e.key.toString();
			if (keyString == "Escape") {
				window.history.back();
				return;
			}
			if (keyString == "ArrowDown") {
				that.selectedIndex++;
				if (that.selectedIndex >= that.results.length) {
					that.selectedIndex = 0;
				}
				return;
			}
			if (keyString == "ArrowUp") {
				that.selectedIndex--;
				if (that.selectedIndex < 0) {
					that.selectedIndex = 0;
				}
				return;
			}
			if (keyString == "Enter") {
				if (that.selectedIndex > -1) {
					window.location = that.results[that.selectedIndex].link;
				}
				return;
			}
			if (e.target.tagName == "INPUT") {
				return;
			}
		})
	});

	this.keyword = "";
	this.results = [];

	var oldKeyword = teaweb.get("globalSearchKeyword");
	if (typeof (oldKeyword) == "string") {
		this.keyword = oldKeyword;
		if (this.keyword.length > 0) {
			this.$delay(function () {
				this.changeKeyword();
			});
		}
	}

	this.changeKeyword = function () {
		teaweb.set("globalSearchKeyword", this.keyword);
		this.selectedIndex = -1;
		if (this.keyword.length == 0) {
			this.results = [];
			return;
		}
		this.$post("$")
			.params({
				"keyword": this.keyword
			})
			.success(function (resp) {
				this.results = resp.data.results;
			});
	};

	this.clearKeyword = function () {
		this.keyword = "";
		this.results = [];
		teaweb.set("globalSearchKeyword", this.keyword);
		this.$find("form input[name='keyword']").focus();
	};
});