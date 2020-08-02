charts.PieChart = function () {
	this.type = "pie";
	this.values = [];
	this.labels = [];
	this.colors = colors.ARRAY;
	this.menus = [];
};

charts.PieChart.prototype = new charts.Chart();