Tea.context(function () {
   this.deletePolicy = function (filename) {
       if (!window.confirm("确定概要删除此缓存策略吗？")) {
            return;
       }

        this.$post(".deletePolicy")
            .params({
                "filename": filename
            })
            .success(function () {
                window.location.reload();
            });
   };
});