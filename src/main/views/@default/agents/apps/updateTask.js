Tea.context(function () {
    this.currentLocation = window.location.toString();
    this.isBooting = this.task.isBooting;
    this.isManual = this.task.isManual;
    this.minYear = new Date().getFullYear();

    this.$delay(function () {
        this.$find("form input[name='name']").focus();
        this.loadEditor();
        this.loadSchedules();
    });

    this.submitSuccess = function (response) {
        alert("保存成功");
        window.location = this.from;
    };

    this.submitBefore = function() {
        // 是否保存了定时任务
        if (this.scheduleAdding) {
            alert("请先确定定时任务");
            return false;
        }
    };

    /**
     * 更多选项
     */
    this.advancedOptionsVisible = false;

    this.showAdvancedOptions = function () {
        this.advancedOptionsVisible = !this.advancedOptionsVisible;
    };

    /**
     * 编辑器
     */
    this.loadEditor = function () {
        var editor = CodeMirror.fromTextArea(document.getElementById("editor-box"), {
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
        editor.setValue(this.task.script);
        editor.save();

        var info = CodeMirror.findModeByMIME("text/x-sh");
        if (info != null) {
            editor.setOption("mode", info.mode);
            CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
            CodeMirror.autoLoadMode(editor, info.mode);
        }

        editor.on("change", function () {
            editor.save();
        });
    };

    /**
     * 环境变量
     */
    this.env = this.task.env;
    this.envAdding = false;
    this.envAddingName = "";
    this.envAddingValue = "";

    this.addEnv = function () {
        this.envAdding = !this.envAdding;
        this.$delay(function () {
            this.$find("form input[name='envAddingName']").focus();
        });
    };

    this.confirmAddEnv = function () {
        if (this.envAddingName.length == 0) {
            alert("请输入变量名");
            this.$find("form input[name='envAddingName']").focus();
        }
        this.env.push({
            "name": this.envAddingName,
            "value": this.envAddingValue
        });
        this.envAdding = false;
        this.envAddingName = "";
        this.envAddingValue = "";
    };

    this.removeEnv = function (index) {
        this.env.$remove(index);
    };

    this.cancelEnv = function () {
        this.envAdding = false;
    };

    /**
     * 时间范围
     */
    this.scheduleTab = "quick";

   this.initSchedule = function () {
       this.quickSecond = 1;
       this.quickMinute = 1;
       this.quickHour = 1;
       this.quickDayHour = 0;
       this.quickDayMinute = 0;
       this.quickDaySecond = 0;

       this.secondLimit = false;
       this.secondTab = "every";
       this.secondEvery = false;
       this.secondPointAdding = false;
       this.secondPointAddingName = "";
       this.secondPoints = [];
       this.secondStepAdding = false;
       this.secondStepAddingStep = {};
       this.secondSteps = [];

       this.minuteLimit = false;
       this.minuteTab = "every";
       this.minuteEvery = false;
       this.minutePointAdding = false;
       this.minutePointAddingName = "";
       this.minutePoints = [];
       this.minuteStepAdding = false;
       this.minuteStepAddingStep = {};
       this.minuteSteps = [];

       this.hourLimit = false;
       this.hourTab = "every";
       this.hourEvery = false;
       this.hourPointAdding = false;
       this.hourPointAddingName = "";
       this.hourPoints = [];
       this.hourStepAdding = false;
       this.hourStepAddingStep = {};
       this.hourSteps = [];

       this.dayLimit = false;
       this.dayTab = "point";
       this.dayPointAdding = false;
       this.dayPointAddingName = "";
       this.dayPoints = [];
       this.dayStepAdding = false;
       this.dayStepAddingStep = {};
       this.daySteps = [];

       this.monthLimit = false;
       this.monthTab = "point";
       this.monthPointAdding = false;
       this.monthPointAddingName = "";
       this.monthPoints = [];
       this.monthStepAdding = false;
       this.monthStepAddingStep = {};
       this.monthSteps = [];

       this.yearLimit = false;
       this.yearTab = "point";
       this.yearPointAdding = false;
       this.yearPointAddingName = "";
       this.yearPoints = [];
       this.yearStepAdding = false;
       this.yearStepAddingStep = {};
       this.yearSteps = [];

       this.weekDayLimit = false;
       this.weekDayTab = "point";
       this.weekDayPointAdding = false;
       this.weekDayPointAddingName = "";
       this.weekDayPoints = [];
       this.weekDayStepAdding = false;
       this.weekDayStepAddingStep = {};
       this.weekDaySteps = [];
   };
   this.initSchedule();

    this.selectScheduleTab = function (tab) {
        this.scheduleTab = tab;
    };

    this.addQuickTime = function (type, step, from, to) {
        if (!/^\d+$/.test(step.toString())) {
            alert("请输入一个整数数字");
            return;
        }

        step = parseInt(step, 10);
        if (from > step || to < step) {
            alert("请输入一个在" + from + "和" + to + "之间的整数数字");
            return;
        }

        this[type + "Steps"].push({
            "from": 0,
            "to": to,
            "step": step
        });
        if (type == "minute") {
            this["secondPoints"].push(0);
        } else if (type == "hour") {
            this["secondPoints"].push(0);
            this["minutePoints"].push(0);
        }
        this.confirmAddSchedule();
    };

    this.addQuickDayTime = function () {
        if (!/^\d+$/.test(this.quickDayHour.toString())) {
            alert("小时请输入一个整数数字");
            return;
        }
        this.quickDayHour = parseInt(this.quickDayHour);
        if (this.quickDayHour < 0 || this.quickDayHour > 23) {
            alert("小时请输入一个在0和23之间的整数数字");
            return;
        }

        if (!/^\d+$/.test(this.quickDayMinute.toString())) {
            alert("分钟请输入一个整数数字");
            return;
        }
        this.quickDayMinute = parseInt(this.quickDayMinute);
        if (this.quickDayMinute < 0 || this.quickDayMinute > 60) {
            alert("分钟请输入一个在0和59之间的整数数字");
            return;
        }

        if (!/^\d+$/.test(this.quickDaySecond.toString())) {
            alert("秒钟请输入一个整数数字");
            return;
        }
        this.quickDaySecond = parseInt(this.quickDaySecond);
        if (this.quickDaySecond < 0 || this.quickDaySecond > 60) {
            alert("秒钟请输入一个在0和59之间的整数数字");
            return;
        }
        this["hourPoints"].push(this.quickDayHour);
        this["minutePoints"].push(this.quickDayMinute);
        this["secondPoints"].push(this.quickDaySecond);
        this.confirmAddSchedule();
    };

    this.addPoint = function (type) {
        this[type + 'PointAdding'] = true;
        this.$delay(function () {
            this.$find("input[name='" + type + "PointAddingName']").focus();
        });
    };

    this.cancelPoint = function (type) {
        this[type + 'PointAdding'] = false;
        this[type + "PointAddingName"] = "";
    };

    this.confirmAddPoint = function (type, from, to) {
        var value = this[type + "PointAddingName"];
        if (value.length == 0) {
            alert("不能为空");
            this.$find("input[name='" + type + "PointAddingName']").focus();
            return;
        }

        if (!/^\d+$/.test(value)) {
            alert("请输入一个数字");
            this.$find("input[name='" + type + "PointAddingName']").focus();
            return;
        }

        value = parseInt(value, 10);
        if (value < from || value > to) {
            alert("请输入一个在" + from + "和" + to + "之间的数字");
            this.$find("input[name='" + type + "PointAddingName']").focus();
            return;
        }

        this[type + "Points"].push(value);
        this.cancelPoint(type);
    };

    this.removePoint = function (type, index) {
        this[type + "Points"].$remove(index);
    };

    this.addStep = function (type) {
        this[type + 'StepAddingStep'] = {
            "from": "",
            "to": "",
            "step": ""
        };
        this[type + 'StepAdding'] = true;
        this.$delay(function () {
            this.$find("input[name='" + type + "StepAddingFrom']").focus();
        });
    };

    this.cancelStep = function (type) {
        this[type + 'StepAdding'] = false;
    };

    this.confirmAddStep = function (type, from, to, maxStep) {
        var stepFrom = this[type + "StepAddingStep"]["from"];
        if (stepFrom.length == 0) {
            alert("不能为空");
            this.$find("input[name='" + type + "StepAddingFrom']").focus();
            return;
        }

        if (!/^\d+$/.test(stepFrom)) {
            alert("请输入一个整数数字");
            this.$find("input[name='" + type + "StepAddingFrom']").focus();
            return;
        }

        stepFrom = parseInt(stepFrom, 10);
        if (stepFrom < from || stepFrom > to) {
            alert("请输入一个在" + from + "和" + to + "之间的整数数字");
            this.$find("input[name='" + type + "StepAddingFrom']").focus();
            return;
        }

        var stepTo = this[type + "StepAddingStep"]["to"];
        if (stepTo.length == 0) {
            alert("不能为空");
            this.$find("input[name='" + type + "StepAddingTo']").focus();
            return;
        }

        if (!/^\d+$/.test(stepTo)) {
            alert("请输入一个整数数字");
            this.$find("input[name='" + type + "StepAddingTo']").focus();
            return;
        }

        stepTo = parseInt(stepTo, 10);
        if (stepTo < from || stepTo > to) {
            alert("请输入一个在" + from + "和" + to + "之间的整数数字");
            this.$find("input[name='" + type + "StepAddingTo']").focus();
            return;
        }

        if (stepTo < stepFrom) {
            alert("结束的数字不能小于开始的整数数字");
            this.$find("input[name='" + type + "StepAddingTo']").focus();
            return;
        }

        var stepStep = this[type + "StepAddingStep"]["step"];
        if (stepStep.length == 0) {
            alert("不能为空");
            this.$find("input[name='" + type + "StepAddingStep']").focus();
            return;
        }

        if (!/^\d+$/.test(stepStep)) {
            alert("请输入一个整数数字");
            this.$find("input[name='" + type + "StepAddingStep']").focus();
            return;
        }

        stepStep = parseInt(stepStep, 10);
        if (stepStep <= 0) {
            alert("必须大于0");
            this.$find("input[name='" + type + "StepAddingStep']").focus();
            return;
        }

        if (stepStep > maxStep) {
            alert("不能大于" + maxStep);
            this.$find("input[name='" + type + "StepAddingStep']").focus();
            return;
        }

        this[type + "Steps"].push({
            "from": stepFrom,
            "to": stepTo,
            "step": stepStep
        });
        this.cancelStep(type);
    };

    this.removeStep = function (type, index) {
        this[type + "Steps"].$remove(index);
    };

    this.scheduleAdding = false;
    this.schedulesJSON = "";
    this.schedules = [];

    this.addSchedule = function () {
        this.scheduleAdding = !this.scheduleAdding;
    };

    this.confirmAddSchedule = function () {
        var summaryStrings = [];
        var secondStrings = [];
        if (this.secondEvery || (this.secondPoints.length == 0 && this.secondSteps.length == 0)) {
            secondStrings.push("每秒钟");
        }
        if (this.secondPoints.length > 0) {
            secondStrings.$pushAll(this.secondPoints.$map(function (k, v) {
                return v + "秒";
            }));
        }
        if (this.secondSteps.length > 0) {
            secondStrings.$pushAll(this.secondSteps.$map(function (k, v) {
                return v["from"] + "秒" + "-" + v["to"] + "秒/每" + v["step"] + "秒";
            }));
        }
        if (secondStrings.length > 0) {
            summaryStrings.push("[" + secondStrings.join("，") + "]");
        }

        var minuteStrings = [];
        if (this.minuteEvery || (this.minutePoints.length == 0 && this.minuteSteps.length == 0)) {
            minuteStrings.push("每分钟");
        }
        if (this.minutePoints.length > 0) {
            minuteStrings.$pushAll(this.minutePoints.$map(function (k, v) {
                return v + "分";
            }));
        }
        if (this.minuteSteps.length > 0) {
            minuteStrings.$pushAll(this.minuteSteps.$map(function (k, v) {
                return v["from"] + "分钟" + "-" + v["to"] + "分钟/每" + v["step"] + "分钟";
            }));
        }
        if (minuteStrings.length > 0) {
            summaryStrings.push("[" + minuteStrings.join("，") + "]");
        }

        var hourStrings = [];
        if (this.hourEvery || (this.hourPoints.length == 0 && this.hourSteps.length == 0)) {
            hourStrings.push("每小时");
        }
        if (this.hourPoints.length > 0) {
            hourStrings.$pushAll(this.hourPoints.$map(function (k, v) {
                return v + "小时";
            }));
        }
        if (this.hourSteps.length > 0) {
            hourStrings.$pushAll(this.hourSteps.$map(function (k, v) {
                return v["from"] + "小时" + "-" + v["to"] + "小时/每" + v["step"] + "小时";
            }));
        }
        if (hourStrings.length > 0) {
            summaryStrings.push("[" + hourStrings.join("，") + "]");
        }

        if (summaryStrings.length == 0) {
            alert("必须设置时分秒中的其中一个");
            return;
        }

        var dayStrings = [];
        if (this.dayPoints.length > 0) {
            dayStrings.$pushAll(this.dayPoints.$map(function (k, v) {
                return v + "日";
            }));
        }
        if (this.daySteps.length > 0) {
            dayStrings.$pushAll(this.daySteps.$map(function (k, v) {
                return v["from"] + "日" + "-" + v["to"] + "日/每" + v["step"] + "天";
            }));
        }
        if (dayStrings.length > 0) {
            summaryStrings.push("[" + dayStrings.join("，") + "]");
        }

        var monthStrings = [];
        if (this.monthPoints.length > 0) {
            monthStrings.$pushAll(this.monthPoints.$map(function (k, v) {
                return v + "月";
            }));
        }
        if (this.monthSteps.length > 0) {
            monthStrings.$pushAll(this.monthSteps.$map(function (k, v) {
                return v["from"] + "月" + "-" + v["to"] + "月/每" + v["step"] + "月";
            }));
        }
        if (monthStrings.length > 0) {
            summaryStrings.push("[" + monthStrings.join("，") + "]");
        }

        var yearStrings = [];
        if (this.yearPoints.length > 0) {
            yearStrings.$pushAll(this.yearPoints.$map(function (k, v) {
                return v + "年";
            }));
        }
        if (this.yearSteps.length > 0) {
            yearStrings.$pushAll(this.yearSteps.$map(function (k, v) {
                return v["from"] + "年" + "-" + v["to"] + "年/每" + v["step"] + "年";
            }));
        }
        if (yearStrings.length > 0) {
            summaryStrings.push("[" + yearStrings.join("，") + "]");
        }

        var weekDayStrings = [];
        if (this.weekDayPoints.length > 0) {
            weekDayStrings.$pushAll(this.weekDayPoints.$map(function (k, v) {
                return "周" + v;
            }));
        }
        if (this.weekDaySteps.length > 0) {
            weekDayStrings.$pushAll(this.weekDaySteps.$map(function (k, v) {
                return "周" + v["from"] + "-周" + v["to"] + "/每" + v["step"] + "天";
            }));
        }
        if (weekDayStrings.length > 0) {
            summaryStrings.push("[" + weekDayStrings.join("，") + "]");
        }
        if (summaryStrings.length == 0) {
            alert("请设置一个时间");
            return;
        }

        this.schedules.push({
            "second": {
                "every": this.secondEvery,
                "points": this.secondPoints,
                "steps": this.secondSteps
            },
            "minute": {
                "every": this.minuteEvery,
                "points": this.minutePoints,
                "steps": this.minuteSteps
            },
            "hour": {
                "every": this.hourEvery,
                "points": this.hourPoints,
                "steps": this.hourSteps
            },
            "day": {
                "points": this.dayPoints,
                "steps": this.daySteps
            },
            "month": {
                "points": this.monthPoints,
                "steps": this.monthSteps
            },
            "year": {
                "points": this.yearPoints,
                "steps": this.yearSteps
            },
            "weekDay": {
                "points": this.weekDayPoints,
                "steps": this.weekDaySteps
            },
            "summary": summaryStrings.join(" ")
        });
        this.schedulesJSON = JSON.stringify(this.schedules);
        this.scheduleAdding = false;
        this.initSchedule();
    };

    this.removeSchedule = function (index) {
        this.schedules.$remove(index);
        this.schedulesJSON = JSON.stringify(this.schedules);
    };

    this.dayOptionsVisible = false;

    this.showDayOptions = function () {
        this.dayOptionsVisible = !this.dayOptionsVisible;
    };

    this.loadSchedules = function () {
        var that = this;
        this.task.schedules.$each(function (k, v) {
            that.schedules.push({
                "second": {
                    "every": v.secondRanges.$find(function (k, v) {
                        return v.every;
                    }) != null,
                    "points": v.secondRanges.$map(function (k, v) {
                        if (v.value < 0) {
                            return Array.$nil;
                        }
                        return v.value;
                    }),
                    "steps": v.secondRanges.$map(function (k, v) {
                        if (v.from < 0) {
                            return Array.$nil;
                        }
                        return {
                            "from": v.from,
                            "to": v.to,
                            "step": v.step
                        };
                    })
                },
                "minute": {
                    "every": v.minuteRanges.$find(function (k, v) {
                        return v.every;
                    }) != null,
                    "points": v.minuteRanges.$map(function (k, v) {
                        if (v.value < 0) {
                            return Array.$nil;
                        }
                        return v.value;
                    }),
                    "steps": v.minuteRanges.$map(function (k, v) {
                        if (v.from < 0) {
                            return Array.$nil;
                        }
                        return {
                            "from": v.from,
                            "to": v.to,
                            "step": v.step
                        };
                    })
                },
                "hour": {
                    "every": v.hourRanges.$find(function (k, v) {
                        return v.every;
                    }) != null,
                    "points": v.hourRanges.$map(function (k, v) {
                        if (v.value < 0) {
                            return Array.$nil;
                        }
                        return v.value;
                    }),
                    "steps": v.hourRanges.$map(function (k, v) {
                        if (v.from < 0) {
                            return Array.$nil;
                        }
                        return {
                            "from": v.from,
                            "to": v.to,
                            "step": v.step
                        };
                    })
                },
                "day": {
                    "every": false,
                    "points": v.dayRanges.$map(function (k, v) {
                        if (v.value < 0) {
                            return Array.$nil;
                        }
                        return v.value;
                    }),
                    "steps": v.dayRanges.$map(function (k, v) {
                        if (v.from < 0) {
                            return Array.$nil;
                        }
                        return {
                            "from": v.from,
                            "to": v.to,
                            "step": v.step
                        };
                    })
                },
                "month": {
                    "every": false,
                    "points": v.monthRanges.$map(function (k, v) {
                        if (v.value < 0) {
                            return Array.$nil;
                        }
                        return v.value;
                    }),
                    "steps": v.monthRanges.$map(function (k, v) {
                        if (v.from < 0) {
                            return Array.$nil;
                        }
                        return {
                            "from": v.from,
                            "to": v.to,
                            "step": v.step
                        };
                    })
                },
                "year": {
                    "every": false,
                    "points": v.yearRanges.$map(function (k, v) {
                        if (v.value < 0) {
                            return Array.$nil;
                        }
                        return v.value;
                    }),
                    "steps": v.yearRanges.$map(function (k, v) {
                        if (v.from < 0) {
                            return Array.$nil;
                        }
                        return {
                            "from": v.from,
                            "to": v.to,
                            "step": v.step
                        };
                    })
                },
                "weekDay": {
                    "every": false,
                    "points": v.weekDayRanges.$map(function (k, v) {
                        if (v.value < 0) {
                            return Array.$nil;
                        }
                        return v.value;
                    }),
                    "steps": v.weekDayRanges.$map(function (k, v) {
                        if (v.from < 0) {
                            return Array.$nil;
                        }
                        return {
                            "from": v.from,
                            "to": v.to,
                            "step": v.step
                        };
                    })
                },
                "summary": v.summary
            });
        });
        this.schedulesJSON = JSON.stringify(this.schedules);
    };
});