'use stict';

function CMS(contentName, i18n, flowchart) {
    try {
        if (!contentName || !i18n || !flowchart) {
            throw 'error: contentName, i18n, flowchart params are mandatories';
        }
        this.constructor(contentName, i18n, flowchart);
    } catch (e) {
        console.error(e);
    }
}

CMS.prototype = {
    contentName: undefined,
    i18n: undefined, 
    flowchart: undefined,
    constructor: function (contentName, i18n, flowchart) {
        this.contentName = contentName;
        this.i18n = i18n;
        this.flowchart = flowchart;
    },
    init: function () {
        (function (scope) {
            var url = 'i18n/' + scope.contentName  + ".json",
                xhttp = new XMLHttpRequest();
            xhttp.onreadystatechange = function () {
                if (this.readyState === 4 && this.status === 200) {
                    console.log(this.response);
                    var json = JSON.parse(this.response);
                    scope.i18n.init(json.dicts);
                    scope.i18n.setLocale("en");
                    scope.i18n.apply();

                    scope.flowchart.drawAll(json.flows, 'canvas');

                    // fix issue with flowchart draw in div/div with display=none
                    $(".collapse").collapse();

                    $('#START').on('click', function (event) {
                        event.preventDefault(); // To prevent following the link (optional)
                        $('#collapse0').collapse('show');
                    });

                    // solo per exercises: change dropdown text
                    $('.dropdown-item').on('click', function (event) {
                        event.preventDefault(); // To prevent following the link (optional)
                        var btnObj = $(this).parent().siblings('button');
                        $(btnObj).text($(this).text());
                        $(btnObj).val($(this).text());
                        $(btnObj).attr("data-selection", $(this).attr("href"));
                    });
                }
            };
            xhttp.open('GET', url, true);
            xhttp.send();
        })(this);
    }

}
