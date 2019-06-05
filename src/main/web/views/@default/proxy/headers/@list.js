Tea.context(function () {
	this.loaded = false;

	this.requestHeaders = [];
	this.headers = [];
	this.ignoreHeaders = [];

	this.query = Tea.serialize(this.headerQuery);
	this.from = encodeURIComponent(window.location.toString());

	this.$delay(function () {
		this.loadData();
	});

	this.loadData = function () {
		this.$get("/proxy/headers/data")
			.params(this.headerQuery)
			.success(function (resp) {
				this.headers = resp.data.headers;
				this.requestHeaders = resp.data.requestHeaders;
				this.ignoreHeaders = resp.data.ignoreHeaders;
				this.loaded = true;
			});
	};

	this.deleteIgnoreHeader = function (header) {
		if (!window.confirm("确定要删除此Header吗？")) {
			return;
		}
		var query = this.headerQuery;
		query["name"] = header;
		this.$post("/proxy/headers/deleteIgnore")
			.params(query)
			.success(function () {
				window.location.reload();
			});
	};

	this.deleteHeader = function (header) {
		if (!window.confirm("确定要删除此Header吗？")) {
			return;
		}
		var query = this.headerQuery;
		query["headerId"] = header.id;
		this.$post("/proxy/headers/delete")
			.params(query)
			.success(function () {
				window.location.reload();
			});
	};

	this.deleteRequestHeader = function (header) {
		if (!window.confirm("确定要删除此请求Header吗？")) {
			return;
		}
		var query = this.headerQuery;
		query["headerId"] = header.id;
		this.$post("/proxy/headers/deleteRequestHeader")
			.params(query)
			.success(function () {
				window.location.reload();
			});
	};
});