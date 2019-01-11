Tea.context(function () {
   this.authUpdating = this.config.username != null && this.config.username.length > 0;
   this.isTesting = false;
   this.testingError = "";
   this.testingSuccess = "";

   this.updateAuth = function () {
        this.authUpdating = !this.authUpdating;
   };

   this.testConnection = function () {
        var params = {
            host: this.config.host,
            port: this.config.port
        };
        if (this.authUpdating) {
            params["username"] = this.config.username;
            params["password"] = this.config.password;
        }

        this.isTesting = true;
        this.testingError = "";
        this.testingSuccess = "";

        this.$get(".test")
            .params(params)
            .success(function () {
                this.testingError = "";
                this.testingSuccess = "连接成功！";
            })
            .fail(function (resp) {
                if (resp) {
                    this.testingError = resp.message;
                }
            })
            .done(function () {
                this.isTesting = false;
            });
   };
});