Tea.context(function () {
   this.statusChanged = false;

   this.servers = this.servers.$map(function (k, server) {
       server.config.backends = server.config.backends.$map(function (_, backend) {
           return backend.address;
       });
       return server;
   });

   this.refreshStatus = function () {
        this.$get("/proxy/status")
            .success(function (response) {
                this.statusChanged = response.data.changed;
            })
            .done(function () {
                this.$delay(function () {
                    this.refreshStatus();
                }, 3000);
            });
   };

   this.refreshStatus();

   this.restart = function () {
       this.$get("/proxy/restart");
   };
});