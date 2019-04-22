Tea.context(function () {
	this.hasInternalGroups = this.groups.$exist(function (k, v) {
		return v.isChecked;
	});
});