Tea.context(function () {
    this.logs = [];
    this.isLoaded = false;
    this.lastId = "";

    this.$delay(function () {
        this.loadLogs();
    });

    this.loadLogs = function () {
        this.$post("/agents/apps/taskLogs")
            .params({
                "agentId": this.agentId,
                "taskId": this.task.id,
                "lastId": this.lastId
            })
            .success(function (resp) {
                if (resp.data.logs.length == 0) {
                    return;
                }
                this.lastId = resp.data.logs.$first().id;
                this.logs = resp.data.logs.$map(function (k, v) {
                    v.datetime = v.timeFormat.second.substr(0, 4) + "-" + v.timeFormat.second.substr(4, 2) + "-" + v.timeFormat.second.substr(6, 2) + " " + v.timeFormat.second.substr(8, 2) + ":" + v.timeFormat.second.substr(10, 2) + ":" + v.timeFormat.second.substr(12);
                    return v;
                }).concat(this.logs);
            })
            .fail(function () {

            })
            .done(function () {
                this.isLoaded = true;

                this.$delay(function () {
                    this.loadLogs();
                }, 3000);
            });
    };
});