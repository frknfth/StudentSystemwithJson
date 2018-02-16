
angular.module("myApp",[]).controller("myCtrl", function ($scope, $http){

    $scope.stname = "";
    $scope.stage = 0;
    $scope.stuni = "";
    $scope.stid = 0;
    $scope.message = "";

    $scope.sendData = function () {
        var url = "/save", data = {"name": $scope.stname,"age": $scope.stage,
                                   "uni": $scope.stuni,"id": $scope.stid};

        $http.post(url, data).then(function mySuccess(response) {
            $scope.message = response.data;
        }, function myError(response) {
            $scope.message = response.statusText;
        });
    }

});