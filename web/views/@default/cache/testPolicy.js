Tea.context(function () {
	this.key = "my-key";

	this.$delay(function () {
		this.$find("textarea[name='value']").focus();

		this.$watch("key", function () {
			this.writeResult = "";
			this.readResult = "";
		});
	});

	this.writeResult = "";
	this.writeOk = false;
	this.writeSuccess = function (resp) {
		this.writeResult = resp.data.result;
		this.writeOk = (resp.code == 200);
	};

	this.readResult = "";
	this.readOk = false;
	this.readSuccess = function (resp) {
		this.readResult = resp.data.result;
		this.readOk = (resp.code == 200);
	};
});