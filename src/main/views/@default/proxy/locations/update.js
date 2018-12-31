Tea.context(function () {
    this.submitSuccess = function () {
        alert("保存成功");
        window.location = this.from;
    };

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
    this.changePatternType(this.location.type);

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
        this.location.index.push(this.addingIndexName);
        this.cancelIndexAdding();
    };

    this.cancelIndexAdding = function () {
        this.indexAdding = !this.indexAdding;
        this.addingIndexName = "";
    };

    this.removeLocationIndex = function (index) {
        this.location.index.$remove(index);
    };
});