Tea.context(function () {
    this.name = "";
    this.developer = "";
    this.site = "http://";
    this.docSite = "http://";
    this.commandName = "";
    this.commandPatterns = [];
    this.commandVersion = "";

    this.detailVisible = false;
    this.showMore = function () {
        this.detailVisible = !this.detailVisible;
    };

    this.addPattern = function () {
        this.commandPatterns.push("");
        this.$delay(function () {
            this.$find("form input[name='commandPatterns']")[this.commandPatterns.length - 1].focus();
        });
    };

    this.removePattern = function (index) {
        this.commandPatterns.$remove(index);
    };

    this.results = [];
    this.loading = false;
    this.test = function () {
        this.loading = true;

        this.results = [];
        this.$post(".addProbe")
            .params({
                "isTesting": 1,
                "name": this.name,
                "developer": this.developer,
                "site": this.site,
                "docSite": this.docSite,
                "commandName": this.commandName,
                "commandVersion": this.commandVersion,
                "commandPatterns": this.commandPatterns
            })
            .success(function (resp) {
                this.results = resp.data.apps;
                this.$delay(function () {
                    window.scroll(0, 10000);
                });
            })
            .done(function () {
                this.loading = false;
            });
    };
});