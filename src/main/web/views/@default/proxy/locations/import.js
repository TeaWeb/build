Tea.context(function () {
	this.file = "";
	this.locationPattern = "";

	this.changeFile = function (v) {
		var that = this;
		this.locationPattern = "";
		if (v.target.files.length == 0) {
			this.file = "";
		} else {
			this.file = v.target.files[0].name;

			if (typeof (FileReader) != "undefined") {
				var reader = new FileReader();
				reader.onload = (function (reader) {
					return function () {
						var content = reader.result;
						if (typeof (content) == "string") {
							content.split("\n").$each(function (k, v) {
								if (v.startsWith("pattern: ")) {
									that.locationPattern = v.substring("pattern: ".length);
								}
							});
						}
					}
				})(reader);

				reader.readAsText(v.target.files[0]);
			}
		}
	};

	this.fileSuccess = function (resp) {
		this.fileData = resp.data.data;
		this.groups = resp.data.groups;
		this.step = "groups";
	};

	this.submitSuccess = function () {
		alert("导入成功");
		window.location.reload();
	};
});