Tea.context(function () {
    this.refreshingAppId = "";
    this.hasFavoredApps = false;

    this.$delay(function () {
        this.watch();
    });

    this.hasFavoredApps = this.apps.$any(function (k, v) {
        return v.isFavored;
    });
    this.hasNotFavoredApps = this.apps.$any(function (k, v) {
        return !v.isFavored;
    });

   this.reloadApp = function (appId) {
       this.refreshingAppId = appId;
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
                   this.refreshingAppId = "";
               }, 500);
           });
   };

   this.watch = function () {
        this.$post(".watch")
            .params({
                "appIds": this.apps.$map(function (k, v) {
                    return v.id;
                })
            })
            .success(function (resp) {
                if (resp.data.isChanged) {
                    this.apps = resp.data.apps;
                }
            })
            .done(function () {
                this.$delay(function () {
                    this.watch();
                }, 1000);
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
               this.detailApp.pluginName = resp.data.plugin;

               this.detailApp.version = this.detailApp.version.replace(/ /g, "&nbsp;").replace(/\n/g, "<br/>");

               this.detailApp.memoryPercent = resp.data.memoryPercent;
               this.detailApp.memoryRSS = resp.data.memoryRSS;
               this.detailApp.memoryVMS = resp.data.memoryVMS;
               this.detailApp.cpuPercent = resp.data.cpuPercent;

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

   this.favor = function (app) {
       this.$post("/apps/app/favor")
           .params({
               "appId": app.id
           })
           .success(function () {
                app.isFavored = true;

               this.hasFavoredApps = this.apps.$any(function (k, v) {
                   return v.isFavored;
               });

               this.hasNotFavoredApps = this.apps.$any(function (k, v) {
                   return !v.isFavored;
               });
           });
   };

    this.cancelFavor = function (app) {
        this.$post("/apps/app/cancelFavor")
            .params({
                "appId": app.id
            })
            .success(function () {
                app.isFavored = false;

                this.hasFavoredApps = this.apps.$any(function (k, v) {
                    return v.isFavored;
                });

                this.hasNotFavoredApps = this.apps.$any(function (k, v) {
                    return !v.isFavored;
                });
            });
    };
});