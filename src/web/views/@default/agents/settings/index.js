Tea.context(function () {
	if (this.agentSpeed <= 1000) {
		this.speedPercent = 100 - this.agentSpeed * 100 / 1000;
	} else {
		this.speedPercent = 1;
	}

	this.putOn = function () {
		this.$post(".on")
			.params({
				"agentId": this.agent.id
			})
			.refresh();
	};

	this.putOff = function () {
		this.$post(".off")
			.params({
				"agentId": this.agent.id
			})
			.refresh();
	};
});