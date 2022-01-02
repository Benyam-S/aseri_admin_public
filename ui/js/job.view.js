// viewJob is a function that views selected job details
function viewJob(id) {

    let csrf = $("#csrf").val()

    $("#view-job-title").text("");
    $("#job-detail-container input").val("");
    $("#view-job-desc").val("")
    $("#view-jobgender-both").prop("checked", false)
    $("#view-jobgender-both").prop("disabled", true)
    $("#view-jobgender-male").prop("checked", false)
    $("#view-jobgender-male").prop("disabled", true)
    $("#view-jobgender-female").prop("checked", false)
    $("#view-jobgender-female").prop("disabled", true)
    $("#view-contact-type").removeAttr('class')
    $("#job-detail-container .tooltiptext").remove();
    $("#view-section .modal-footer").hide()
    $("#view-modal-spacer").show()

    $.ajax({
        url: `../../staff/job/${id}/${csrf}`,
        type: "GET",
        success: successFunc,
        error: errorFunc,
    });

    function successFunc(result, status, xhr) {

        let job = JSON.parse(result);

        $("#job-detail-container").show();
        $("#view-error-div").hide();

        setViewDisplay(job)

    }

    function errorFunc(xhr, status, error) {

        $("#job-detail-container").hide();
        $("#view-section .modal-dialog").addClass("width-600")
        $("#view-error-div").show();

    }

}

function setIDStatus(status, element) {

    if (status == "P") {
        element.css({
            "border-color": "orange"
        })
    } else if (status == "O") {
        element.css({
            "border-color": "#0073aa"
        })
    } else if (status == "C") {
        element.css({
            "border-color": "#1ab31a"
        })
    } else if (status == "D") {
        element.css({
            "border-color": "tomato"
        })
    }

    element.css({ "border-width": "1.5px" })
}

function setViewDisplay(job) {

    $('.removable').hide()

    if (job.PostType == "User") {
        $('#view-contact-type-group').show()
        $('#view-job-employer-group').show()
        $('#view-job-type-group').show()
        $('#view-job-education-group').show()
        $('#view-job-experience-group').show()
        $('#view-job-sector-group').show()
        $('#view-job-gender-group').show()
        $('#view-due-date-group').show()

        $($("#job-detail-container .row")[4]).removeClass("d-flex justify-content-center")

    } else if (job.PostType == "Internal") {
        $('#view-contact-info-group').show()
        $('#view-job-employer-group').show()
        $('#view-job-type-group').show()
        $('#view-job-education-group').show()
        $('#view-job-experience-group').show()
        $('#view-job-sector-group').show()
        $('#view-job-gender-group').show()
        $('#view-due-date-group').show()

        $($("#job-detail-container .row")[4]).removeClass("d-flex justify-content-center")

    } else if (job.PostType == "External") {
        $('#view-job-employer-group').show()
        $('#view-job-link-group').show()
        $('#view-job-type-group').show()
        $('#view-job-education-group').show()
        $('#view-job-experience-group').show()
        $('#view-job-sector-group').show()
        $('#view-due-date-group').show()

        $($("#job-detail-container .row")[4]).addClass("d-flex justify-content-center")

    }

    // Tooltip incase of visibility problem
    let jobID = job.ID ? $("<span></span>").text(job.ID).addClass("tooltiptext") : null;
    let employer = job.Employer ? $("<span></span>").text(job.Employer).addClass("tooltiptext") : null;
    let type = job.Type ? $("<span></span>").text(reConstructor(job.Type)).addClass("tooltiptext") : null;
    let education = job.EducationLevel ? $("<span></span>").text(job.EducationLevel).addClass("tooltiptext") : null;
    let experience = job.Experience ? $("<span></span>").text(job.Experience).addClass("tooltiptext") : null;
    let sector = job.Sector ? $("<span></span>").text(reConstructor(job.Sector)).addClass("tooltiptext") : null;
    let link = job.Link ? $("<span></span>").text(job.Link).addClass("tooltiptext") : null;
    let postType = job.PostType ? $("<span></span>").text(job.PostType).addClass("tooltiptext") : null;

    let createdAt = moment(job.CreatedAt).year() > 10 ?
        $("<span></span>").text(moment(job.CreatedAt).format("dddd, MMMM Do YYYY, h:mm:ss a")).
            addClass("tooltiptext") : null;

    let dueDate = moment(job.DueDate).year() > 10 ?
        $("<span></span>").text(moment(job.DueDate).format("dddd, MMMM Do YYYY, h:mm:ss a")).
            addClass("tooltiptext") : null;


    $("#view-job-id").parent().addClass("tooltip-div").append(jobID)
    $("#view-job-created").parent().addClass("tooltip-div").append(createdAt)
    $("#view-job-employer").parent().addClass("tooltip-div").append(employer)
    $("#view-job-type").parent().addClass("tooltip-div").append(type)
    $("#view-job-education").parent().addClass("tooltip-div").append(education)
    $("#view-job-experience").parent().addClass("tooltip-div").append(experience)
    $("#view-job-sector").parent().addClass("tooltip-div").append(sector)
    $("#view-job-link").parent().addClass("tooltip-div").append(link)
    $("#view-due-date").parent().addClass("tooltip-div").append(dueDate)
    $("#view-post-type").parent().addClass("tooltip-div").append(postType)

    $("#view-job-title").text(job.Title.toTitleCase())
    $("#view-job-desc").val(job.Description)
    setIDStatus(job.Status, $("#view-job-id").val(job.ID))
    $("#view-job-created").val(
        moment(job.CreatedAt).year() > 10 ?
            moment(job.CreatedAt).format("dddd, MMMM Do YYYY, h:mm:ss a") : null)
    $("#view-due-date").val(
        moment(job.DueDate).year() > 10 ?
            moment(job.DueDate).format("dddd, MMMM Do YYYY, h:mm:ss a") : null)
    $("#view-job-employer").val(job.Employer)
    $("#view-job-type").val(reConstructor(job.Type))
    $("#view-job-education").val(job.EducationLevel)
    $("#view-job-experience").val(job.Experience)
    $("#view-job-sector").val(reConstructor(job.Sector))
    $("#view-job-link").val(job.Link)
    $("#view-contact-info").text(job.ContactInfo)
    $("#view-post-type").val(job.PostType)

    if (job.ContactType == "Send CV") {
        $("#view-contact-type").addClass("custom-icon icon-cv")
    } else if (job.ContactType == "Via Telegram Account") {
        $("#view-contact-type").addClass("custom-icon icon-telegram-black")
    }

    if (job.Gender == "M") {
        $("#view-jobgender-male").prop("checked", true)
        $("#view-jobgender-male").prop("disabled", false)

        $("#view-jobgender-female").prop("checked", false)
        $("#view-jobgender-female").prop("disabled", true)

        $("#view-jobgender-both").prop("checked", false)
        $("#view-jobgender-both").prop("disabled", true)


    } else if (job.Gender == "F") {
        $("#view-jobgender-female").prop("checked", true)
        $("#view-jobgender-female").prop("disabled", false)

        $("#view-jobgender-male").prop("checked", false)
        $("#view-jobgender-male").prop("disabled", true)

        $("#view-jobgender-both").prop("checked", false)
        $("#view-jobgender-both").prop("disabled", true)

    } else if (job.Gender == "B") {
        $("#view-jobgender-both").prop("checked", true)
        $("#view-jobgender-both").prop("disabled", false)

        $("#view-jobgender-male").prop("checked", false)
        $("#view-jobgender-male").prop("disabled", true)

        $("#view-jobgender-female").prop("checked", false)
        $("#view-jobgender-female").prop("disabled", true)

    }

    // Adding approve and decline buttons
    if (job.Status == "P") {
        $("#view-section .modal-footer").show()
        $("#view-modal-spacer").hide()
    }
}