'use stict';

function Route(name, htmlName, defaultRoute, cms) {
    try {
        if(!name || !htmlName) {
            throw 'error: name and htmlName params are mandatories';
        }
        this.constructor(name, htmlName, defaultRoute, cms);
    } catch (e) {
        console.error(e);
    }
}

Route.prototype = {
    name: undefined,
    htmlName: undefined,
    default: undefined,
    cms: undefined,
    constructor: function (name, htmlName, defaultRoute, cms) {
        this.name = name;
        this.htmlName = htmlName;
        this.default = defaultRoute;
        this.cms = cms;
    },
    isActiveRoute: function (hashedPath) {
        return hashedPath.replace('#', '') === this.name; 
    }
}