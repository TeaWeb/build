var charts = {};

charts.Chart = function () {
	this.options = {
		"name": "",
		"columns": 1
	};

	this.render = function () {
		//stub
		renderChart(this);
	};
};