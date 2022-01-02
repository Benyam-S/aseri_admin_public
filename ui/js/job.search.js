// CURRENTPAGEVALUE is a variable that can be altered by the pervious and next buttons
var CURRENTPAGEVALUE = 0;
// CURRENTTOTALPAGECOUNT is a variable that holds the current result total page
var CURRENTTOTALPAGECOUNT = 0;
// CURRENTPAGEMODE is a variable that determines whether the user is on search mode or preview all mode
var CURRENTPAGEMODE = "all";
// CURRENTSTATUS is a variable that holds the current status type
var CURRENTSTATUS = "Any";
// CURRENTSEARCHKEY is variable that holds the current search term
var CURRENTSEARCHKEY = "";


$(document).ready(onload);

function onload() {

    // Marking job link
    $("#job-link").addClass("active-link")

    // Adding internal options
    addInternalOptions()

    // Keybindings
    $(document).on("keyup", (e) => {
        if (e.altKey && e.which == 65) {
            $("#view-all-btn").click()
        } else if (e.altKey && e.which == 78) {
            initAddJob()
        } else if (e.altKey && e.which == 80) {
            initUserProfile()
        } else if (e.altKey && e.which == 83) {
            initSearchJob()
        }
    })

    let displaySection = $("#display-section").val()
    if (displaySection != undefined && displaySection != null && displaySection) {
        if (displaySection.includes("add-section")) {
            let employer = displaySection.replace("add-section@", "")
            initAddJob(employer, "User")
            return
        } else if (displaySection.includes("view-section")) {

            let identifier = displaySection.replace("view-section@", "")
            $("#job-detail").modal("show")
            viewJob(identifier)

            // Adding search link to the navigation bar
            removeBaseURI()
            addBaseURI("Search")

            // Showing searching for the job
            hideAllExcept($("#search-section"))
            $("#search-bar").val(identifier)
            search(document.getElementById("search-button"))

            // Selecting default status
            $("#status-menu p").removeClass("selected-status")
            $(`#status-all`).addClass("selected-status")

            return
        }
    }

    $('.selectpicker').selectpicker() // initialization of select picker

    // Registering function for add job
    $('#add-post-type').on('change', function () {
        var selected = $('#add-post-type option:selected').val()
        changePostType(selected)
    })

    // Adding search link to the navigation bar
    removeBaseURI()
    addBaseURI("Search")

    // Showing all jobs
    hideAllExcept($("#search-section"));
    allJobs(CURRENTSTATUS, CURRENTPAGEVALUE);

    // Selecting default status
    $("#status-menu p").removeClass("selected-status")
    $(`#status-all`).addClass("selected-status")

}

// initSearchJob is a function that initiate or display job search window
function initSearchJob() {

    hideAllExcept($("#search-section"));

    removeBaseURI()
    addBaseURI("Search")

    // Resetting global variables
    CURRENTPAGEVALUE = 0
    CURRENTTOTALPAGECOUNT = 0
    CURRENTPAGEMODE = "all"
    CURRENTSTATUS = "Any"

    // Showing all jobs
    allJobs(CURRENTSTATUS, CURRENTPAGEVALUE);

    // Selecting default status
    $("#search-bar").val("")
    $("#status-menu p").removeClass("selected-status")
    $(`#status-all`).addClass("selected-status")
}

function next() {

    if (CURRENTPAGEVALUE <= 0) {
        if (CURRENTPAGEMODE == "all") {

            allJobs(CURRENTSTATUS, 1)
            ++CURRENTPAGEVALUE

        } else if (CURRENTPAGEMODE == "search") {

            searchJob(CURRENTSEARCHKEY, CURRENTSTATUS, 1)
            ++CURRENTPAGEVALUE
        }

    } else if (CURRENTPAGEVALUE < CURRENTTOTALPAGECOUNT - 1) {
        if (CURRENTPAGEMODE == "all") {
            allJobs(CURRENTSTATUS, ++CURRENTPAGEVALUE)

        } else if (CURRENTPAGEMODE == "search") {
            searchJob(CURRENTSEARCHKEY, CURRENTSTATUS, ++CURRENTPAGEVALUE)
        }
    }
}

function previous() {

    if (CURRENTPAGEVALUE > 0) {
        if (CURRENTPAGEMODE == "all") {
            allJobs(CURRENTSTATUS, --CURRENTPAGEVALUE)
        } else if (CURRENTPAGEMODE == "search") {
            searchJob(CURRENTSEARCHKEY, CURRENTSTATUS, --CURRENTPAGEVALUE)
        }
    }
}

function showStatusMenu() {
    $("#status-menu").toggle();
    $(".toggle-menu").not(`#status-menu`).hide()
}

function changeStatus(element) {

    if (element.id == "status-pending") {
        CURRENTSTATUS = "P"
    } else if (element.id == "status-opened") {
        CURRENTSTATUS = "O"
    } else if (element.id == "status-closed") {
        CURRENTSTATUS = "C"
    } else if (element.id == "status-declined") {
        CURRENTSTATUS = "D"
    } else if (element.id == "status-all") {
        CURRENTSTATUS = "Any"
    }

    $("#status-menu p").removeClass("selected-status")
    $(`#${element.id}`).addClass("selected-status")

    if (CURRENTPAGEMODE == "all") {
        allJobs(CURRENTSTATUS, CURRENTPAGEVALUE);
    } else if (CURRENTPAGEMODE == "search") {
        searchJob(CURRENTSEARCHKEY, CURRENTSTATUS, CURRENTPAGEVALUE);
    }

}

// appendJobs is a function that appends set of jobs from the result container
function appendJobs(tableBody, jobs, key) {

    jobs.forEach(
        (job) => {

            let reducedTitle = reduceTo(job.Title, 20, true)
            let reducedKey = reduceTo(key != undefined ? key : "", 20)

            let id = highlightElementEx($("<span></span>").text(job.ID), key)
            let title = highlightElementRx($("<span></span>").text(reducedTitle), reducedKey)
            let viewIcon = $("<span></span>").addClass("material-icons").text("visibility")
            let employer = job.Employer ? highlightElementEx($("<span></span>").
                text(job.Employer), key).addClass("view-employer") :
                $("<span></span>").text("*UNAVAILABLE*")

            let type = highlightElementEx(job.Type ? $("<span></span>").
                text(reConstructor(job.Type).length > 25 ?
                    reConstructor(job.Type).substring(0, 25).concat('...') : reConstructor(job.Type)) :
                $("<span></span>").text("*UNSET*").css("background-color", "#ff63477d"), key)

            let sector = $("<span></span>").css("text-transform", "capitalize").
                text(reConstructor(job.Sector).length > 25 ?
                    reConstructor(job.Sector).substring(0, 25).concat('...') : reConstructor(job.Sector))
            let status = $("<span></span>").css("text-transform", "capitalize").
                text(job.Status).addClass("status-cell").
                attr("id", "status-" + job.ID)

            let initiatorViewer = $("<span></span>").addClass("material-icons").text("account_circle")

            let row = $("<tr> </tr>")

            let col1 = $("<td> </td>").append(id)
            let col2 = $("<td> </td>").append(title)
            let col3 = $("<td> </td>").append(viewIcon)
            let col4 = $("<td> </td>").append(employer)
            let col5 = $("<td> </td>").append(type).
                attr({ "data-toggle": "tooltip", "data-placement": "top", "title": reConstructor(job.Type) })
            let col6 = $("<td> </td>").append(sector).
                attr({ "data-toggle": "tooltip", "data-placement": "top", "title": reConstructor(job.Sector)})
            let col7 = $("<td> </td>").append(setStatus(status))
            let col8 = $("<td> </td>").append(initiatorViewer).
                attr({ "data-toggle": "tooltip", "data-placement": "top", "title": "view initiator" })
            let col9 = createMoreMenu(job, $("<td> </td>"))

            // Setting onclick functions
            viewIcon.attr({ "data-toggle": "modal", "data-target": "#job-detail" })
            viewIcon.click(() =>
                viewJob(job.ID)
            )

            employer.attr({ "data-toggle": "modal", "data-target": "#user-profile" })
            employer.click(() =>
                viewUser(job.Employer, "job-employer")
            )

            initiatorViewer.attr({ "data-toggle": "modal", "data-target": "#user-profile" })
            initiatorViewer.click(() =>
                viewUser(job.InitiatorID, "job-initiator")
            )

            row.append(col1, col2, col3, col4, col5, col6, col7, col8, col9)
            tableBody.append(row)

        }
    )
}

// viewAll is a function that is fired when view all icon is shown
function viewAll() {

    CURRENTPAGEVALUE = 0

    allJobs(CURRENTSTATUS, CURRENTPAGEVALUE)
}

function allJobs(status, page) {

    // Changing the current page mode
    CURRENTPAGEMODE = "all"
    let csrf = $("#csrf").val()

    // Loading table
    loadingTable()

    $.ajax({
        url: `../../staff/jobs/all/${csrf}`,
        data: { status, page },
        type: "GET",
        success: successFunc,
        error: errorFunc,
    });

    function successFunc(result, status, xhr) {

        // Unloading table
        unLoadingTable()

        let resultContainer = JSON.parse(result);
        let currentPage = resultContainer.CurrentPage;
        let pageCount = resultContainer.PageCount;
        let jobs = resultContainer.Result;

        if (pageCount <= 1) {
            $("#pagination-div").hide();
        } else {
            $("#pagination-div").show();

            if (currentPage > 0 && currentPage < pageCount) {
                $("#prev-icon").show();
            } else {
                $("#prev-icon").hide();
            }

            if (currentPage >= 0 && currentPage < pageCount - 1) {
                $("#next-icon").show();
            } else {
                $("#next-icon").hide();
            }
        }

        CURRENTTOTALPAGECOUNT = resultContainer.PageCount
        let tableBody = $("#search-table-body").empty();
        $("#search-error-div").hide();

        if (jobs.length == 0) {
            $("#search-noresult-div").show();
            $("#search-noresult-div-msg").text(`Sorry we couldn't find any job.`);
            return;
        } else {
            $("#search-noresult-div").hide();
        }

        appendJobs(tableBody, jobs);
    }

    function errorFunc(xhr, status, error) {

        // Unloading table
        unLoadingTable()

        $("#search-table-body").empty();
        $("#pagination-div").hide();
        $("#search-noresult-div").hide();
        $("#search-error-div").show();
    }

}

// search is a function that is fired when the search button is clicked
function search(btn) {

    CURRENTSEARCHKEY = $("#search-bar").val().trim()
    CURRENTPAGEVALUE = 0

    loading(btn)

    searchJob(CURRENTSEARCHKEY, CURRENTSTATUS, CURRENTPAGEVALUE)
}

function searchJob(key, status, page) {

    // Changing the current page mode
    CURRENTPAGEMODE = "search"
    let csrf = $("#csrf").val()
    let searchButton = document.getElementById("search-button")

    // Loading table
    loadingTable()

    if (CURRENTSEARCHKEY.match(/^\s*$/)) {
        unLoading(searchButton, "Search");
        unLoadingTable();
        return
    }

    $.ajax({
        url: `../../staff/jobs/search/${csrf}`,
        data: { key, status, page },
        type: "GET",
        success: successFunc,
        error: errorFunc,
    });

    function successFunc(result, status, xhr) {

        let resultContainer = JSON.parse(result);
        let currentPage = resultContainer.CurrentPage;
        let pageCount = resultContainer.PageCount;
        let jobs = resultContainer.Result;

        // Unloading table
        unLoadingTable()
        unLoading(searchButton, "Search");

        if (pageCount <= 1) {
            $("#pagination-div").hide();
        } else {
            $("#pagination-div").show();

            if (currentPage > 0 && currentPage < pageCount) {
                $("#prev-icon").show();
            } else {
                $("#prev-icon").hide();
            }

            if (currentPage >= 0 && currentPage < pageCount - 1) {
                $("#next-icon").show();
            } else {
                $("#next-icon").hide();
            }
        }

        CURRENTTOTALPAGECOUNT = resultContainer.PageCount
        let tableBody = $("#search-table-body").empty();
        $("#search-error-div").hide();

        if (jobs.length == 0) {
            tableBody.empty();
            $("#search-noresult-div").show();
            let reducedKey = reduceTo(key != undefined ? key : "", 20, true)
            let keySpan = $("<span></span>").css("color", "#006391").text(reducedKey)
            $("#search-noresult-div-msg").text(`Sorry we couldn't find any matches for `)
            $("#search-noresult-div-msg").append(keySpan);
            return;
        } else {
            $("#search-noresult-div").hide();
        }

        appendJobs(tableBody, jobs, key);
    }

    function errorFunc(xhr, status, error) {

        // Unloading table
        unLoadingTable()
        unLoading(searchButton, "Search");

        $("#search-table-body").empty();
        $("#pagination-div").hide();
        $("#search-noresult-div").hide();
        $("#search-error-div").show();

    }

}

function hideAllExcept(section) {

    $("#search-section").hide()
    $("#add-section").hide()
    $("#profile-section").hide()

    section.show()
}

// addInternalOptions is a function that creates and add the needed internal option
function addInternalOptions() {

    let div1 = $("<div> </div>").addClass("tooltip-div")
    let div1Span1 = $("<span> </span>").addClass("material-icons").text("search").click(() => initSearchJob())
    let div1Span2 = $("<span> </span>").addClass("tooltiptext").text("search job").
        css({ "left": "50%", "width": "110px", "top": "25px", "margin-left": "-40px" })

    let div2 = $("<div> </div>").addClass("tooltip-div")
    let div2Span1 = $("<span> </span>").addClass("material-icons").text("addchart").click(() => initAddJob())
    let div2Span2 = $("<span> </span>").addClass("tooltiptext").text("add job").
        css({ "left": "50%", "width": "100px", "top": "25px", "margin-left": "-50px" })

    let div3 = $(`<div> </div>`).addClass("tooltip-div")
    let div3Span1 = $(`<span onclick="selectFile()"> </span>`).addClass("material-icons").text("reorder")
    let div3Span2 = $("<span> </span>").addClass("tooltiptext").text("bulk import").
        css({ "left": "50%", "width": "100px", "top": "25px", "margin-left": "-50px" }).
        append($(`<input id="bulk-job-uploader" class="job-uploader" type="file" 
        name="jobs" accept=".xls,.xlsx" hidden />`))

    div1.append(div1Span1, div1Span2)
    div2.append(div2Span1, div2Span2)
    div3.append(div3Span1, div3Span2)

    $("#internal-option-div").append(div1, div2, div3)

}

function showMoreMenu(id, e) {
    $(`#more-menu-${id}`).toggle().css("top", e.clientY + 10)
    $(".toggle-menu").not(`#more-menu-${id}`).hide()
}

function createMoreMenu(job, col) {

    let moreIcon = $("<span></span>").addClass("material-icons").text("more_vert").
        attr("id", "more-menu-icon-" + job.ID)
    let moreMenu = $("<div> </div>").addClass("card toggle-menu more-menu").
        attr("id", "more-menu-" + job.ID)

    if (job.Status == "P") {

        let approveJob = $("<p> </p>").append($("<i></i>").addClass("material-icons").
            text("check"), $("<span></span>").text("Approve"))

        let declineJob = $("<p> </p>").append($("<i> </i>").text("block").
            addClass("material-icons"), $("<span></span>").text("Decline"))

        approveJob.attr({ "data-toggle": "modal", "data-target": "#change-status-modal" })
        approveJob.click(() =>
            initChangeJobStatus('approve', 'search', job.ID)
        )

        declineJob.attr({ "data-toggle": "modal", "data-target": "#change-status-modal" })
        declineJob.click(() =>
            initChangeJobStatus('decline', 'search', job.ID)
        )

        moreMenu.append(approveJob, declineJob)
        col.append(moreIcon, moreMenu).
            addClass("toggle-menu-parent").click(
                (e) => {
                    showMoreMenu(job.ID, e)
                }
            )

    }

    if (job.Status == "O") {

        let re_pushJob = $("<p> </p>").append($("<i></i>").addClass("material-icons").
            text("replay"), $("<span></span>").text("Re-push"))

        let closeJob = $("<p> </p>").append($("<i></i>").addClass("material-icons").
            text("block"), $("<span></span>").text("Close"))

        re_pushJob.click(() =>
            rePushJob(job.ID)
        )

        closeJob.attr({ "data-toggle": "modal", "data-target": "#change-status-modal" })
        closeJob.click(() =>
            initChangeJobStatus('close', 'search', job.ID)
        )

        moreMenu.append(re_pushJob, closeJob)
        col.append(moreIcon, moreMenu).
            addClass("toggle-menu-parent").click(
                (e) => {
                    showMoreMenu(job.ID, e)
                }
            )
    }

    return col

}

function initChangeJobStatus(status, from, jobID) {

    // Resetting the change status modal
    $("#change-status-container").show();
    $("#change-status-modal .modal-footer").show();
    $("#change-status-error").hide();

    if (from == 'view') {
        jobID = $("#view-job-id").val()
    }

    $("#target-job-id").val(jobID)
    $("#change-status-from").val(from)

    if (status == "approve") {
        $("#target-job-status").css({ "color": "#1ab31a" }).text(status)
    } else if (status == "decline") {
        $("#target-job-status").css({ "color": "tomato" }).text(status)
    } else if (status == "close") {
        $("#target-job-status").css({ "color": "#0073aa" }).text(status)
    }

}

function changeJobStatus(btn) {

    let id = $("#target-job-id").val()
    let status = $("#target-job-status").text()
    let from = $("#change-status-from").val()
    let csrf = $("#csrf").val()

    // Adding loading spinner
    loadButtonW(btn)

    $.ajax({
        url: `../../staff/job/${id}/status/${csrf}`,
        data: { status },
        type: "PUT",
        success: successFunc,
        error: errorFunc,
    });

    function successFunc(result, status, xhr) {

        let job = JSON.parse(result)

        let statusParent = $("#status-" + id).parent()
        let moreMenuParent = $("#more-menu-" + id).parent()
        let row = moreMenuParent.parent()

        $("#status-" + id).remove()
        moreMenuParent.remove()

        let statusSpan = $("<span></span>").css("text-transform", "capitalize").
            text(job.Status).addClass("status-cell").
            attr("id", "status-" + job.ID)

        statusParent.append(setStatus(statusSpan))
        row.append(createMoreMenu(job, $("<td> </td>")))

        // Closing the 'Process Job' modal
        unLoadingButtonW(btn, "Proceed")
        $('#change-status-modal').modal('hide')

        if (from == 'view') {
            $('#job-detail').css('overflow-y', 'auto')
            setIDStatus(job.Status, $("#view-job-id").val(job.ID))
            $("#view-section .modal-footer").hide()
            $("#view-modal-spacer").show()
        }

    }

    function errorFunc(xhr, status, error) {

        // Removing spinner
        unLoadingButtonW(btn, "Proceed")

        // Instead of creating another modal we will be using the current modal
        $("#change-status-container").hide();
        $("#change-status-modal .modal-footer").hide();

        if (xhr.status == 400) {
            let errMap = JSON.parse(xhr.responseText)
            $("#change-status-error-msg").text(errorToMessage(errMap["error"]))
        } else {
            $("#change-status-error-msg").text("Unable to perform operation.")
        }
        $("#change-status-error").show();
    }
}