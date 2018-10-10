Tea.context(function () {
    this.locationAdding = false;
    this.pattern = "";
    this.typeId = 1;
    this.isReverse = false;
    this.isCaseInsensitive = false;

   this.addLocation = function () {
        this.locationAdding = !this.locationAdding;
   };

   this.reverse = function () {
       this.isReverse = !this.isReverse;
   };

   this.switchCaseInsensitive = function () {
       this.isCaseInsensitive = !this.isCaseInsensitive;
   };

   this.locationSave = function () {
       this.$post("/proxy/locations/add")
           .params({
               "filename": this.filename,
               "pattern": this.pattern,
               "typeId": this.typeId,
               "reverse": this.isReverse ? 1 : 0,
               "caseInsensitive": this.isCaseInsensitive ? 1 : 0
           });
   };

   this.deleteLocation = function (index) {
       if (!window.confirm("确定要删除此路径配置吗？")) {
           return;
       }
       this.$post("/proxy/locations/delete")
           .params({
               "filename": this.filename,
               "index": index
           });
   };

   this.moveUp = function (index) {
       this.$post("/proxy/locations/moveUp")
           .params({
               "filename": this.filename,
               "index": index
           });
   };

    this.moveDown = function (index) {
        this.$post("/proxy/locations/moveDown")
            .params({
                "filename": this.filename,
                "index": index
            });
    };
});