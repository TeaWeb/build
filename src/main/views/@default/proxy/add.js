Tea.context(function () {
    this.description = "";
    this.serviceType = 1;
    this.root = "";

    this.name = "";
    this.nameArray = [];

    this.listen = "";
    this.listenArray = [];

    this.backend = "";
    this.backendArray = [];

    this.localPaths = [];

    this.$delay(function () {
        this.$find("form input[name='description']").focus();

        this.$watch("name", function (newValue) {
            this.nameArray = newValue.trim().split(/\s+/);
        });

        this.$watch("listen", function (newValue) {
            this.listenArray = newValue.trim().split(/\s+/);
        });

        this.$watch("backend", function (newValue) {
            this.backendArray = newValue.trim().split(/\s+/);
        });
    });

    this.changeRoot = function (root) {
        this.$get(".localPath")
            .params({
                "prefix": root
            })
            .success(function (resp) {
                this.localPaths = resp.data.paths;
            });
    };

    this.selectRoot = function (root) {
        this.root =root;
        this.localPaths = [];
        this.$delay(function () {
            this.$find("#web-root-input").focus();
        });
    };
});