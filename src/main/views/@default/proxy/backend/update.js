Tea.context(function () {
    this.advancedOptionsVisible = false;

    if (this.server.requestGroups != null) {
    	var selectedRequestGroupIds = (this.backend.requestGroupIds == null) ? [] : this.backend.requestGroupIds;
		this.server.requestGroups.$each(function (k, v) {
			v.isChecked = selectedRequestGroupIds.$contains(v.id);
		});
	}

    this.showAdvancedOptions = function () {
        this.advancedOptionsVisible = !this.advancedOptionsVisible;
    };

    this.updateSuccess = function () {
        alert("保存成功");
        window.location = this.from;
    };
});