charts.PieChart = function () {
	this.type = "pie";
	this.values = [];
	this.labels = [];
	this.colors = [];
};

charts.PieChart.prototype = new charts.Chart();