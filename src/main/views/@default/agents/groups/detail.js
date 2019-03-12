Tea.context(function () {
	this.keyword = "";
	this.filterOnline = "-1";
	var allAgents = this.agents;
	this.hasAgents = allAgents.length > 0;

	this.changeKeyword = function () {
		this.filter();
	};

	this.resetKeyword = function () {
		this.keyword = "";
		this.filterOnline = "-1";
		this.$delay(function () {
			this.changeKeyword();
		});
	};

	this.changeOnlineFilter = function () {
		this.filter();
	};

	this.filter = function () {
		if (this.keyword.length == 0 && this.filterOnline == "-1") {
			this.agents = allAgents;
			return;
		}
		var keyword = this.keyword;
		var filterOnline = this.filterOnline;
		this.agents = allAgents.$filter(function (k, v) {
			if (keyword.length > 0) {
				if (!teaweb.match(v.name + " " + v.host, keyword)) {
					return false;
				}
			}
			if (filterOnline != "-1") {
				if (filterOnline == "0" && v.isWaiting) {
					return false;
				}
				if (filterOnline == "1" && !v.isWaiting) {
					return false;
				}
			}
			return true;
		});
	};
});