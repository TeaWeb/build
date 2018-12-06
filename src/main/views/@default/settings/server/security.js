Tea.context(function () {

    this.addAllowIP = function () {
        this.security.allow.push("");
        this.$delay(function () {
            var input = this.$find("#allow-ips-box input[name='allowIPs']").last();
            if (input.length > 0) {
                input.focus();
            }
        });
    };

    this.deleteAllowIP = function (index) {
        this.security.allow.$remove(index);
    };

    this.addDenyIP = function () {
        this.security.deny.push("");
        this.$delay(function () {
            var input = this.$find("#deny-ips-box input[name='denyIPs']").last();
            if (input.length > 0) {
                input.focus();
            }
        });
    };

    this.deleteDenyIP = function (index) {
        this.security.deny.$remove(index);
    };
});