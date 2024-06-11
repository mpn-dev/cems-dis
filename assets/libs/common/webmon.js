var mainApp = angular.module("mainApp", ['ngUpload']);

mainApp.service("rights", function($http) {
  var data = {};

  data.load = function(scopeId, successFn = null) {
    $http.get('/api/v1/access-right/' + scopeId, requestConfig()).then(
      function(response) {
        Object.assign(data, response.data.data);
        if(typeof successFn === 'function') {
          successFn();
        }
      }
    )
  }

  return data;
});
