// initAddStaff is a function that initiate or display the add staff member window
function initAddStaff() {

    hideAllExcept($("#add-section"));

    // cleaning up before processing
    $(".error").css("visibility", "hidden");
    $(".error").text("");
    $("#add-member-body input").val("")
    $("#add-member-role").selectpicker('val', 'none');

    removeBaseURI()
    addBaseURI("Add")

}

function addStaff(btn) {

    let csrf = $("#csrf").val()

    // cleaning up before processing
    $(".error").css("visibility", "hidden");
    $(".error").text("");
    $("#added-member-id").text("")
    $("#added-member-name").text("")
    $("#added-member-role").text("")

    // Adding loading spinner
    loadButtonW(btn)

    let firstName = $("#add-first-name").val()
    let lastName = $("#add-last-name").val()
    let email = $("#add-member-email").val()
    let phoneNumber = $("#add-phone-number").val()
    let role = $("#add-member-role").find(":selected").val()

    let password = $("#add-member-password").val()
    let vPassword = $("#add-verify-password").val()

    let errMap1 = validateStaffProfile(firstName, lastName, email, phoneNumber)
    let errMap2 = verifyPassword(password, vPassword)
    let errMap = { ...errMap1, ...errMap2 }

    let formData = {};

    if (!Object.keys(errMap).length) {
        formData = {
            "first_name": firstName, "last_name": lastName,
            "email": email, "phone_number": phoneNumber,
            "password": password, "vpassword": vPassword
        }
    } else {

        // Removing spinner
        unLoadingButtonW(btn, "ADD ", "custom-icon icon-admin")

        displayAddError(errMap)
        return
    }

    if (role == "staff") {
        $.ajax({
            url: `../../admin/staff/${csrf}`,
            data: formData,
            type: "POST",
            success: successStaffFunc,
            error: errorFunc,
        });

    } else if (role == "admin") {

        $.ajax({
            url: `../../admin/register/init/${csrf}`,
            data: formData,
            type: "POST",
            success: successInitFunc,
            error: errorFunc,
        });

        // cleaning up
        $("#add-admin-nonce").val("")
        $("#add-admin-otp").val("")
        $("#add-otp-error").css("visibility", "hidden")

    } else {

        // Removing spinner
        unLoadingButtonW(btn, "ADD ", "custom-icon icon-admin")

        $("#add-role-error").text("select role").css("visibility", "visible")
        return
    }

    function successStaffFunc(result, status, xhr) {

        // Removing spinner
        unLoadingButtonW(btn, "ADD ", "custom-icon icon-admin")

        let addedMember = JSON.parse(xhr.responseText)

        $("#added-member-id").text(addedMember.ID)
        $("#added-member-name").text(addedMember.FirstName + " " + addedMember.LastName)
        $("#added-member-role").css("text-transform", "uppercase").text(addedMember.Role)

        $("#add-modal-error").hide()
        $("#add-modal-admin").hide()
        $("#add-member-modal .modal-footer").hide()
        $("#add-modal-success").show()
        $("#add-member-modal").modal("show")

        // Cleaning up
        initAddStaff()

    }

    function errorFunc(xhr, status, error) {

        // Removing spinner
        unLoadingButtonW(btn, "ADD ", "custom-icon icon-admin")

        if (xhr.status == 400) {
            errMap = JSON.parse(xhr.responseText)
            displayAddError(errMap)

        } else {
            $("#add-modal-success").hide()
            $("#add-modal-admin").hide()
            $("#add-member-modal .modal-footer").hide()
            $("#add-modal-error").show()
            $("#add-modal-error-msg").text("Unable to add the provided staff member.")
            $("#add-member-modal").modal("show")
        }

    }

    // Initiating Admin Adding process
    function successInitFunc(result, status, xhr) {
        // Removing spinner
        unLoadingButtonW(btn, "ADD ", "custom-icon icon-admin")

        let nonceContainer = JSON.parse(xhr.responseText)

        $("#add-admin-nonce").val(nonceContainer.nonce)

        $("#add-modal-error").hide()
        $("#add-modal-success").hide()
        $("#add-modal-admin").show()
        $("#add-member-modal .modal-footer").show()
        $("#add-member-modal").modal("show")
    }

}

function addAdmin(btn) {

    let csrf = $("#csrf").val()

    // Adding loading spinner
    loading(btn)

    let otp = $("#add-admin-otp").val()
    let nonce = $("#add-admin-nonce").val()


    if (otp.length == 0) {
        // Removing spinner
        unLoading(btn, "Submit")
    }

    $.ajax({
        url: `../../admin/register/finish/${csrf}`,
        data: { otp, nonce },
        type: "POST",
        success: successFinishFunc,
        error: errorFinishFunc,
    });

    function successFinishFunc(result, status, xhr) {

        // Removing spinner
        unLoading(btn, "Submit")

        let addedMember = JSON.parse(xhr.responseText)

        $("#added-member-id").text(addedMember.ID)
        $("#added-member-name").text(addedMember.FirstName + " " + addedMember.LastName)
        $("#added-member-role").css("text-transform", "uppercase").text(addedMember.Role)

        $("#add-modal-error").hide()
        $("#add-modal-admin").hide()
        $("#add-member-modal .modal-footer").hide()
        $("#add-modal-success").show()
        $("#add-member-modal").modal("show")

        // Cleaning up
        $(".error").css("visibility", "hidden");
        $(".error").text("");
        $("#add-member-body input").val("")
        $("#add-member-role").selectpicker('val', 'none');

    }

    function errorFinishFunc(xhr, status, error) {

        // Removing spinner
        unLoading(btn, "Submit")

        if (xhr.status == 400) {
            $("#add-otp-error").css("visibility", "visible").text("Incorrect OTP Used!")

        } else {
            $("#add-modal-success").hide()
            $("#add-modal-admin").hide()
            $("#add-member-modal .modal-footer").hide()
            $("#add-modal-error").show()
            $("#add-modal-error-msg").text("Unable to add the provided staff member.")
            $("#add-member-modal").modal("show")
        }

    }


}

function displayAddError(errMap) {

    for (let prop in errMap) {

        switch (prop) {
            case "first_name":
                $("#add-fname-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "last_name":
                $("#add-lname-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "email":
                $("#add-email-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "phone_number":
                $("#add-phone-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "password":
                if (errMap[prop] == "password does not match") {
                    $("#add-vpassword-error").text(errMap[prop]).css("visibility", "visible")
                } else {
                    $("#add-password-error").text(errMap[prop]).css("visibility", "visible")
                }
                break
            case "vPassword":
                $("#add-vpassword-error").text(errMap[prop]).css("visibility", "visible")
                break
        }
    }
}