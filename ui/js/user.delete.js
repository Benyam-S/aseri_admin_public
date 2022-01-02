function initDeleteUser(id){

    $("#delete-password").val("");
    $("#delete-user-id").text(id);
    $("#delete-error-div").hide();
    $("#delete-password-error").css("visibility","hidden");
    $("#delete-user .modal-footer").show();
    $("#delete-user-container").show();
}

function deleteUser(btn){

    let csrf = $("#csrf").val()
    let id = $("#delete-user-id").text();
    let password = $("#delete-password").val();

    if (password.match(/^\s*$/)) {
        return
    }

    loading(btn);

    // cleaning up before processing
    $("#delete-password-error").css("visibility","hidden");

    $.ajax({
        url: `../../staff/user/${id}/${csrf}`,
        data: {password},
        type: "POST",
        success: successFunc,
        error: errorFunc,
    });

    function successFunc(result, status, xhr) {
        $('#delete-user').modal('hide');
        unLoading(btn, "Delete");

        // Reloading the table content
        initSearchUser()
    }

    function errorFunc(xhr, status, error) {

        unLoading(btn, "Delete");

        if (xhr.status == 400){
            $("#delete-password-error").text("Incorrect password!")
            $("#delete-password-error").css("visibility","visible");
        }else{
            $("#delete-user-container").hide();
            $("#delete-user .modal-footer").hide();
            $("#delete-error-div-msg").text("Unable to delete the selected user.")
            $("#delete-error-div").show();
        }

    }
}