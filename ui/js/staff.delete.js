function initDeleteStaff(id){

    $("#delete-password").val("");
    $("#delete-member-id").text(id);
    $("#delete-error-div").hide();
    $("#delete-password-error").css("visibility","hidden");
    $("#delete-staff-member .modal-footer").show();
    $("#delete-member-container").show();
}

function deleteStaff(btn){

    let csrf = $("#csrf").val()
    let id = $("#delete-member-id").text();
    let password = $("#delete-password").val();

    if (password.match(/^\s*$/)) {
        return
    }

    loading(btn);

    // cleaning up before processing
    $("#delete-password-error").css("visibility","hidden");

    $.ajax({
        url: `../../admin/staff/${id}/${csrf}`,
        data: {password},
        type: "POST",
        success: successFunc,
        error: errorFunc,
    });

    function successFunc(result, status, xhr) {
        $('#delete-staff-member').modal('hide');
        unLoading(btn, "Delete");

        // Reloading the table content
        initSearchStaff()
    }

    function errorFunc(xhr, status, error) {

        unLoading(btn, "Delete");

        if (xhr.status == 400){
            $("#delete-password-error").text("Incorrect password!")
            $("#delete-password-error").css("visibility","visible");
        }else{
            $("#delete-member-container").hide();
            $("#delete-staff-member .modal-footer").hide();
            $("#delete-error-div-msg").text("Unable to delete the selected staff member.")
            $("#delete-error-div").show();
        }

    }
}