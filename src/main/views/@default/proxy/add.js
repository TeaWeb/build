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
});