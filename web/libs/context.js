var context = {};

context.features = [];
context.server = null;
context.agent = null;
context.app = null;
context.item = null;
context.timeType = "";
context.timePast = "";
context.timeUnit = "";
context.dayFrom = "";
context.dayTo = "";

context.hasTimeRange = function () {
	return this.timeType == "past" || this.timeType == "range";
};