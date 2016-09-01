
// ( window.location.hostname === 'shangqu.local.wboll.com' ) { 
  $.ajaxSettings = $.extend(
				$.ajaxSettings,
				{
					xhrFields: {withCredentials: true},
					crossDomain: true,
					contentType:'application/json',
					beforeSend: function(xhr, settings) {
						// setting.url = 'http://www.mylvfa.com'+setting.url
					}
				}
	); 
// }
