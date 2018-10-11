Tea.context(function () {
   this.authUpdating = this.config.username != null && this.config.username.length > 0;

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

        this.$get(".test")
            .params(params)
            .success(function () {
                alert("连接成功");
            });
   };
});