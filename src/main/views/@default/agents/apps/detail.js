Tea.context(function () {
	this.from = encodeURIComponent(window.location.toString());

	this.putOn = function () {
		this.$post(".on")
			.params({
				"agentId": this.agentId,
				"appId": this.app.id
			})
			.refresh();
	};

	this.putOff = function () {
		this.$post(".off")
			.params({
				"agentId": this.agentId,
				"appId": this.app.id
			})
			.refresh();
	};
});