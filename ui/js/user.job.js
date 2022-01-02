// CTPAGEVALUE is a variable that can be altered by the perviousJobs and nextJobs buttons
var CTPAGEVALUE = 0;
// CTTOTALPAGECOUNT is a variable that holds the current user jobs result total page
var CTTOTALPAGECOUNT = 0;


function initViewUserJobs(userID) {

    $("#user-job-table").show();
    $("#uj-noresult-div").hide();
    $("#view-ujerror-div").hide();
    $("#user-job .modal-dialog").removeClass("width-600")

    CTPAGEVALUE = 0;
    CTTOTALPAGECOUNT = 0;

    $("#hidden-user-id").val(userID)

    allUserJobs(CTPAGEVALUE)
}

function allUserJobs(page) {

    let csrf = $("#csrf").val()
    let userID = $("#hidden-user-id").val()

    $.ajax({
        url: `../../staff/jobs/search/${csrf}`,
        data: { "key": userID, page },
        type: "GET",
        success: successFunc,
        error: errorFunc,
    });

    function successFunc(result, status, xhr) {

        let resultContainer = JSON.parse(result);
        let currentPage = resultContainer.CurrentPage;
        let pageCount = resultContainer.PageCount;
        let jobs = resultContainer.Result;

        if (pageCount <= 1) {
            $("#uj-pagination-div").hide();
        } else {
            $("#uj-pagination-div").show();

            if (currentPage > 0 && currentPage < pageCount) {
                $("#uj-prev-icon").show();
            } else {
                $("#uj-prev-icon").hide();
            }

            if (currentPage >= 0 && currentPage < pageCount - 1) {
                $("#uj-next-icon").show();
            } else {
                $("#uj-next-icon").hide();
            }
        }

        CTTOTALPAGECOUNT = resultContainer.PageCount
        let tableBody = $("#uj-table-body").empty();
        $("#view-ujerror-div").hide();

        if (jobs.length == 0) {
            $("#user-job .modal-dialog").addClass("width-600")
            $("#user-job-table").hide();
            $("#uj-noresult-div").show();
            $("#uj-noresult-div-msg").text(`Sorry we couldn't find any job.`);
            return;
        } else {
            $("#uj-noresult-div").hide();
            $("#user-job-table").show();
        }

        appendJobs(tableBody, jobs);
    }

    function errorFunc(xhr, status, error) {

        $("#user-job .modal-dialog").addClass("width-600")

        $("#uj-table-body").empty();
        $("#user-job-table").hide();
        $("#uj-pagination-div").hide();
        $("#uj-noresult-div").hide();
        $("#view-ujerror-div").show();
    }

}

function nextJobs() {

    if (CTPAGEVALUE <= 0) {
        allUserJobs(1)
        ++CTPAGEVALUE

    } else if (CTPAGEVALUE < CTTOTALPAGECOUNT - 1) {
        allUserJobs(++CTPAGEVALUE)
    }
}

function previousJobs() {

    if (CTPAGEVALUE > 0) {
        allUserJobs(--CTPAGEVALUE)
    }
}

// appendJobs is a function that appends set of jobs from the result container
function appendJobs(tableBody, jobs) {

    jobs.forEach(
        (job) => {

            let reducedTitle = reduceTo(job.Title, 10, true)

            let id = $("<span></span>").text(job.ID)
            let title = $("<span></span>").text(reducedTitle)
            let viewIcon = $("<span></span>").addClass("material-icons").text("visibility")
            let type = $("<span></span>").text(job.Type)
            let sector = $("<span></span>").css("text-transform", "capitalize").text(job.Sector)
            let status = $("<span></span>").css("text-transform", "capitalize").
                text(job.Status).addClass("status-cell").
                attr("id", "status-" + job.ID)


            let row = $("<tr> </tr>")

            let col1 = $("<td> </td>").append(id)
            let col2 = $("<td> </td>").append(title)
            let col3 = $("<td> </td>").append(viewIcon)
            let col4 = $("<td> </td>").append(type)
            let col5 = $("<td> </td>").append(sector)
            let col6 = $("<td> </td>").append(setStatus(status))

            // Setting onclick functions
            viewIcon.click(() =>
                navigate('../../staff/job', `view-section@${job.ID}`, "_blank")
            )

            row.append(col1, col2, col3, col4, col5, col6)
            tableBody.append(row)

        }
    )
}