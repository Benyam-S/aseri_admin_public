$(document).ready(onload);

function onload() {

    // Showing all job attributes
    hideAllExcept($("#result-section"))
    allAttributes();

    // Keybindings
    $(document).on("keyup", (e) => {
        if (e.altKey && e.which == 80) {
            initUserProfile()
        }
    })

}

function allAttributes() {

    let csrf = $("#csrf").val()

    $.ajax({
        url: `../../staff/attributes/${csrf}`,
        type: "GET",
        success: successFunc,
        error: errorFunc
    });

    function successFunc(result, status, xhr) {

        // Clearing table errors
        $("#result-sector-section").show();
        $("#result-both-section").show();
        $("#result-error-div").hide();

        let resultContainer = JSON.parse(result);
        let jobSectors = resultContainer.JobSectors;
        let jobTypes = resultContainer.JobTypes;
        let educationLevels = resultContainer.EducationLevels;

        let sectorTableBody = $("#table-sector-body").empty();
        let typeTableBody = $("#table-type-body").empty();
        let levelTableBody = $("#table-level-body").empty();

        if (jobSectors.length == 0) {
            $("#noresult-sector-div").show();
            $("#noresult-sector-div-msg").text(`Sorry we couldn't find any job sectors.`);
        } else {
            $("#noresult-sector-div").hide();
            appendJobAttributes(sectorTableBody, jobSectors, "job sectors");
        }

        if (jobTypes.length == 0) {
            $("#noresult-type-div").show();
            $("#noresult-type-div-msg").text(`Sorry we couldn't find any job types.`);
        } else {
            $("#noresult-type-div").hide();
            appendJobAttributes(typeTableBody, jobTypes, "job types");
        }

        if (educationLevels.length == 0) {
            $("#noresult-level-div").show();
            $("#noresult-level-div-msg").text(`Sorry we couldn't find any education levels.`);
        } else {
            $("#noresult-level-div").hide();
            appendJobAttributes(levelTableBody, educationLevels, "education levels");
        }

    }

    function errorFunc(xhr, status, error) {
        $("#result-sector-section").hide();
        $("#result-both-section").hide();
        $("#result-error-div").css({"margin-top":"90px"}).show();
    }

}

// appendJobAttributes is a function that appends set of job attributes from the result container
function appendJobAttributes(tableBody, attributes, category) {

    attributes.forEach(
        (attribute) => {

            let id = $("<span></span>").text(attribute.ID)
            let name = $("<span></span>").text(attribute.Name)
            let deleteIcon = $("<span></span>").addClass("material-icons").text("delete")

            let row = $("<tr> </tr>")

            let col1 = $("<td> </td>").append(id)
            let col2 = $("<td> </td>").append(name)
            let col3 = $("<td> </td>").append(deleteIcon)

            deleteIcon.attr({ "data-toggle": "modal", "data-target": "#delete-attribute" })
            deleteIcon.click(() =>
                initDeleteAttribute(attribute.ID, category)
            )

            row.append(col1, col2, col3)

            tableBody.append(row)
        }
    )
}

function hideAllExcept(section) {

    $("#result-section").hide()
    $("#profile-section").hide()

    section.show()
}