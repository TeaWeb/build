Tea.context(function () {
	this.addresses = [];
	if (this.server != null && this.server.http != null && this.server.http.listen != null) {
		this.addresses = this.server.http.listen;
	}
});