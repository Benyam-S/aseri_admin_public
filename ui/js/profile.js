// initUserProfile is a function that initiate or display the the user profile
function initUserProfile() {

    let csrf = $("#csrf").val()
    let id = $("#current-user-id").val()

    hideAllExcept($("#profile-section"));

    // cleaning up before processing
    $("#profile-error").removeClass("success").addClass("error")
    $("#profile-container .error").css("visibility", "hidden");
    $("#profile-container .error").text("");
    $("#user-profile-pic").css("background-image", `url(../../ui/images/user.svg)`)
    $("#profile-body input").val("")
    $("#profile-delete-collapser").attr("aria-expanded", "false")
    $("#delete-profile-info").removeClass("show")
    cleanProfileFrame()

    $("#profile-body").show();
    $("#profile-error-div").hide();

    let element1 = $("<div> </div>").addClass("internal").text("/");
    let element2 = $("<a> </a>").addClass("top-nav-link internal").text("Profile")

    // Removing base uri
    removeBaseURI()
    $("#uri-div").append(element1, element2)

    $.ajax({
        url: `../../staff/profile/${csrf}`,
        type: "GET",
        success: successFunc,
        error: errorFunc,
    });

    function successFunc(result, status, xhr) {

        let staffMember = JSON.parse(result);

        if (staffMember.ProfilePic) {
            $("#profile-remove-photo").show()
            $("#profile-hidden-pic").val(staffMember.ProfilePic)
            $("#user-profile-pic").css("background-image", `url(../../assets/profilepics/${staffMember.ProfilePic})`)
        } else {
            $("#profile-remove-photo").hide()
        }

        $("#profile-first-name").css("text-transform", "capitalize").val(staffMember.FirstName)
        $("#profile-last-name").css("text-transform", "capitalize").val(staffMember.LastName)
        $("#profile-id").val(staffMember.ID)
        $("#profile-password").val("********")
        $("#profile-role").css({ "text-transform": "uppercase" }).val(staffMember.Role)
        $("#profile-email").val(staffMember.Email)
        $("#profile-phone-number").val(staffMember.PhoneNumber.replace("+251", "0"))

    }

    function errorFunc(xhr, status, error) {

        $("#profile-body").hide();
        $("#profile-error-div").show();

    }

}

function cleanProfileFrame() {
    $("#profile-uploading-photo").hide()
    $("#profile-remove-photo").hide()
    $("#profile-reload-photo").hide()
    $("#profile-photo-error").text("").css("visibility", "hidden");
    $("#user-profile-pic > div").css("visibility", "hidden")
}

function profileUploadPhoto() {
    $("#profile-uploader").click()
}

function profilePreviewPhoto(uploader) {
    if (uploader.files && uploader.files[0]) {
        var imageFile = uploader.files[0];
        var reader = new FileReader();
        reader.onload = function (e) {
            $("#user-profile-pic").css("background-image", `url(${e.target.result})`);
        };

        reader.readAsDataURL(imageFile);

        // Uploading to the database
        let csrf = $("#csrf").val()
        let id = $("#current-user-id").val()

        let spinner = $("<span> </span>").addClass("spinner-border spinner-border-sm").
            attr({ "role": "status", "aria-hidden": true })
        $("#profile-uploader").val("")
        $("#profile-upload-photo").hide()
        $("#profile-remove-photo").hide()
        $("#user-profile-pic > div").css("visibility", "visible")
        $("#profile-uploading-photo").show().append(spinner)

        let formData = new FormData()
        formData.append("profile_pic", imageFile)

        $.ajax({
            url: `../../staff/profile/pic/${csrf}`,
            type: "POST",
            data: formData,
            processData: false,
            contentType: false,
            success: successFunc,
            error: errorFunc,
        });

        function successFunc(result, status, xhr) {

            let fileNameContainer = JSON.parse(xhr.responseText)

            cleanProfileFrame()
            $("#profile-uploading-photo").find(".spinner-border").remove()
            $("#profile-remove-photo").show()
            $("#profile-upload-photo").show()
            $("#profile-hidden-pic").val(fileNameContainer["file_name"])

            // Changing header user profile pic
            $("#navbar-top #profile-pic").css("background-image", `url(../../assets/profilepics/${fileNameContainer["file_name"]})`);

        }

        function errorFunc(xhr, status, error) {

            let err = ""
            if (xhr.status == 400) {
                let errMap = JSON.parse(xhr.responseText)
                err = errMap["error"]

            } else {
                err = "Unable To Upload Photo"
            }

            cleanProfileFrame()
            $("#profile-uploading-photo").find(".spinner-border").remove()
            $("#profile-reload-photo").show()
            $("#profile-upload-photo").show()
            $("#profile-photo-error").css("text-transform", "capitalize").text(err).css("visibility", "visible");
        }

    }
}

function profileRemovePhoto() {

    let csrf = $("#csrf").val()

    let spinner = $("<span> </span>").addClass("spinner-border spinner-border-sm").
        attr({ "role": "status", "aria-hidden": true })
    $("#profile-uploader").val("")
    $("#profile-upload-photo").hide()
    $("#profile-remove-photo").hide()
    $("#user-profile-pic > div").css("visibility", "visible")
    $("#profile-uploading-photo").show().append(spinner)


    if ($("#profile-hidden-pic").val()) {

        $.ajax({
            url: `../../staff/profile/pic/${csrf}`,
            type: "DELETE",
            success: successFunc,
            error: errorFunc,
        });

        function successFunc(result, status, xhr) {
            cleanProfileFrame()
            $("#profile-uploading-photo").find(".spinner-border").remove()
            $("#profile-upload-photo").show()
            $("#profile-hidden-pic").val("")
            $("#user-profile-pic").css("background-image", `url(../../ui/images/user.svg)`)

            // Changing header user profile pic
            $("#navbar-top #profile-pic").css("background-image", `url(../../ui/images/user.svg)`);

        }

        function errorFunc(xhr, status, error) {
            cleanProfileFrame()
            $("#profile-uploading-photo").find(".spinner-border").remove()
            $("#profile-upload-photo").show()
            $("#profile-remove-photo").show()
            $("#profile-photo-error").css("text-transform", "capitalize").text("Unable To Delete Profile Pic").css("visibility", "visible");
        }

    } else {
        cleanProfileFrame()
        $("#profile-upload-photo").show()
        $("#user-profile-pic").css("background-image", `url(../../ui/images/user.svg)`)
    }


}

function profileReloadPhoto() {

    cleanProfileFrame()
    $("#profile-upload-photo").show()

    let profilePic = $("#profile-hidden-pic").val()

    if (profilePic) {
        $("#profile-remove-photo").show()
        $("#user-profile-pic").css("background-image", `url(../../assets/profilepics/${profilePic})`)
    } else {
        $("#profile-remove-photo").hide()
        $("#user-profile-pic").css("background-image", `url(../../ui/images/user.svg)`)
    }
}

function updateProfile(btn) {

    let csrf = $("#csrf").val()

    // cleaning up before processing
    $("#profile-error").removeClass("success").addClass("error")
    $("#profile-container .error").css("visibility", "hidden");
    $("#profile-container .error").text("");

    // Adding loading spinner
    loadButtonW(btn)

    let firstName = $("#profile-first-name").val()
    let lastName = $("#profile-last-name").val()
    let email = $("#profile-email").val()
    let phoneNumber = $("#profile-phone-number").val()


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
        displayProfileError(errMap)
        return
    }


    $.ajax({
        url: `../../staff/profile/${csrf}`,
        data: formData,
        type: "PUT",
        success: successFunc,
        error: errorFunc,
    });


    function successFunc(result, status, xhr) {

        // Removing spinner
        unLoadingButtonW(btn, "UPDATE", "material-icons", "edit")
        $("#profile-error").removeClass("error").addClass("success").
            text("Successfully updated staff member profile").css("visibility", "visible")

        // Changing headers user name
        $("#navbar-top #user-name").text(firstName + " " + lastName)

    }

    function errorFunc(xhr, status, error) {

        // Removing spinner
        unLoadingButtonW(btn, "UPDATE", "material-icons", "edit")

        if (xhr.status == 400) {
            errMap = JSON.parse(xhr.responseText)
            displayProfileError(errMap)

        } else {
            $("#profile-error").removeClass("success").addClass("error").
                text("Unable to update staff member profile").css("visibility", "visible")
        }

    }


}

function displayProfileError(errMap) {

    for (let prop in errMap) {

        switch (prop) {
            case "first_name":
                $("#profile-fname-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "last_name":
                $("#profile-lname-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "email":
                $("#profile-email-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "phone_number":
                $("#profile-phone-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "error":
                $("#profile-member-error").removeClass("success").addClass("error").
                    text("Unable to update staff member profile").css("visibility", "visible")
                break
        }
    }
}

function initProfileChangePassword() {

    $("#profile-password-modal input").val("")
    $("#profile-password-modal .error").text("")
    $("#profile-password-modal .error").css("visibility", "hidden");

    $("#profile-password-modal").modal("show")
}

function changeUserPassword(btn) {


    let csrf = $("#csrf").val()

    // cleaning up before processing
    $("#profile-password-modal .error").css("visibility", "hidden");
    $("#profile-password-modal .error").text("");

    // Adding loading spinner
    loading(btn)

    let currentPassword = $("#profile-current-password").val()
    let newPassword = $("#profile-new-password").val()
    let newVerifyPassword = $("#profile-verify-password").val()


    if (currentPassword.match(/^\s*$/)) {
        unLoading(btn, "Change");
        return
    }


    let errMap = verifyPassword(newPassword, newVerifyPassword)
    let formData = {};

    if (!Object.keys(errMap).length) {
        formData = {
            "current_password": currentPassword,
            "new_password": newPassword,
            "new_vPassword": newVerifyPassword
        }
    } else {

        // Removing spinner
        unLoading(btn, "Change")

        for (let prop in errMap) {
            if (prop == "vPassword") {
                $("#profile-vrpassword-error").text(errMap[prop]).css("visibility", "visible")
            } else {
                $("#profile-nwpassword-error").text(errMap[prop]).css("visibility", "visible")
            }
        }


        return
    }


    $.ajax({
        url: `../../staff/profile/password/${csrf}`,
        data: formData,
        type: "PUT",
        success: successFunc,
        error: errorFunc,
    });


    function successFunc(result, status, xhr) {
        // Removing spinner
        unLoading(btn, "Change")
        $("#profile-password-modal").modal("hide")
    }

    function errorFunc(xhr, status, error) {

        // Removing spinner
        unLoading(btn, "Change")

        if (xhr.status == 400) {

            errMap = JSON.parse(xhr.responseText)

            for (let prop in errMap) {

                switch (prop) {

                    case "current_password":
                        $("#profile-crpassword-error").text("Invalid Password Used").css("visibility", "visible")
                        break

                    case "password":
                        if (errMap[prop] == "password does not match") {
                            $("#profile-vrpassword-error").text(errMap[prop]).css("visibility", "visible")
                        } else {
                            $("#profile-nwpassword-error").text(errMap[prop]).css("visibility", "visible")
                        }
                        break
                }
            }

        } else {
            $("#profile-modal-error").text("Unable to change staff member password").css("visibility", "visible")
        }

    }
}

function initDeleteBody() {
    // cleaning up before processing
    $("#profile-delete-div .error").css("visibility", "hidden");
    $("#profile-delete-div .error").text("");

    $("#delete-profile-password").val("")
}

function deleteProfile(btn) {

    let csrf = $("#csrf").val()

    // cleaning up before processing
    $("#profile-delete-div .error").css("visibility", "hidden");
    $("#profile-delete-div .error").text("");

    // Adding loading spinner
    loading(btn)

    let currentPassword = $("#delete-profile-password").val()

    if (currentPassword.match(/^\s*$/)) {
        unLoading(btn, "Delete");
        return
    }

    $.ajax({
        url: `../../staff/profile/${csrf}`,
        data: { "password": currentPassword },
        type: "POST",
        success: successFunc,
        error: errorFunc,
    });


    function successFunc(result, status, xhr) {
        // Removing spinner
        unLoading(btn, "Delete")

        window.location.replace("../../staff/login")

    }

    function errorFunc(xhr, status, error) {

        // Removing spinner
        unLoading(btn, "Delete")

        if (xhr.status == 400) {
            $("#delete-prpassword-error").text("Invalid Password Used!").css("visibility", "visible")
        } else {
            $("#delete-prpassword-error").text("Unable to delete staff member password").css("visibility", "visible")
        }
    }

}