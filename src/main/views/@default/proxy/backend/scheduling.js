Tea.context(function () {
    var that = this;

    this.type = this.scheduling.code;
    this.schedulingTypeDescription = null;

    this.changeSchedulingType = function () {
        this.schedulingTypeDescription = this.schedulingTypes.$find(function (k, v) {
            return v.code == that.type;
        }).description;
    };
    this.changeSchedulingType();

    // hash
    this.hashKey = "";
    this.hashVar = "";
    if (this.scheduling.code == "hash") {
        this.hashKey = this.scheduling.options.key;
    } else {
        this.hashKey = "${remoteAddr}";
    }

    this.changeHashVar = function () {
        if (this.hashVar.length > 0) {
            this.hashKey = this.hashVar;
        }
    };

    // sticky
    if (this.scheduling.code == "sticky") {
        this.stickyType = this.scheduling.options.type;
        this.stickyParam = this.scheduling.options.param;
    } else {
        this.stickyType = "cookie";
        this.stickyParam = "Backend";
    }

    this.saveSuccess = function () {
        alert("保存成功");
        window.location = this.from;
    };
});