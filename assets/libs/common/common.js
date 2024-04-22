function getCsrfToken() {
  return document.querySelector('meta[name="csrf-token"]').getAttribute('content');
}

function requestConfig(params = null) {
  return {
    headers: {
      "Content-Type":   "application/json", 
      "Token":          getCsrfToken()
    }
  };
}

function strToUtcTimestamp(value) {
  var d1 = new Date(value);
  if(isNaN(d1)) {
    return 0;
  }
  var d2 = new Date(d1.getFullYear(), d1.getMonth(), d1.getDate());
  // return (d2.valueOf() - (d2.getTimezoneOffset() * 60000)) / 1000;
  return d2.valueOf() / 1000;
}

function getDateRange(datepicker1, datepicker2, name1, name2, allowEmpty = false) {
  var ts1 = strToUtcTimestamp(datepicker1.value());
  var ts2 = strToUtcTimestamp(datepicker2.value());
  if((ts1 == 0) && (ts2 == 0)) {
    if(allowEmpty) {
      return null;
    } else {
      return new Error(name1 + " dan " + name2 + " wajib diisi");
    }
  } else if(ts1 == 0) {
    return new Error(name1 + " tidak valid");
  } else if(ts2 == 0) {
    return new Error(name2 + " tidak valid");
  } else if(ts2 < ts1) {
    return new Error(name2 + " tidak boleh lebih kecil dari " + name1);
  }

  return {start: ts1, end: ts2 + 24 * 60 * 60 - 1};
}

function defaultPagingOptions() {
  return {
    first: '<i class="fa fa-step-backward" />', 
    prev: '<i class="fa fa-chevron-left" />', 
    next: '<i class="fa fa-chevron-right" />', 
    last: '<i class="fa fa-step-forward" />'
  };
}

function defaultErrorHandler(resp) {
  responses = Array.isArray(resp) ? resp : [resp];
  responses.forEach(function(r) {
    if(r.status != 200) {
      var errmsg = "Unknown error occured";
      if(typeof r.data === 'string' || r.data instanceof String)
        errmsg = r.data;
      else if("error" in r.data) {
        errmsg = r.data.error;
      }
      window.alert("Error " + r.status + ": " + errmsg);
      return;
    }
  });
}
