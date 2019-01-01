Tea.context(function () {
    this.location = null;

    this.showSpeical = false;
    this.pattern = "";
    this.type = 2;
    this.root = "";
    this.charset = "utf-8";
    this.indexes = [];
    this.on = true;
    this.isCaseInsensitive = false;
    this.isReverse = false;

    this.submitSuccess = function () {
        alert("保存成功");
        window.location = "/proxy/locations?server=" + this.server.filename;
    };

    this.$delay(function () {
        this.$find("form input[name='pattern']").focus();
    });

    /**
     * advanced settings
     */
    this.advancedOptionsVisible = false;

    this.showAdvancedOptions = function () {
        this.advancedOptionsVisible = !this.advancedOptionsVisible;
    };

    /**
     * special settings
     */
    this.specialSettingsVisible = this.showSpecial; // 从参数中获取

    this.showSpecialSettings = function () {
        this.specialSettingsVisible = !this.specialSettingsVisible;
    };

    /**
     * pattern
     */
    this.patternDescription = "";
    this.changePatternType = function (patternType) {
        this.patternDescription = this.patternTypes.$find(function (k, v) {
            return v.type == patternType;
        }).description;
    };
    this.changePatternType(this.type);

    /**
     * index
     */
    this.indexAdding = false;
    this.addingIndexName = "";

    this.addIndex = function () {
        this.indexAdding = true;
        this.$delay(function () {
            this.$find("form input[name='addingIndexName']").focus();
        });
    };

    this.confirmAddIndex = function () {
        this.addingIndexName = this.addingIndexName.trim();
        if (this.addingIndexName.length == 0) {
            alert("文件名不能为空");
            this.$find("form input[name='addingIndexName']").focus();
            return;
        }
        this.indexes.push(this.addingIndexName);
        this.cancelIndexAdding();
    };

    this.cancelIndexAdding = function () {
        this.indexAdding = !this.indexAdding;
        this.addingIndexName = "";
    };

    this.removeLocationIndex = function (index) {
        this.indexes.$remove(index);
    };
});