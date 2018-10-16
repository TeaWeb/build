Tea.context(function () {
   this.startMongo = function () {
       this.$post(".install")
           .success(function () {
               window.location.reload();
           });
   };
});