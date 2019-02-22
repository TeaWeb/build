charts.Clock = function () {
	this.type = "clock";
	this.timestamp = new Date().getTime() / 1000;
};

charts.Clock.prototype = new charts.Chart();
