var appAuth = angular.module('appAuth', []);
appAuth.controller('ctrlReg', function($scope, $http){
    $scope.R = {
        FirstName: "",
        LastName: "",
        Email: "",
        BirthDate: new Date(),
        Password: "",
        ConfirmPassword:""
    };

    var alertReg = $("#alertReg").hide();
    $("#registration").removeClass("ng-invalid");
    var email = $("#inputEmail");

    email.on("click",function(){
        email.removeClass("email-invalid");
    });

    $scope.SubmitForm = function(){
        var isValid = false;
        var isPasswordConfirm = false;
        alertReg.hide("fade");
        isValid = validate($("#registration"));
        email.removeClass("email-invalid");
        if (isValid && $scope.R.Password == $scope.R.ConfirmPassword){
            $http.post("/api/register", $scope.R)
            .then(function(response) {
                //First function handles success

                if (response.data.ID>0){
                    //$scope.R = response.data;
                    location.href = "/";
                    
                }else{
                    alertReg.show("fade");
                    email.addClass("email-invalid");
                }
                
                console.log(response);
            }, function(response) {
                //Second function handles error
                $scope.content = "Something went wrong";
            });
        }
        
    };
});

appAuth.controller('ctrlLogin', function($scope, $http){
    $scope.L = {
        Email: "",
        Password: ""
    };

    $scope.SubmitForm = function(){
        
        $http.post("/api/login", $scope.L)
        .then(function(response) {
            //First function handles success
            console.log(response);
            if (response.data.ID>0){
                //$scope.R = response.data;
                location.href = "/";
                
            }else{
                //alertReg.show("fade");
                //email.addClass("email-invalid");
                alert("invalid login!");
            }
            
            console.log(response);
        }, function(response) {
            //Second function handles error
            $scope.content = "Something went wrong";
        });
    };
});