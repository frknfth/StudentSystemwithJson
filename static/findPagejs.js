
angular.module("myApp",[]).controller("myCtrl", function($scope, $http){

    $scope.id=0;

    $scope.student = "";
    
    $scope.order=function (x) {
        $http.get("/get",{headers:{'Id':$scope.id}}).then(function mySuccess(response) {
            $scope.student = response.data;
        }, function myError(response) {
            $scope.student = response.statusText;
        });
    }

})