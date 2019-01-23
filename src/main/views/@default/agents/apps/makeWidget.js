Tea.context(function () {
    this.$delay(function () {
        this.$find("input[name='name']").focus();
        this.loadEditor();
    });

    this.widgetName = "";
    this.widgetCode = "";
    this.widgetAuthor = "";
    this.widgetVersion = "";
    this.widgetDescription = "";
    this.chartCode = "";

    this.submitSuccess = function () {
    	alert("保存成功");
    	window.location = "/agents/apps/addWidget?agentId=" + this.agentId + "&appId=" + this.app.id;
	};

    /**
     * 参数
     */
    this.params = [];
    this.paramAdding = false;
    this.addingParamName = "";
    this.addingParamDescription = "";
    this.addingParamDefault = "";
    this.addingParamCode = "";

    this.addParam = function () {
        this.paramAdding = !this.paramAdding;
        this.addingParamName = "";
        this.addingParamDescription = "";
        this.addingParamDefault = "";
        this.addingParamCode = "";

        if (this.paramAdding) {
            this.$delay(function () {
                this.$find("form input[name='addingParamName']").focus();
            });
        }
    };

    this.confirmAddParam = function () {
        if (this.addingParamName.length == 0) {
            alert("请输入参数名");
            return;
        }
        if (this.addingParamCode.length == 0) {
            alert("请输入参数代号");
            return;
        }
        this.params.push({
            "name": this.addingParamName,
            "description": this.addingParamDescription,
            "default": this.addingParamDefault,
            "code": this.addingParamCode,
			"value": this.addingParamDefault // 当前的Value
        });
        this.paramAdding = false;
    };

    this.removeParam = function (index) {
        this.params.$remove(index);
    };

    /**
     * 编辑器
     */
    this.loadEditor = function () {
        var editor = CodeMirror.fromTextArea(document.getElementById("code-editor"), {
            theme: "idea",
            lineNumbers: true,
            value: "",
            readOnly: false,
            showCursorWhenSelecting: true,
            height: "auto",
            //scrollbarStyle: null,
            viewportMargin: Infinity,
            lineWrapping: true,
            highlightFormatting: false,
            indentUnit: 4,
            indentWithTabs: true
        });
        editor.setValue( '\
\n{\
\n	var chart = new charts.LineChart();\
\n	chart.options.name = "测试图表";\
\n	chart.options.columns = 2;\
\n\
\n	var line = new charts.Line();\
\n	line.values = new values.Query().latestValues(10);\
\n	chart.addLine(line);\
\n\
\n	chart.render();\
\n}\
\n');
        editor.save();
        this.chartCode = editor.getValue();

        var info = CodeMirror.findModeByMIME("text/javascript");
        if (info != null) {
            editor.setOption("mode", info.mode);
            CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
            CodeMirror.autoLoadMode(editor, info.mode);
        }

        var that = this;
        editor.on("change", function () {
            editor.save();
            that.chartCode = editor.getValue();
        });
    };

	/**
	 * 测试
	 */
	this.isTesting = false;
	this.charts = [];

	this.test = function () {
		this.isTesting = !this.isTesting;
		if (this.isTesting) {
			this.$delay(function () {
				window.scroll(0, 1000);
			});
		}

		this.charts = [];
	};

	this.testSuccess = function (resp) {
		// output
		resp.data.output.$each(function (k, v) {
			console.log("[widget]" + v);
		});

		// charts
		this.charts = resp.data.charts;
		new ChartRender(this.charts);
		this.$delay(function () {
			window.scroll(0, 1000);
		});
	};


	this.testingValueAdding = false;
	this.testingValues = [];

	this.initTestingValue = function () {
		this.testingValue = "";
		this.testingHour = 0;
		this.testingMinute = 0;
		this.testingSecond = 0;
	};

	this.initTestingValue();

	this.addTestingValue = function () {
		this.testingValueAdding = !this.testingValueAdding;
		if (this.testingValueAdding) {
			this.$delay(function () {
				this.$find("form textarea[name='testingNewValue']").focus();
			});
		}
	};

	this.confirmAddTestingValue = function () {
		var value = this.testingValue;
		try {
			value = JSON.parse(this.testingValue);
		} catch(e) {

		}
		this.testingValues.push({
			"value": value,
			"valueString": JSON.stringify(value),
			"year": this.testingYear,
			"month": this.testingMonth,
			"day": this.testingDay,
			"hour": this.testingHour,
			"minute": this.testingMinute,
			"second": this.testingSecond
		});
		this.testingValueAdding = false;
		this.initTestingValue();
	};

	this.removeTestingValue = function (index) {
		this.testingValues.$remove(index);
	};
});