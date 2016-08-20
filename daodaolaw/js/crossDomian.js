
// ( window.location.hostname === 'shangqu.local.wboll.com' ) { 
  $.ajaxSettings = $.extend(
				$.ajaxSettings,
				{
					xhrFields: {withCredentials: true},
					crossDomain: true,
					contentType:'application/json',
					beforeSend: function(xhr, settings) {
						// settings.url = (settings.url).indexOf('resource')!==-1? settings.url : 'http://www.mylvfa.com' + settings.url;
						// settings.url = (settings.url).indexOf('resource')!==-1? settings.url : 'http://shangqu.dev.wboll.com' + settings.url;
						// setting.url = 'http://www.mylvfa.com'+setting.url
					}
				}
	); 
// }
