Tea.context(function () {
    this.isRefreshing = false;

   this.reloadApp = function (appId) {
       this.isRefreshing = true;
       this.$get("/apps/app/reload")
           .params({
               "appId": appId
           })
           .success(function (resp) {
                var index = -1;
                var app = resp.data.app;
                this.apps.$each(function (k, v) {
                    if (v.id == app.id) {
                        index = k;
                    }
                });
                if (index >= 0) {
                    this.apps[index] = app;
                    this.$set(this.apps, index, app);
                }
           })
           .done(function () {
               this.$delay(function () {
                   this.isRefreshing = false;
               }, 500);
           });
   };

   this.showAppDetail = false;
   this.detailTab = "system";

   this.showApp = function (appId) {
       this.$get("/apps/app")
           .params({
               "appId": appId
           })
           .success(function (resp) {
               this.detailTab = "system";
               this.detailApp = resp.data.app;

               this.detailApp.version = this.detailApp.version.replace(/ /g, "&nbsp;").replace(/\n/g, "<br/>");

               this.detailApp.memory = resp.data.memory;
               this.detailApp.memoryRSS = resp.data.memoryRSS;
               this.detailApp.memoryVMS = resp.data.memoryVMS;
               this.detailApp.cpu = resp.data.cpu;

               var listens = [];
               var connections = [];
               var openFiles = [];
               this.detailApp.processes.$each(function (k, v) {
                   v.listens.$each(function (_, listen) {
                       if (!listens.$contains(listen.network + " " + listen.addr)) {
                           listens.push(listen.network + " " + listen.addr);
                       }
                   });

                   v.connections.$each(function (_, connection) {
                        if (!connections.$contains(connection)) {
                            connections.push(connection);
                        }
                   });

                   v.openFiles.$each(function (_, openFile) {
                       if (!openFiles.$contains(openFile)) {
                           openFiles.push(openFile);
                       }
                   });
               });
               this.detailApp.listens = listens;
               this.detailApp.connections = connections;
               this.detailApp.openFiles = openFiles;

               this.showAppDetail = true;
           });
   };

   this.closeApp = function () {
       this.showAppDetail = false;
   };
});