Tea.context(function () {
	this.hasInternalGroups = this.groups.$exist(function (k, v) {
		return v.isChecked;
	});

	this.updatesVisible = false;

	this.showUpdates = function () {
		this.updatesVisible = !this.updatesVisible;
	};

	this.mergeTemplate = function () {
		this.$post(".mergeTemplate")
			.params({
				"wafId": this.config.id
			})
			.refresh();
	};
});