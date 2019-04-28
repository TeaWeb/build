charts.StackBarChart = function () {
	this.type = "stackbar";
	this.values = [];
	this.labels = [];
	this.colors = colors.ARRAY;
};

charts.StackBarChart.prototype = new charts.Chart();