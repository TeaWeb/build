Tea.context(function () {
	this.addresses = (this.server != null) ? this.server.http.listen.join("\n") : [];
});