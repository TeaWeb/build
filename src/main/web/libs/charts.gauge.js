charts.GaugeChart = function () {
	this.type = "gauge";
	this.value = 0;
	this.label = "";
	this.min = 0;
	this.max = 0;
	this.unit = "";
	this.menus = [];
};

charts.GaugeChart.prototype = new charts.Chart();
