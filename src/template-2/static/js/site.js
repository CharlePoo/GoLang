var mainApp = angular.module("mainApp",[]);
mainApp.controller("mainCtr",function($scope,$http){
    $http.get("/api/initialize")
    .then(function(response) {
        //console.log(response);
        $scope.childItems = {};
        $scope.column = [];
        $scope.data = response.data;
 
        

        $scope.row = [];
        for (var rowIndex=0; rowIndex < response.data.Desktop.RowSize; rowIndex++){
            //$scope.row.push(row);
            var colList = [];
            for (var colIndex=0; colIndex < response.data.Desktop.ColumnSize; colIndex++){
                //$scope.column.push({col: col, item: null});
                colList.push({col: colIndex, item: null});
            }
            $scope.row.push({row:rowIndex,column:colList });
        }
        
        //$scope.columnSize = response.data.Desktop.ColumnSize;
        //$scope.rowSize = response.data.Desktop.RowSize;
        console.log($scope);
        
    }, function(response) {
        //Second function handles error
        $scope.content = "Something went wrong";
    });
});

mainApp.directive('folderContainer', function ($compile,$http) {
    return function (scope, element, attrs){

        
        allowDrop(scope, element, attrs);
        startDragOver(scope, element, attrs);
        allowDrag(scope, element, attrs);
        

        _.find(scope.data.Items,function(item){
            var id = item.RowIndex.toString() + item.ColumnIndex.toString();
            
            if (attrs.indexid==id){
                scope.c.item = item;
                
                
                _.defer(function(){
                    var jElement = $(element);
                    //ondrop
                    
                    jElement.find(".item-element").on( "dblclick", function() {

                        if (item.IsFolder){
                            $http.post("/api/getSubFiles",item)
                            .then(function(response) {
    
                                //Add child items
                                scope.childItems["item"+item.IdPath] = {details: item , items: response.data};
    
                                //console.log("child",scope);
                            }, function(response) {
                                //Second function handles error
                                //$scope.content = "Something went wrong";
                                console.log(response);
                            });
    
                            
                            
                            var _compiled = _.template($("#modalTemplate").html());
                            var _template = $(_compiled({scopeItem: "childItems."+"item"+item.IdPath}));
                            
    
                            //https://docs.angularjs.org/guide/compiler
                            var newEl = angular.element(_template);
                            element.append(newEl);
                            $compile(newEl)(scope);
                            setTimeout(function(){
                                _template.dialog({
                                    resizable: false,
                                    height: "auto",
                                    width: 400,
                                    modal: false,
                                    buttons: {
                                      "Delete all items": function() {
                                        $( this ).dialog( "close" );
                                      },
                                      Cancel: function() {
                                        $( this ).dialog( "close" );
                                      }
                                    }
                                  });
                            },100);
                        }else{
                            alert("This is a file");
                            console.log("This is a file");
                        }
                        
                        
                        scope.$apply(); //This will apply all changes for scope
                    });
                    
                });
                
                //element.on( "click", function() {
                //    alert(id);
                //});
                return true;
            }
            
            return false;
        });

        
    };
});

mainApp.directive("childItems",function($compile, $http){
    return function(scope, element, attrs){
        allowDrag(scope, element, attrs);
        var jElement = $(element);
        jElement.on("dblclick", function(){
            console.log(jElement.parent());
            //scope.$parent.c.item.Name = "test";
            //scope.$apply();

        });
    };
});

mainApp.directive("initializeChildContainer",function(){
    return function(scope, element, attrs){
        allowDrop(scope, element, attrs);
        startDragOver(scope, element, attrs);
    };
});

function allowDrop(scope, element, attrs){
    element.on("drop",function(e){
        
        e.preventDefault();
        var data = e.dataTransfer.getData("text");
        var elem = document.getElementById(data);
        e.target.appendChild(elem);
        
        var $target = $(e.target);
        var $elem = $(elem);
        

        if ($target.hasClass("desktop")){
            //do logic for changing parent id
            var children = $target.children();
            console.log($target.children().length);
            if (children.length > 1 && children.hasClass("item")){
                $elem.remove();
            }else{
                //logic
            }
        }else if ($target.closest(".folder-container").hasClass("folder-container")){
            console.log("first");
            //do logic for changing parent id
            //logic for reloading folder items
        }else if ($target.hasClass("item") || $target.closest(".item").hasClass("item")){
            //do logic for changing parent id
            console.log("second");
            elem.remove();
        }

    });
}

function startDragOver(scope, element, attrs){
    element.on("dragover",function(e){
        e.preventDefault();
    });
}

function allowDrag(scope, element, attrs){
    element.on("dragstart", function(e){
        e.dataTransfer.setData("text", e.target.id);
    });
}





/*mainApp.directive("folderContainer",function(){
    
    var directive = {};
    console.log(1);
    //directive.restrict = 'E';
    directive.template = $("#item-template").html();
    directive.compile = function (element, attributes) {
        
        //if (scope.$last) 
        var linkFunc = function($scope, element, attrs){
            
        };
        return linkFunc;
    };
    return directive;
});*/
