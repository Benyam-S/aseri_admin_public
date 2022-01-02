// initEditStaff is a function that initiate or display the edit staff member window
function initEditStaff(id) {

    let csrf = $("#csrf").val()

    hideAllExcept($("#edit-section"));

    // cleaning up before processing
    $("#edit-member-error").removeClass("success").addClass("error")
    $(".error").css("visibility", "hidden");
    $(".error").text("");
    $("#edit-profile-pic").css("background-image", `url(../../ui/images/user.svg)`)
    $("#edit-member-body input").val("")
    cleanFrame()

    $("#edit-member-body").show();
    $("#edit-error-div").hide();

    removeBaseURI()
    addBaseURI("Edit")

    $.ajax({
        url: `../../admin/staff/${id}/${csrf}`,
        type: "GET",
        success: successFunc,
        error: errorFunc,
    });

    function successFunc(result, status, xhr) {

        let staffMember = JSON.parse(result);

        if (staffMember.ProfilePic) {
            $("#edit-remove-photo").show()
            $("#edit-hidden-profile").val(staffMember.ProfilePic)
            $("#edit-profile-pic").css("background-image", `url(../../assets/profilepics/${staffMember.ProfilePic})`)
        } else {
            $("#edit-remove-photo").hide()
        }

        $("#edit-first-name").css("text-transform", "capitalize").val(staffMember.FirstName)
        $("#edit-last-name").css("text-transform", "capitalize").val(staffMember.LastName)
        $("#edit-member-id").val(staffMember.ID)
        $("#edit-member-password").val("********")
        $("#edit-member-role").css({ "text-transform": "uppercase" }).val(staffMember.Role)
        $("#edit-member-email").val(staffMember.Email)
        $("#edit-phone-number").val(staffMember.PhoneNumber.replace("+251", "0"))

    }

    function errorFunc(xhr, status, error) {

        $("#edit-member-body").hide();
        $("#edit-error-div").show();

    }

}

function editPreviewPhoto(uploader) {
    if (uploader.files && uploader.files[0]) {
        var imageFile = uploader.files[0];
        var reader = new FileReader();
        reader.onload = function (e) {
            $("#edit-profile-pic").css("background-image", `url(${e.target.result})`);
        };

        reader.readAsDataURL(imageFile);

        // Uploading to the database
        let csrf = $("#csrf").val()
        let id = $("#edit-member-id").val()

        let spinner = $("<span> </span>").addClass("spinner-border spinner-border-sm").
            attr({ "role": "status", "aria-hidden": true })
        $("#edit-uploader").val("")
        $("#edit-upload-photo").hide()
        $("#edit-remove-photo").hide()
        $("#edit-profile-pic > div").css("visibility", "visible")
        $("#edit-uploading-photo").show().append(spinner)

        let formData = new FormData()
        formData.append("profile_pic", imageFile)

        $.ajax({
            url: `../../admin/staff/${id}/profile/pic/${csrf}`,
            type: "POST",
            data: formData,
            processData: false,
            contentType: false,
            success: successFunc,
            error: errorFunc,
        });

        function successFunc(result, status, xhr) {

            let fileNameContainer = JSON.parse(xhr.responseText)
            cleanFrame()
            $("#edit-uploading-photo").find(".spinner-border").remove()
            $("#edit-remove-photo").show()
            $("#edit-upload-photo").show()
            $("#edit-hidden-profile").val(fileNameContainer["file_name"])

        }

        function errorFunc(xhr, status, error) {

            let err = ""
            if (xhr.status == 400) {
                let errMap = JSON.parse(xhr.responseText)
                err = errMap["error"]

            } else {
                err = "Unable To Upload Photo"
            }
            cleanFrame()
            $("#edit-uploading-photo").find(".spinner-border").remove()
            $("#edit-reload-photo").show()
            $("#edit-upload-photo").show()
            $("#edit-photo-error").css("text-transform", "capitalize").text(err).css("visibility", "visible");
        }

    }
}

function editRemovePhoto() {

    let csrf = $("#csrf").val()
    let id = $("#edit-member-id").val()

    let spinner = $("<span> </span>").addClass("spinner-border spinner-border-sm").
        attr({ "role": "status", "aria-hidden": true })
    $("#edit-uploader").val("")
    $("#edit-upload-photo").hide()
    $("#edit-remove-photo").hide()
    $("#edit-profile-pic > div").css("visibility", "visible")
    $("#edit-uploading-photo").show().append(spinner)


    if ($("#edit-hidden-profile").val()) {

        $.ajax({
            url: `../../admin/staff/${id}/profile/pic/${csrf}`,
            type: "DELETE",
            success: successFunc,
            error: errorFunc,
        });

        function successFunc(result, status, xhr) {
            cleanFrame()
            $("#edit-uploading-photo").find(".spinner-border").remove()
            $("#edit-upload-photo").show()
            $("#edit-hidden-profile").val("")
            $("#edit-profile-pic").css("background-image", `url(../../ui/images/user.svg)`)

        }

        function errorFunc(xhr, status, error) {
            cleanFrame()
            $("#edit-uploading-photo").find(".spinner-border").remove()
            $("#edit-upload-photo").show()
            $("#edit-remove-photo").show()
            $("#edit-photo-error").css("text-transform", "capitalize").text("Unable To Delete Profile Pic").css("visibility", "visible");
        }

    } else {
        cleanFrame()
        $("#edit-upload-photo").show()
        $("#edit-profile-pic").css("background-image", `url(../../ui/images/user.svg)`)
    }


}

function editStaff(btn) {

    let csrf = $("#csrf").val()
    let id = $("#edit-member-id").val()

    // cleaning up before processing
    $("#edit-member-error").removeClass("success").addClass("error")
    $(".error").css("visibility", "hidden");
    $(".error").text("");

    // Adding loading spinner
    loadButtonW(btn)

    let firstName = $("#edit-first-name").val()
    let lastName = $("#edit-last-name").val()
    let email = $("#edit-member-email").val()
    let phoneNumber = $("#edit-phone-number").val()


    let errMap = validateStaffProfile(firstName, lastName, email, phoneNumber)
    let formData = {};

    if (!Object.keys(errMap).length) {
        formData = {
            "first_name": firstName, "last_name": lastName,
            "email": email, "phone_number": phoneNumber
        }
    } else {

        // Removing spinner
        unLoadingButtonW(btn, "UPDATE", "material-icons", "edit")
        displayEditError(errMap)
        return
    }


    $.ajax({
        url: `../../admin/staff/${id}/${csrf}`,
        data: formData,
        type: "PUT",
        success: successFunc,
        error: errorFunc,
    });


    function successFunc(result, status, xhr) {

        // Removing spinner
        unLoadingButtonW(btn, "UPDATE", "material-icons", "edit")
        $("#edit-member-error").removeClass("error").addClass("success").
            text("Successfully updated staff member profile").css("visibility", "visible")

    }

    function errorFunc(xhr, status, error) {

        // Removing spinner
        unLoadingButtonW(btn, "UPDATE", "material-icons", "edit")

        if (xhr.status == 400) {
            errMap = JSON.parse(xhr.responseText)
            displayEditError(errMap)

        } else {
            $("#edit-member-error").removeClass("success").addClass("error").
                text("Unable to update staff member profile").css("visibility", "visible")
        }

    }


}

function displayEditError(errMap) {

    for (let prop in errMap) {

        switch (prop) {
            case "first_name":
                $("#edit-fname-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "last_name":
                $("#edit-lname-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "email":
                $("#edit-email-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "phone_number":
                $("#edit-phone-error").text(errMap[prop]).css("visibility", "visible")
                break

            case "error":
                $("#edit-member-error").removeClass("success").addClass("error").
                    text("Unable to update staff member profile").css("visibility", "visible")
                break
        }
    }
}

function initEditChangePassword() {

    let id = $("#edit-member-id").val()

    $("#edit-password-modal input").val("")
    $("#edit-password-modal .error").text("")
    $("#edit-password-modal .error").css("visibility", "hidden");

    $("#edit-password-id").text(id)
    $("#edit-password-modal").modal("show")
}

function changeStaffPassword(btn) {

    let csrf = $("#csrf").val()
    let id = $("#edit-member-id").val()

    // cleaning up before processing
    $("#edit-password-modal .error").css("visibility", "hidden");
    $("#edit-password-modal .error").text("");

    // Adding loading spinner
    loading(btn)

    let adminPassword = $("#edit-admin-password").val()
    let newPassword = $("#edit-new-password").val()
    let newVerifyPassword = $("#edit-verify-password").val()

    if (adminPassword.match(/^\s*$/)) {
        unLoading(btn, "Change");
        return
    }


    let errMap = verifyPassword(newPassword, newVerifyPassword)
    let formData = {};

    if (!Object.keys(errMap).length) {
        formData = {
            "admin_password": adminPassword,
            "new_password": newPassword,
            "new_vPassword": newVerifyPassword
        }
    } else {

        // Removing spinner
        unLoading(btn, "Change")

        for (let prop in errMap) {
            if (prop == "vPassword") {
                $("#edit-vrpassword-error").text(errMap[prop]).css("visibility", "visible")
            } else {
                $("#edit-nwpassword-error").text(errMap[prop]).css("visibility", "visible")
            }
        }

        return
    }


    $.ajax({
        url: `../../admin/staff/${id}/profile/password/${csrf}`,
        data: formData,
        type: "PUT",
        success: successFunc,
        error: errorFunc,
    });


    function successFunc(result, status, xhr) {
        // Removing spinner
        unLoading(btn, "Change")
        $("#edit-password-modal").modal("hide")
    }

    function errorFunc(xhr, status, error) {

        // Removing spinner
        unLoading(btn, "Change")

        if (xhr.status == 400) {

            errMap = JSON.parse(xhr.responseText)

            for (let prop in errMap) {

                switch (prop) {

                    case "admin_password":
                        $("#edit-adpassword-error").text("Invalid Admin Password Used").css("visibility", "visible")
                        break

                    case "password":
                        if (errMap[prop] == "password does not match") {
                            $("#edit-vrpassword-error").text(errMap[prop]).css("visibility", "visible")
                        } else {
                            $("#edit-nwpassword-error").text(errMap[prop]).css("visibility", "visible")
                        }
                        break
                }
            }

        } else {
            $("#edit-modal-error").text("Unable to change staff member password").css("visibility", "visible")
        }

    }
}