// CURRENTPAGEVALUE is a variable that can be altered by the pervious and next buttons
var CURRENTPAGEVALUE = 0;
// CURRENTTOTALPAGECOUNT is a variable that holds the current result total page
var CURRENTTOTALPAGECOUNT = 0;
// CURRENTPAGEMODE is a variable that determines whether the user is on search mode or preview all mode
var CURRENTPAGEMODE = "all";
// CURRENTSTATUS is a variable that holds the current status type
var CURRENTSTATUS = "Both";
// CURRENTSEARCHKEY is variable that holds the current search term
var CURRENTSEARCHKEY = "";


$(document).ready(onload);

function onload() {

    // Marking feedback link
    $("#feedback-link").addClass("active-link")

    // Keybindings
    $(document).on("keyup", (e) => {
        if (e.altKey && e.which == 65) {
            $("#view-all-btn").click()
        } else if (e.altKey && e.which == 80) {
            initUserProfile()
        } else if (e.altKey && e.which == 83) {
            initSearchFeedback()
        }
    })

    // Adding search link to the navigation bar
    removeBaseURI()
    addBaseURI("Search")

    // Showing all feedbacks
    hideAllExcept($("#search-section"));
    allFeedbacks(CURRENTSTATUS, CURRENTPAGEVALUE);

    // Selecting default status
    $("#status-menu p").removeClass("selected-status")
    $(`#status-both`).addClass("selected-status")

}

// initSearchFeedback is a function that initiate or display feedback search window
function initSearchFeedback() {

    hideAllExcept($("#search-section"));

    removeBaseURI()
    addBaseURI("Search")

    // Resetting global variables
    CURRENTPAGEVALUE = 0
    CURRENTTOTALPAGECOUNT = 0
    CURRENTPAGEMODE = "all"
    CURRENTSTATUS = "Both"

    // Showing all feedbacks
    allFeedbacks(CURRENTSTATUS, CURRENTPAGEVALUE);

    // Selecting default status
    $("#search-bar").val("")
    $("#status-menu p").removeClass("selected-status")
    $(`#status-both`).addClass("selected-status")
}

function next() {

    if (CURRENTPAGEVALUE <= 0) {
        if (CURRENTPAGEMODE == "all") {

            allFeedbacks(CURRENTSTATUS, 1)
            ++CURRENTPAGEVALUE

        } else if (CURRENTPAGEMODE == "search") {

            searchFeedback(CURRENTSEARCHKEY, CURRENTSTATUS, 1)
            ++CURRENTPAGEVALUE
        }

    } else if (CURRENTPAGEVALUE < CURRENTTOTALPAGECOUNT - 1) {
        if (CURRENTPAGEMODE == "all") {
            allFeedbacks(CURRENTSTATUS, ++CURRENTPAGEVALUE)

        } else if (CURRENTPAGEMODE == "search") {
            searchFeedback(CURRENTSEARCHKEY, CURRENTSTATUS, ++CURRENTPAGEVALUE)
        }
    }
}

function previous() {

    if (CURRENTPAGEVALUE > 0) {
        if (CURRENTPAGEMODE == "all") {
            allFeedbacks(CURRENTSTATUS, --CURRENTPAGEVALUE)
        } else if (CURRENTPAGEMODE == "search") {
            searchFeedback(CURRENTSEARCHKEY, CURRENTSTATUS, --CURRENTPAGEVALUE)
        }
    }
}

function showStatusMenu() {
    $("#status-menu").toggle();
    $(".toggle-menu").not(`#status-menu`).hide()
}

function changeStatus(element) {

    if (element.id == "status-unseen") {
        CURRENTSTATUS = "Unseen"
    } else if (element.id == "status-seen") {
        CURRENTSTATUS = "Seen"
    } else if (element.id == "status-both") {
        CURRENTSTATUS = "Both"
    }

    $("#status-menu p").removeClass("selected-status")
    $(`#${element.id}`).addClass("selected-status")

    if (CURRENTPAGEMODE == "all") {
        allFeedbacks(CURRENTSTATUS, CURRENTPAGEVALUE);
    } else if (CURRENTPAGEMODE == "search") {
        searchFeedback(CURRENTSEARCHKEY, CURRENTSTATUS, CURRENTPAGEVALUE);
    }

}

// appendFeedbacks is a function that appends set of feedbacks from the result container
function appendFeedbacks(tableBody, feedbacks, key) {

    feedbacks.forEach(
        (feedback) => {

            let reducedComment = reduceTo(feedback.Comment, 30, true)
            let reducedKey = reduceTo(key != undefined ? key : "", 30)

            let id = highlightElementEx($("<span></span>").text(feedback.ID), key)
            let comment = highlightElementRx($("<span></span>").text(reducedComment), reducedKey)
            let viewIcon = $("<span></span>").addClass("material-icons").text("visibility")
            let userID = highlightElementEx($("<span></span>").
                text(feedback.UserID ? feedback.UserID : "*UNAVAILABLE*"), key)

            let createdAt = moment(feedback.CreatedAt).year() > 10 ?
                $("<span></span>").text(moment(feedback.CreatedAt).format("dddd, MMM Do YYYY")).
                    addClass("tooltiptext") : null;


            let row = $("<tr> </tr>")

            let col1 = $("<td> </td>").append(id)
            let col2 = $("<td> </td>").append(comment)
            let col3 = $("<td> </td>").append(viewIcon)
            let col4 = $("<td> </td>").append(userID)
            let col5 = $("<td> </td>").append(createdAt)
            let col6 = $("<td> </td>")

            // Setting onclick functions
            viewIcon.attr({ "data-toggle": "modal", "data-target": "#feedback-detail" })
            viewIcon.click(() =>
                viewFeedback(feedback.ID)
            )

            if (!feedback.Seen) {
                let newDot = $("<div> </div>").css({
                    "background-color": "orange", "width": "12px",
                    "height": "12px", "border-radius": "50%"
                }).attr("id", "unseen-" + feedback.ID)
                col6.append(newDot)
            }

            row.append(col1, col2, col3, col4, col5, col6)
            tableBody.append(row)

        }
    )
}

// viewAll is a function that is fired when view all icon is clicked
function viewAll() {

    CURRENTPAGEVALUE = 0

    allFeedbacks(CURRENTSTATUS, CURRENTPAGEVALUE)
}

function allFeedbacks(status, page) {

    // Changing the current page mode
    CURRENTPAGEMODE = "all"
    let csrf = $("#csrf").val()

    // Loading table
    loadingTable()

    $.ajax({
        url: `../../staff/feedbacks/all/${csrf}`,
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
        let feedbacks = resultContainer.Result;

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

        if (feedbacks.length == 0) {
            $("#search-noresult-div").show();
            $("#search-noresult-div-msg").text(`Sorry we couldn't find any feedback.`);
            return;
        } else {
            $("#search-noresult-div").hide();
        }

        appendFeedbacks(tableBody, feedbacks);
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

    searchFeedback(CURRENTSEARCHKEY, CURRENTSTATUS, CURRENTPAGEVALUE)
}

function searchFeedback(key, status, page) {

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
        url: `../../staff/feedbacks/search/${csrf}`,
        data: { key, status, page },
        type: "GET",
        success: successFunc,
        error: errorFunc,
    });

    function successFunc(result, status, xhr) {

        let resultContainer = JSON.parse(result);
        let currentPage = resultContainer.CurrentPage;
        let pageCount = resultContainer.PageCount;
        let feedbacks = resultContainer.Result;

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

        if (feedbacks.length == 0) {
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

        appendFeedbacks(tableBody, feedbacks, key);
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
    $("#profile-section").hide()

    section.show()
}