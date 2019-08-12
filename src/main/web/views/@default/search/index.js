Tea.context(function () {
	var that = this;

	this.$delay(function () {
		this.$find("form input[name='keyword']").focus();

		document.addEventListener("keyup", function (e) {
			if (e == null || e.key == null || e.target == null) {
				return;
			}
			if (e.key.toString() == "Escape") {
				window.history.back();
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