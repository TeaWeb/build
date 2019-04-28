charts.Line = function () {
	this.name = "";
	this.values = [];
	this.isFilled = false;
	this.showItems = false;
	this.color = colors.BLUE;
	this.smooth = false;
};

charts.LineChart = function () {
	this.type = "line";
	this.lines = [];
	this.labels = [];
	this.max = 0;
	this.xShowTick = true;
	this.xTickCount = 0;
	this.yShowTick = true;

	this.addLine = function (line) {
		this.lines.push(line);
	};

	this.addLines = function (lines) {
		this.lines.$pushAll(lines);
	};
};

charts.LineChart.prototype = new charts.Chart();