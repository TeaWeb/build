charts.ProgressBarChart = function () {
	this.type = "progress";
	this.value = 0;
	this.menus = [];
};

charts.ProgressBarChart.prototype = new charts.Chart();