var folderList = [];
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

        
        allowDrop(scope, element, attrs,$http);
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
                            //$http.post("/api/getSubFiles",item)
                            $http.get("/api/getSubFilesDirectPath?path="+item.Path)
                            .then(function(response) {
    
                                console.log("response",response.data);
                                //Add child items
                                
                                //scope.childItems["item"+item.IdPath] = {details: item , items: response.data.Items, backItem: item, forwaredItem: null};
                                scope.childItems["item"+item.IdPath] = response.data;
                                scope.childItems["item"+item.IdPath+"previous"] = [];
                                //scope.childItems["item"+item.ParentPath] = {details: item , items: response.data, backItem: item, forwaredItem: null};
                                

                                //var uniqueFolderID = Date.UTC(dtm.getFullYear(), dtm.getMonth(), dtm.getDate(),dtm.getHours(),dtm.getMinutes(),dtm.getSeconds(),dtm.getMilliseconds());var uniqueFolderID = Date.UTC(dtm.getFullYear(), dtm.getMonth(), dtm.getDate(),dtm.getHours(),dtm.getMinutes(),dtm.getSeconds(),dtm.getMilliseconds());
                                //scope.childItems["item"+uniqueFolderID] = {folderId:uniqueFolderID, details: item, items: response.data, backItem: item, forwaredItem: null};
    
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
        var originPath = "item"+scope.$parent.c.item.IdPath;



        jElement.on("dblclick", function(){
            
            if (scope.item.IsFolder){

                var parentContainer = jElement.closest(".folder-container");
                //parentContainer.attr("previous",parentContainer.attr("path"));
                
    
                scope.childItems[originPath+"previous"].push(parentContainer.attr("path"));
                parentContainer.attr("path",scope.item.Path);
                //$http.post("/api/getSubFiles",scope.item)
                $http.get("/api/getSubFilesDirectPath?path="+scope.item.Path)
                .then(function(response) {
                    //parentContainer.attr("path",response.data.IdPath);
                    
                    //scope.childItems[originPath] = {details: scope.item , items: response.data, backItem: scope.item, forwaredItem: null};
                    scope.childItems[originPath] = response.data;

                    var btnBack = parentContainer.find(".back");
                    if (!btnBack.is("[click]")){
                        btnBack.attr("click","true");
    
                            btnBack.on("click",function(){
                                
                                var prev = scope.childItems[originPath+"previous"];
                                $http.get("/api/getSubFilesDirectPath?path="+prev[prev.length-1])
                                .then(function(responseBack) {
                                    //scope.childItems[originPath] = {details: responseBack.data , items: responseBack.data.Items, backItem: scope.item, forwaredItem: null};
                                    scope.childItems[originPath] = responseBack.data;
                                    prev.pop();
                                    scope.$apply();
                                }, function(response) {
                                    console.log("error",response);
                                });
                            });
                    }
    
                }, function(response) {
                    console.log(response);
                });
                

                scope.$apply();
            }
            else{
                alert("This is a file.");
            }

        });
    };
});




mainApp.directive("initializeChildContainer",function($http){
    return function(scope, element, attrs){
        allowDrop(scope, element, attrs,$http);
        startDragOver(scope, element, attrs);
    };
});

function allowDrop(scope, element, attrs,$http){
    element.on("drop",function(e){
        
        e.preventDefault();
        var data = e.dataTransfer.getData("text");
        var elem = document.getElementById(data);
        

        

        //e.target.appendChild(elem);
        
        var $target = $(e.target);
        var $elem = $(elem);

        $("#moveConfirmation").dialog({
            resizable: false,
            height: "auto",
            width: 400,
            modal: false,
            buttons: {
              "Copy File": function() {
                executeMoveOrCopy(e, scope, element, attrs, $http, data, elem, "copy");
                $( this ).dialog( "close" );
              },
              "Move File": function() {
                executeMoveOrCopy(e, scope, element, attrs, $http, data, elem, "move");
                $( this ).dialog( "close" );
              },
              Cancel: function() {
                $( this ).dialog( "close" );
              }
            }
          });
        

        

    });
}

function executeMoveOrCopy(e, scope, element, attrs, $http, data, elem, action){
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
        //from desktop to folder
        //do logic for changing parent id
        //logic for reloading folder items

        //console.log("move",$elem.attr("path"))
        //console.log("current folder",$target.closest(".folder-container").attr("path"))

        var url = "";
        if (action=="move"){
            url = "/api/moveFile";
        }else{
            url = "/api/copyFile";
        }

        $http.post(url,{"Source":$elem.attr("path"), "Destination":$target.closest(".folder-container").attr("path")+"/"+$elem.attr("name")})
        .then(function(responseBack) {
            console.log("Move file message",responseBack)
            if (action=="move"){
                e.target.appendChild(elem);
            }
        }, function(response) {
            console.log("error",response);
        });


        


    }else if ($target.hasClass("item") || $target.closest(".item").hasClass("item")){
        //do logic for changing parent id
        console.log("second");
        elem.remove();
    }
};

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
