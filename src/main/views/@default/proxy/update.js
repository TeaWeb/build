Tea.context(function () {
    this.names = this.server.name.$map(function (_, name) {
        return {
            "key": Tea.key(),
            "name": name
        };
    });

    this.listenArray = this.server.listen.$map(function (_, address) {
        return {
            "key": Tea.key(),
            "address": address,
            "host": address.substring(0, address.lastIndexOf(":")),
            "port": address.substring(address.lastIndexOf(":") + 1)
        };
    });
    this.backendsArray = this.server.backends.$map(function (_, backend) {
        var address = backend.address;
        return {
            "key": Tea.key(),
            "host": address.substring(0, address.lastIndexOf(":")),
            "port": address.substring(address.lastIndexOf(":") + 1)
        };
    });

    this.addName = function () {
        this.names.push({
            "key": Tea.key(),
            "name": ""
        });
    };

    this.removeName = function (index) {
        this.names.$remove(index);
    };

    this.addListen = function () {
        this.listenArray.push({
            "key": Tea.key(),
            "address": ":",
            "host": "",
            "port": ""
        });
    };

    this.removeListen = function (index) {
        this.listenArray.$remove(index);
    };

    this.addBackend = function () {
        this.backendsArray.push({
            "key": Tea.key(),
            "host": "",
            "port": ""
        });
    };

    this.removeBackend = function (index) {
        this.backendsArray.$remove(index);
    };

});