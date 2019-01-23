charts.PieChart = function () {
	this.type = "pie";
	this.values = [];
	this.labels = [];
	this.colors = colors.ARRAY;
};

charts.PieChart.prototype = new charts.Chart();