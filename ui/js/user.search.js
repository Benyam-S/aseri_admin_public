// CURRENTPAGEVALUE is a variable that can be altered by the pervious and next buttons
var CURRENTPAGEVALUE = 0;
// CURRENTTOTALPAGECOUNT is a variable that holds the current result total page
var CURRENTTOTALPAGECOUNT = 0;
// CURRENTPAGEMODE is a variable that determines whether the user is on search mode or preview all mode
var CURRENTPAGEMODE = "all";
// CURRENTCATEGORY is a variable that holds the current category type
var CURRENTCATEGORY = "Any";
// CURRENTSEARCHKEY is variable that holds the current search term
var CURRENTSEARCHKEY = "";

$(document).ready(onload);

function onload() {

    // Marking user link
    $("#user-link").addClass("active-link")

    // Adding internal options
    addInternalOptions()

    // Keybindings
    $(document).on("keyup", (e) => {
        if (e.altKey && e.which == 65) {
            $("#view-all-btn").click()
        } else if (e.altKey && e.which == 80) {
            initUserProfile()
        } else if (e.altKey && e.which == 83) {
            initSearchUser()
        }
    })

    // Adding search link to the navigation bar
    removeBaseURI()
    addBaseURI("Search")

    // Showing all users
    hideAllExcept($("#search-section"));
    allUsers(CURRENTCATEGORY, CURRENTPAGEVALUE);

    // Selecting default category
    $("#category-menu p").removeClass("selected-category")
    $(`#category-all`).addClass("selected-category")

}

// initSearchUser is a function that initiate or display user search window
function initSearchUser() {

    hideAllExcept($("#search-section"));

    removeBaseURI()
    addBaseURI("Search")

    // Resetting global variables
    CURRENTPAGEVALUE = 0
    CURRENTTOTALPAGECOUNT = 0
    CURRENTPAGEMODE = "all"
    CURRENTCATEGORY = "Any";

    // Showing all users
    allUsers(CURRENTCATEGORY, CURRENTPAGEVALUE);

    // Clearing extra search tag
    $("#search-bar").val("")
    $("#category-menu p").removeClass("selected-category")
    $(`#category-all`).addClass("selected-category")
}

function next() {

    if (CURRENTPAGEVALUE <= 0) {
        if (CURRENTPAGEMODE == "all") {

            allUsers(CURRENTCATEGORY, 1)
            ++CURRENTPAGEVALUE

        } else if (CURRENTPAGEMODE == "search") {

            searchUser(CURRENTSEARCHKEY, CURRENTCATEGORY, 1)
            ++CURRENTPAGEVALUE
        }

    } else if (CURRENTPAGEVALUE < CURRENTTOTALPAGECOUNT - 1) {
        if (CURRENTPAGEMODE == "all") {
            allUsers(CURRENTCATEGORY, ++CURRENTPAGEVALUE)

        } else if (CURRENTPAGEMODE == "search") {
            searchUser(CURRENTSEARCHKEY, CURRENTCATEGORY, ++CURRENTPAGEVALUE)
        }
    }
}

function previous() {

    if (CURRENTPAGEVALUE > 0) {
        if (CURRENTPAGEMODE == "all") {
            allUsers(CURRENTCATEGORY, --CURRENTPAGEVALUE)
        } else if (CURRENTPAGEMODE == "search") {
            searchUser(CURRENTSEARCHKEY, CURRENTCATEGORY, --CURRENTPAGEVALUE)
        }
    }
}

function showCategoryMenu() {
    $("#category-menu").toggle();
    $(".toggle-menu").not("#category-menu").hide()
}

function changeCategory(element) {

    if (element.id == "category-aseri") {
        CURRENTCATEGORY = "Aseri"
    } else if (element.id == "category-agent") {
        CURRENTCATEGORY = "Agent"
    } else if (element.id == "category-jobseeker") {
        CURRENTCATEGORY = "JobSeeker"
    } else if (element.id == "category-all") {
        CURRENTCATEGORY = "Any"
    }

    $("#category-menu p").removeClass("selected-category")
    $(`#${element.id}`).addClass("selected-category")

    if (CURRENTPAGEMODE == "all") {
        allUsers(CURRENTCATEGORY, CURRENTPAGEVALUE);
    } else if (CURRENTPAGEMODE == "search") {
        searchUser(CURRENTSEARCHKEY, CURRENTCATEGORY, CURRENTPAGEVALUE);
    }

}

// appendUsers is a function that appends set of users from the result container
function appendUsers(tableBody, users, key) {

    users.forEach(
        (user) => {

            let id = highlightElementEx($("<span></span>").text(user.ID), key)
            let userName = highlightElementRx($("<span></span>").
                css("text-transform", "capitalize").text(user.UserName), key)
            let viewIcon = $("<span></span>").addClass("material-icons").text("visibility")
            let phoneNumber = highlightElementEx(
                highlightElementEx($("<span></span>").text(user.PhoneNumber), key)
                    .text(user.PhoneNumber.replace("+251", "0")), key)

            let category = $("<span></span>").text(user.Category).css(
                { "font-size": "15px", "text-transform": "capitalize" })

            let deleteIcon = $("<span></span>").addClass("material-icons").text("delete")
            let jobsIcon = $("<span></span>").addClass("material-icons").text("work")
            let moreIcon = $("<span></span>").addClass("material-icons").text("more_vert")


            let row = $("<tr> </tr>")

            let col1 = $("<td> </td>").append(id)
            let col2 = $("<td> </td>").append(userName)
            let col3 = $("<td> </td>").append(viewIcon)
            let col4 = $("<td> </td>").append(phoneNumber)
            let col5 = $("<td> </td>").append(setCategory(category))
            let col6 = $("<td> </td>").append(deleteIcon)
            let col7 = $("<td> </td>").append(jobsIcon).
                attr({ "data-toggle": "tooltip", "data-placement": "top", "title": "view job" })
            let col8 = $("<td> </td>")

            // Setting onclick functions
            viewIcon.attr({ "data-toggle": "modal", "data-target": "#user-profile" })
            viewIcon.click(() =>
                viewUser(user.ID)
            )

            deleteIcon.attr({ "data-toggle": "modal", "data-target": "#delete-user" })
            deleteIcon.click(() =>
                initDeleteUser(user.ID)
            )

            jobsIcon.attr({ "data-toggle": "modal", "data-target": "#user-job" })
            jobsIcon.click(() =>
                initViewUserJobs(user.ID)
            )

            if (user.Category == "Agent") {
                col8.append(moreIcon, createMoreMenu(user)).
                    addClass("toggle-menu-parent").click(
                        (e) => {
                            showMoreMenu(user.ID, e)
                        }
                    )
            }

            row.append(col1, col2, col3, col4, col5, col6, col7, col8)

            tableBody.append(row)

        }
    )
}

// viewAll is a function that is fired when view all icon is shown
function viewAll() {

    CURRENTPAGEVALUE = 0

    allUsers(CURRENTCATEGORY, CURRENTPAGEVALUE)
}

function allUsers(category, page) {

    // Changing the current page mode
    CURRENTPAGEMODE = "all"
    let csrf = $("#csrf").val()

    // Loading table
    loadingTable()

    $.ajax({
        url: `../../staff/users/all/${csrf}`,
        data: { page, category },
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
        let users = resultContainer.Result;

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

        if (users.length == 0) {
            $("#search-noresult-div").show();
            $("#search-noresult-div-msg").text(`Sorry we couldn't find any user.`);
            return;
        } else {
            $("#search-noresult-div").hide();
        }

        appendUsers(tableBody, users);
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

    searchUser(CURRENTSEARCHKEY, CURRENTPAGEVALUE)
}

function searchUser(key, category, page) {

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
        url: `../../staff/users/search/${csrf}`,
        data: { key, page, category },
        type: "GET",
        success: successFunc,
        error: errorFunc,
    });

    function successFunc(result, status, xhr) {

        let resultContainer = JSON.parse(result);
        let currentPage = resultContainer.CurrentPage;
        let pageCount = resultContainer.PageCount;
        let users = resultContainer.Result;


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

        if (users.length == 0) {
            tableBody.empty();
            $("#search-noresult-div").show();
            let keySpan = $("<span></span>").css("color", "#006391").text(key + ".")
            $("#search-noresult-div-msg").text(`Sorry we couldn't find any matches for `)
            $("#search-noresult-div-msg").append(keySpan);
            return;
        } else {
            $("#search-noresult-div").hide();
        }

        appendUsers(tableBody, users, key);
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

function showMoreMenu(id, e) {

    if (e.clientY != undefined && e.clientY != null) {
        $(`#more-menu-${id}`).toggle().css("top", e.clientY + 10)
        $(".toggle-menu").not(`#more-menu-${id}`).hide()
    }
}

function createMoreMenu(user) {

    let moreMenu = $("<div> </div>").addClass("card toggle-menu more-menu").
        attr("id", "more-menu-" + user.ID).
        append($(`<input id="job-uploader-${user.ID}" class="job-uploader" type="file" name="jobs" accept=".xls,.xlsx" hidden />`))


    let addJob = $("<p> </p>").append($("<i> </i>").
        addClass("material-icons").text("post_add"), $("<span></span>").text("Add Job")).
        click(() => { navigate('../../staff/job', `add-section@${user.ID}`, "_blank") })

    let importJobs = $("<p> </p>").append($("<i> </i>").
        addClass("material-icons").text("reorder"), $("<span></span>").text("Import Jobs")).
        click(() => { selectFile(user.ID) })

    moreMenu.append(addJob)
    moreMenu.append(importJobs)

    return moreMenu

}

function hideAllExcept(section) {

    $("#search-section").hide()
    $("#profile-section").hide()

    section.show()
}

// addInternalOptions is a function that creates and add the needed internal option
function addInternalOptions() {

    let div1 = $("<div> </div>").addClass("tooltip-div")
    let div1Span1 = $("<span> </span>").addClass("material-icons").text("search").click(() => initSearchUser())
    let div1Span2 = $("<span> </span>").addClass("tooltiptext").text("search users").
        css({ "left": "50%", "width": "100px", "top": "25px", "margin-left": "-25px" })

    div1.append(div1Span1, div1Span2)

    $("#internal-option-div").append(div1)

}