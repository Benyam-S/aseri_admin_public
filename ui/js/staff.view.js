// viewStaffMember is a function that views selected staff member profile
function viewStaffMember(id) {

    let csrf = $("#csrf").val()

    $("#member-profile-container input").val("");
    $("#member-profile-container .tooltiptext").remove();
    $("#view-profile-pic").css("background-image", `url(../../ui/images/user.svg)`)
    $("#view-section .modal-dialog").removeClass("width-600")

    $.ajax({
        url: `../../admin/staff/${id}/${csrf}`,
        type: "GET",
        success: successFunc,
        error: errorFunc,
    });

    function successFunc(result, status, xhr) {

        let staffMember = JSON.parse(result);

        $("#member-profile-container").show();
        $("#view-error-div").hide();

        if (staffMember.ProfilePic) {
            $("#view-profile-pic").css("background-image", `url(../../assets/profilepics/${staffMember.ProfilePic})`)
        }

        // Tooltip incase of invisibility problem
        let firstName = staffMember.FirstName ? $("<span></sapn>").text(staffMember.FirstName).addClass("tooltiptext") : null;
        let lastName = staffMember.LastName ? $("<span></sapn>").text(staffMember.LastName).addClass("tooltiptext") : null;
        let id = staffMember.ID ? $("<span></sapn>").text(staffMember.ID).addClass("tooltiptext") : null;
        let role = staffMember.Role ? $("<span></sapn>").text(staffMember.Role).addClass("tooltiptext") : null;
        let email = staffMember.Email ? $("<span></sapn>").text(staffMember.Email).addClass("tooltiptext") : null;
        let phoneNumbber = staffMember.PhoneNumber ? $("<span></sapn>").text(staffMember.PhoneNumber).addClass("tooltiptext") : null;

        $("#view-first-name").parent().addClass("tooltip-div").append(firstName)
        $("#view-last-name").parent().addClass("tooltip-div").append(lastName)
        $("#view-member-id").parent().addClass("tooltip-div").append(id)
        $("#view-member-role").parent().addClass("tooltip-div").append(role)
        $("#view-member-email").parent().addClass("tooltip-div").append(email)
        $("#view-phone-number").parent().addClass("tooltip-div").append(phoneNumbber)

        $("#view-first-name").css("text-transform", "capitalize").val(staffMember.FirstName)
        $("#view-last-name").css("text-transform", "capitalize").val(staffMember.LastName)
        $("#view-member-id").val(staffMember.ID)
        $("#view-member-role").css({ "text-transform": "uppercase", "font-size": "14px" }).val(staffMember.Role)
        $("#view-member-email").val(staffMember.Email)
        $("#view-phone-number").val(staffMember.PhoneNumber.replace("+251", "0"))

    }

    function errorFunc(xhr, status, error) {

        $("#member-profile-container").hide();
        $("#view-section .modal-dialog").addClass("width-600")
        $("#view-error-div").show();

    }

}