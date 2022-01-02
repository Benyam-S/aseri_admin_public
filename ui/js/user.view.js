// viewUser is a function that views selected user profile
function viewUser(id, viewType) {

    let csrf = $("#csrf").val()
    let url = `../../staff/user/${id}/${csrf}`
    let errorText = "Unable to find selected user profile."

    $("#user-profile-container input").val("");
    $("#user-profile-container .tooltiptext").remove();
    $("#view-section .modal-dialog").removeClass("width-600")

    $("#user-profile-container").show();
    $("#view-uerror-div").hide();
    $("#view-user-notfound").hide();
    $("#user-profile-title").text("User Profile")

    if (viewType == "job-initiator") {
        url = `../../staff/job/initiator/${id}/${csrf}`
        $("#user-profile-title").text("Initiator Profile")
        errorText = "Unable to find selected initiator profile."

    }else if (viewType == "job-employer") {
        $("#user-profile-title").text("Employer Profile")
        errorText = "Unable to find selected employer profile."
    }

    $.ajax({
        url: url,
        type: "GET",
        success: successFunc,
        error: errorFunc,
    });

    function successFunc(result, status, xhr) {

        let user = JSON.parse(result);

        $("#view-uerror-div").hide(200);
        $("#view-user-notfound").hide(200);

        // Tooltip incase of invisibility problem
        let userName = user.UserName ? $("<span></span>").text(user.UserName).addClass("tooltiptext") : null;
        let id = user.ID ? $("<span></span>").text(user.ID).addClass("tooltiptext") : null;
        let phoneNumber = user.PhoneNumber ? $("<span></span>").text(user.PhoneNumber).addClass("tooltiptext") : null;
        let category = user.Category ? $("<span></span>").text(user.Category).addClass("tooltiptext") : null;

        $("#view-user-name").parent().addClass("tooltip-div").append(userName)
        $("#view-user-id").parent().addClass("tooltip-div").append(id)
        $("#view-user-phone").parent().addClass("tooltip-div").append(phoneNumber)
        $("#view-user-category").parent().addClass("tooltip-div").append(category)

        $("#view-user-name").css("text-transform", "capitalize").val(user.UserName)
        $("#view-user-id").val(user.ID)
        $("#view-user-phone").val(user.PhoneNumber.replace("+251", "0"))

        if (user.Category == "Aseri") {
            $("#view-user-category").val("አሰሪ")
        } else if (user.Category == "Agent") {
            $("#view-user-category").val("Agent")
        } else if (user.Category == "JobSeeker") {
            $("#view-user-category").val("Job Seeker")
        } else {
            $("#view-user-category").val(user.Category)
        }

    }

    function errorFunc(xhr, status, error) {

        $("#user-profile-container").hide(200);
        $("#view-section .modal-dialog").addClass("width-600")

        if (xhr.status == 400) {
            $("#view-uerror-div").hide(200);
            $("#view-user-notfound").show(300);
            $("#view-user-notfound-msg").text(errorText)
        } else {
            $("#view-user-notfound").hide(200);
            $("#view-uerror-div").show(300);
        }

    }

}