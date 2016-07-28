if( window.location.hostname === 'shangqu.local.wboll.com' ) {
    $.ajaxSettings = $.extend(
        $.ajaxSettings,
        {
            xhrFields: {withCredentials: true},
            crossDomain: true,
            contentType:'application/json',
            beforeSend: function(xhr, settings) {
                settings.url = (settings.url).indexOf('resource')!==-1? settings.url : 'http://shangqu.dev.wboll.com/shangqu' + settings.url;
                // settings.url = (settings.url).indexOf('resource')!==-1? settings.url : 'http://shangqu.dev.wboll.com' + settings.url;
            }
        }
    );
}else if(window.location.hostname === 'uat.shennongke.com' ){
    $.ajaxSettings = $.extend(
        $.ajaxSettings,
        {
            xhrFields: {withCredentials: true},
            crossDomain: true,
            contentType:'application/json',
            beforeSend: function(xhr, settings) {
                settings.url = (settings.url).indexOf('resource')!==-1? settings.url : 'http://uat.shennongke.com/shangqu' + settings.url;
                // settings.url = (settings.url).indexOf('resource')!==-1? settings.url : 'http://shangqu.dev.wboll.com' + settings.url;
            }
        }
    );
} else if(window.location.hostname === 'wb.shennongke.com' ){
    $.ajaxSettings = $.extend(
        $.ajaxSettings,
        {
            xhrFields: {withCredentials: true},
            crossDomain: true,
            contentType:'application/json',
            beforeSend: function(xhr, settings) {
                settings.url = (settings.url).indexOf('resource')!==-1? settings.url : 'http://wb.shennongke.com/' + settings.url;
                // settings.url = (settings.url).indexOf('resource')!==-1? settings.url : 'http://shangqu.dev.wboll.com' + settings.url;
            }
        }
    );
}

