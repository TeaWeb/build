Tea.context(function () {
	this.currentDBType = this.dbType;

	this.goNext = function () {
		window.location = "/settings/" + this.dbType + "/update?action=switchType";
	};
});