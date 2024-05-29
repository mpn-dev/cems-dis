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

function getSensors() {
  return [
    {code: "so2",         name: "SO2",          unit: "", gauge_min: 0, gauge_max: 100, major: 10, minor: 5, dials: [{min: 0, max: 100, color: "#00ff00"}]}, 
    {code: "nox",         name: "NOX",          unit: "", gauge_min: 0, gauge_max: 100, major: 10, minor: 5, dials: [{min: 0, max: 100, color: "#00ff00"}]}, 
    {code: "pm",          name: "PM",           unit: "", gauge_min: 0, gauge_max: 100, major: 10, minor: 5, dials: [{min: 0, max: 100, color: "#00ff00"}]}, 
    {code: "h2s",         name: "H2S",          unit: "", gauge_min: 0, gauge_max: 100, major: 10, minor: 5, dials: [{min: 0, max: 100, color: "#00ff00"}]}, 
    {code: "opacity",     name: "Opacity",      unit: "", gauge_min: 0, gauge_max: 100, major: 10, minor: 5, dials: [{min: 0, max: 100, color: "#00ff00"}]}, 
    {code: "flow",        name: "Flow",         unit: "", gauge_min: 0, gauge_max: 100, major: 10, minor: 5, dials: [{min: 0, max: 100, color: "#00ff00"}]}, 
    {code: "o2",          name: "O2",           unit: "", gauge_min: 0, gauge_max: 100, major: 10, minor: 5, dials: [{min: 0, max: 100, color: "#00ff00"}]}, 
    {code: "temperature", name: "Temperature",  unit: "", gauge_min: 0, gauge_max: 100, major: 10, minor: 5, dials: [{min: 0, max: 100, color: "#00ff00"}]}, 
    {code: "pressure",    name: "Pressure",     unit: "", gauge_min: 0, gauge_max: 100, major: 10, minor: 5, dials: [{min: 0, max: 100, color: "#00ff00"}]}, 
  ];
}
