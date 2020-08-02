Tea.context(function () {
	this.formatVariables = function (code) {
		if (code == null) {
			return "";
		}
		code = code.replace(/\${(.+?)}/g, "<em>${<a>$1</a>}</em>");
		return code;
	};

	// syslog
	var that = this;
	this.priorityName = function (priority) {
		for (var i = 0; i < that.syslogPriorities.length; i++) {
			if (that.syslogPriorities[i].value == priority) {
				return that.syslogPriorities[i].name;
			}
		}
		return "";
	};
});