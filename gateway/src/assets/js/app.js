'use strict';

(function () {
    function init() {
        var router = new Router([   
            new Route('home', 'home.html', true),       
            new Route('services', 'services/search.html'),
            new Route('users/admin', 'users/admin.html'),
            // new Route('lession2', 'lezione.html',false,new CMS("lesson2", i18n, flowchart)),
            // new Route('lession3', 'lezione.html',false,new CMS("lesson3", i18n, flowchart))
        ]);
    }
    init();
}());
