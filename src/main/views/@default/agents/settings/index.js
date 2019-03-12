Tea.context(function () {
	if (this.agentSpeed <= 1000) {
		this.speedPercent = 100 - this.agentSpeed * 100 / 1000;
	} else {
		this.speedPercent = 0;
	}
});