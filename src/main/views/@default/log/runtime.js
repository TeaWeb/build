Tea.context(function () {
    var that = this;

    this.logText = "";
    this.stopScrolling = false;

    this.$delay(function () {
        this.$find(".log-box").bind("wheel", function () {
            if (this.scrollHeight < this.scrollTop + this.offsetHeight + 10) {
                that.stopScrolling = false;
            } else {
                that.stopScrolling = true;
            }
        });
        this.$find(".log-box").bind("scroll", function () {

        });
        this.load();
    });

    this.load = function () {
        if (!this.stopScrolling) {
            this.$find(".log-box")[0].scrollTop = 10000;
        }

        this.$post(".runtime")
            .success(function (resp) {
                this.logText = this.logText + resp.data.data
                    .replace(/ /g, "&nbsp;")
                    .replace(/\t/, "&nbsp; &nbsp; ")
                    .replace(/\n/g, "<br/>");
            })
            .done(function () {
                this.$delay(function () {
                    this.load();
                }, 1000);
            });
    };
});