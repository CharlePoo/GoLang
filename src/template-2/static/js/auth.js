var appAuth = angular.module('appAuth', []);

appAuth.controller('ctrlReg', function($scope, $http){
    $scope.R = {
        FirstName: "",
        LastName: ""
    };

    
    $("#registration").removeClass("ng-invalid");
    

    $scope.SubmitForm = function(){

        
        validate($("#registration"));
        alert("ayos");
        return;
        $http.post("/api/register", $scope.R)
        .then(function(response) {
            //First function handles success
            //$scope.content = response.data;
            console.log($scope);
        }, function(response) {
            //Second function handles error
            $scope.content = "Something went wrong";
        });
        
    };
});

$( function() {
    /*$( "#inputBirthdate" ).datepicker({
      changeMonth: true,
      changeYear: true,
      yearRange: "1950:2017"
    });*/
  } );