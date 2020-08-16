var charts = {};

charts.Chart = function () {
	this.html = "";
	this.options = {
		"name": "",
		"columns": 1,
		"events": []
	};

	this.menus = [];

	this.addMenu = function () {
		var menu = new charts.Menu(this);
		this.menus.push(menu);
		return menu;
	};

	this.render = function () {
		//stub
		callChartRender(this);
	};
};