Tea.context(function () {
	this.file = "";
	this.step = "file";
	this.fileData = "";
	this.groups = [];

	this.changeFile = function (v) {
		if (v.target.files.length == 0) {
			this.file = "";
		} else {
			this.file = v.target.files[0].name;
		}
	};

	this.fileSuccess = function (resp) {
		this.fileData = resp.data.data;
		this.groups = resp.data.groups;
		this.step = "groups";
	};

	/**
	 * 分组
	 */
	this.countGroups = 0;
	this.countSets = 0;

	this.goFile = function () {
		this.step = "file";
	};

	this.groupsSuccess = function (resp) {
		this.step = "finish";
		this.countGroups = resp.data.countGroups;
		this.countSets = resp.data.countSets;
	};
});