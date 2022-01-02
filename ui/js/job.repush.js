// rePushJob is a function that re-pushes job
function rePushJob(id) {

    let csrf = $("#csrf").val()

    // Loading table
    loadingTable()

    $.ajax({
        url: `../../staff/job/${id}/push/${csrf}`,
        type: "GET",
        success: successFunc,
        error: errorFunc,
    });

    function successFunc(result, status, xhr) {

        // Unloading table
        unLoadingTable()

        $("#push-success-div-msg").
            text(`You have successfully re-pushed the selected job!`)
        hideAllExceptInModalP("push-success-div")

    }

    function errorFunc(xhr, status, error) {

        // Unloading table
        unLoadingTable()

        $("#push-error-div-msg").
            text(`Unable to push the selected job!`)
        hideAllExceptInModalP("push-error-div")
    }

}

function hideAllExceptInModalP(elementID) {
    $("#push-error-div").hide()
    $("#push-success-div").hide()

    $(`#${elementID}`).show()
    $("#push-job").modal("show")
}