var validate = function(container){
    if (container.find("input.ng-invalid").length > 0) return false;
    else return true;
};