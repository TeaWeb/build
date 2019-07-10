Tea.context(function () {
	this.formatVariables = function (code) {
		if (code == null) {
			return "";
		}
		code = code.replace(/\${(.+?)}/g, "<em>${<a>$1</a>}</em>");
		return code;
	};
});