// initAddJob is a function that initiate or display the add job window
function initAddJob(employer, postType) {

    hideAllExcept($("#add-section"));
    $('#add-due-date-tooltip').tooltip()

    // cleaning up before processing
    $("#add-section .error").css("visibility", "hidden");
    $("#add-section .error").text("");
    $("#add-job-body input").val("")
    $("#add-job-body input[name=gender]").prop("checked", false)
    $("#add-gender-male").val("M")
    $("#add-gender-female").val("F")
    $("#add-gender-both").val("B")
    $("#add-section textarea").val("")
    $("select option").remove()

    let dueDateInit = new Date(Date.now())
    dueDateInit.setMonth(dueDateInit.getMonth() + 6)
    $('#add-due-date').datetimepicker({
        format: 'M d, Y H:i',
        minDateTime: true,
        defaultSelect: false,
        value: dueDateInit,
    });

    $('#add-education-level').append($('<option>', {
        value: "",
        text: "Please Select",
        disabled: 'disabled',
        selected: true,
    }));

    $('#add-work-experience').append($('<option>', {
        value: "",
        text: "Please Select",
        disabled: 'disabled',
        selected: true,
    }));

    $('#add-post-type').append($('<option>', {
        value: "User",
        text: "User",
    }), $('<option>', {
        value: "Internal",
        text: "Internal",
    }), $('<option>', {
        value: "External",
        text: "External",
    }));

    // Setting the employer if provided
    if (employer != null && employer != undefined && employer != "") {
        $('#add-job-employer').val(employer)
    }

    changePostType(postType)
    removeBaseURI()
    addBaseURI("Add")

    $('.selectpicker').selectpicker('refresh')  // refreshing the select pickers
    fetchJobAttributes()

}

function changePostType(type) {

    $('.removable').hide()
    $('.removable-asterisk').hide()

    // cleaning up before processing
    $("#add-section .error").css("visibility", "hidden")
    $("#add-section .error").text("")

    // Setting default post type
    if (type != "User" && type != "Internal" && type != "External") {
        type = "User"
    }

    if (type == "User") {
        $('#add-job-employer-group').show(200)
        $('#add-contact-type-group').show(200)
        $('#add-job-type-group').show(200)
        $('#add-education-level-group').show(200)
        $('#add-work-experience-group').show(200)
        $('#add-job-sector-group').show(200)
        $('#add-gender-group').show(200)
        $('#add-due-date-group').show(200)
        $('.removable-asterisk').show()

        $($("#add-job-body .row")[4]).removeClass("d-flex justify-content-center")


    } else if (type == "Internal") {
        $('#add-job-employer-group').show(200)
        $('#add-contact-info-group').show(200)
        $('#add-job-type-group').show(200)
        $('#add-education-level-group').show(200)
        $('#add-work-experience-group').show(200)
        $('#add-job-sector-group').show(200)
        $('#add-gender-group').show(200)
        $('#add-due-date-group').show(200)
        $('.removable-asterisk').show()

        $($("#add-job-body .row")[4]).removeClass("d-flex justify-content-center")

    } else if (type == "External") {
        $('#add-job-employer-group').show(200)
        $('#add-contact-link-group').show(200)
        $('#add-job-type-group').show(200)
        $('#add-education-level-group').show(200)
        $('#add-work-experience-group').show(200)
        $('#add-job-sector-group').show(200)
        $('#add-due-date-group').show(200)

        $($("#add-job-body .row")[4]).addClass("d-flex justify-content-center")

    }

    // Setting the selected option
    onSelectOption("add-post-type", type)
}

function fetchJobAttributes() {

    let csrf = $("#csrf").val()

    $.ajax({
        url: `../../staff/job/init/add/${csrf}`,
        type: "GET",
        success: successFunc,
    });

    function successFunc(result, status, xhr) {
        let resultContainer = JSON.parse(result);
        let validJobSectors = resultContainer.ValidJobSectors
        let validJobTypes = resultContainer.ValidJobTypes
        let validEducationLevels = resultContainer.ValidEducationLevels
        let validWorkExperiences = resultContainer.ValidWorkExperiences
        let contactTypes = resultContainer.ContactTypes

        contactTypes.forEach(contactType => {
            $('#add-contact-type').append($('<option>', {
                value: contactType,
                text: contactType
            }));
        });

        $('#add-contact-type').selectpicker('refresh') // We separately refresh each selectors to avoid 'Nothing Selected' text of default values

        validJobTypes.forEach(validJobType => {
            $('#add-job-type').append($('<option>', {
                id: validJobType.ID,
                value: validJobType.Name,
                text: validJobType.Name == "Other" ? "Other (please specify)" : validJobType.Name
            }));
        });

        validEducationLevels.forEach(validEducationLevel => {
            $('#add-education-level').append($('<option>', {
                id: validEducationLevel.ID,
                value: validEducationLevel.Name,
                text: validEducationLevel.Name == "Other" ? "Other (please specify)" : validEducationLevel.Name
            }));
        });

        validWorkExperiences.forEach(validWorkExperience => {
            $('#add-work-experience').append($('<option>', {
                value: validWorkExperience,
                text: validWorkExperience
            }));
        });

        validJobSectors.forEach(validJobSector => {
            $('#add-job-sector').append($('<option>', {
                id: validJobSector.ID,
                value: validJobSector.Name,
                text: validJobSector.Name == "Other" ? "Other (please specify)" : validJobSector.Name
            }));
        });

        $('.selectpicker').selectpicker('refresh');

    }

}

function addJob(btn) {

    let csrf = $("#csrf").val()

    // cleaning up before processing
    $("#add-section .error").css("visibility", "hidden");
    $("#add-section .error").text("");
    $("#added-job-id").text("")
    $("#added-job-title").text("")
    $("#added-job-desc").text("")

    // Adding loading spinner
    loadButtonW(btn)

    let title = $("#add-job-title").val()
    let description = $("#add-job-desc").val()
    let employer = $("#add-job-employer").val()
    let contactType = $("#add-contact-type").val()
    let jobType = $("#add-job-type").val()
    let educationLevel = $("#add-education-level").val()
    let workExperience = $("#add-work-experience").val()
    let jobSector = $("#add-job-sector").val()
    let gender = $("#add-job-body input[name=gender]:checked").val()
    let contactInfo = $("#add-contact-info").val()
    let postType = $("#add-post-type").val()
    let link = $("#add-contact-link").val()
    let dueDate = $('#add-due-date').datetimepicker('getValue')

    let errMap = validateJobEntries(postType, title, description, employer, contactType,
        jobType, educationLevel, workExperience, jobSector, gender, contactInfo, link, dueDate)

    let formData = {};

    if (!Object.keys(errMap).length) {
        formData = {
            "title": title, "description": description, "employer": employer,
            "contact_type": contactType, "job_type": `${jobType}`, "education_level": educationLevel,
            "work_experience": workExperience, "job_sector": `${jobSector}`, "gender": gender,
            "contact_info": contactInfo, "post_type": postType, "link": link,
            "due_date": dueDate != "" && dueDate != null ? (new Date(dueDate)).toISOString() : "",
        }
    } else {

        // Removing spinner
        unLoadingButtonW(btn, "ADD ", "custom-icon icon-job-edit")

        displayAddError(errMap)
        return
    }

    $.ajax({
        url: `../../staff/job/${csrf}`,
        data: formData,
        type: "POST",
        success: successFunc,
        error: errorFunc,
    });

    function successFunc(result, status, xhr) {

        // Removing spinner
        unLoadingButtonW(btn, "ADD ", "custom-icon icon-job-edit")

        let addedJob = JSON.parse(xhr.responseText)

        $("#added-job-id").text(addedJob.ID)
        $("#added-job-title").text(addedJob.Title)
        $("#added-job-desc").text(addedJob.Description)

        $("#add-modal-error").hide()
        $("#add-modal-success").show()
        $("#add-job-modal").modal("show")

        // Cleaning up
        initAddJob(employer, postType)

    }

    function errorFunc(xhr, status, error) {

        // Removing spinner
        unLoadingButtonW(btn, "ADD ", "custom-icon icon-job-edit")

        if (xhr.status == 400) {
            errMap = JSON.parse(xhr.responseText)
            displayAddError(errMap)

        } else {
            $("#add-modal-success").hide()
            $("#add-modal-error").show()
            $("#add-modal-error-msg").text("Unable to add the provided job.")
            $("#add-job-modal").modal("show")
        }

    }

}

function displayAddError(errMap) {

    for (let prop in errMap) {

        switch (prop) {

            case "title":
                $("#add-title-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "description":
                $("#add-desc-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "employer":
                $("#add-employer-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "contact_type":
                $("#add-contact-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "type":
                $("#add-type-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "experience":
                $("#add-experience-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "education_level":
                $("#add-level-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "sector":
                $("#add-sector-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "gender":
                $("#add-gender-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "contact_info":
                $("#add-contact-info-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "link":
                $("#add-link-error").text(errMap[prop]).css("visibility", "visible")
                break
            case "due_date":
                $("#add-due-error").text(errMap[prop]).css("visibility", "visible")
                break
        }
    }
}