charts.TableChart = function () {
	this.type = "table";
	this.rows = [];
	this.widths = [];
	this.menus = [];

	this.addRow = function () {
		var cols = [];
		for (var i = 0; i < arguments.length; i++) {
			cols.push(arguments[i]);
		}
		this.rows.push(cols);
	};

	this.setWidth = function (index, width) {
		if (this.widths.length > index) {
			this.widths[index] = width;
		} else {
			for (var i = this.widths.length; i < index; i++) {
				this.widths.push(null);
			}
			this.widths.push(width);
		}
	};
};

charts.TableChart.prototype = new charts.Chart();