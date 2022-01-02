// CURRENTPAGEVALUE is a variable that can be altered by the pervious and next buttons
var CURRENTPAGEVALUE = 0;
// CURRENTTOTALPAGECOUNT is a variable that holds the current result total page
var CURRENTTOTALPAGECOUNT = 0;
// CURRENTPAGEMODE is a variable that determines whether the user is on search mode or preview all mode
var CURRENTPAGEMODE = "all";
// CURRENTROLE is a variable that holds the current role type
var CURRENTROLE = "Any";
// CURRENTSEARCHKEY is variable that holds the current search term
var CURRENTSEARCHKEY = "";

$(document).ready(onload);

function onload() {

    // Marking staff link
    $("#staff-link").addClass("active-link")

    // Adding search link to the navigation bar
    removeBaseURI()
    addBaseURI("Search")

    // Adding internal options
    addInternalOptions()

    // Showing all staff members
    hideAllExcept($("#search-section"));
    allStaff(CURRENTROLE, CURRENTPAGEVALUE);

    // Selecting default role
    $("#role-menu p").removeClass("selected-role")
    $(`#role-both`).addClass("selected-role")

    // Keybindings
    $(document).on("keyup", (e) => {
        if (e.altKey && e.which == 65) {
            $("#view-all-btn").click()
        } else if (e.altKey && e.which == 78) {
            initAddStaff()
        } else if (e.altKey && e.which == 80) {
            initUserProfile()
        } else if (e.altKey && e.which == 83) {
            initSearchStaff()
        }
    })

}

// initSearchStaff is a function that initiate or display staff search window
function initSearchStaff() {

    hideAllExcept($("#search-section"));

    removeBaseURI()
    addBaseURI("Search")

    // Resetting golbal variables
    CURRENTPAGEVALUE = 0
    CURRENTTOTALPAGECOUNT = 0
    CURRENTPAGEMODE = "all"
    CURRENTROLE = "Any"

    // Showing all staff members
    allStaff(CURRENTROLE, CURRENTPAGEVALUE);

    // Selecting default role
    $("#search-bar").val("")
    $("#role-menu p").removeClass("selected-role")
    $(`#role-both`).addClass("selected-role")
}

function next() {

    if (CURRENTPAGEVALUE <= 0) {
        if (CURRENTPAGEMODE == "all") {

            allStaff(CURRENTROLE, 1)
            ++CURRENTPAGEVALUE

        } else if (CURRENTPAGEMODE == "search") {

            searchStaff(CURRENTSEARCHKEY, CURRENTROLE, 1)
            ++CURRENTPAGEVALUE
        }

    } else if (CURRENTPAGEVALUE < CURRENTTOTALPAGECOUNT - 1) {
        if (CURRENTPAGEMODE == "all") {
            allStaff(CURRENTROLE, ++CURRENTPAGEVALUE)

        } else if (CURRENTPAGEMODE == "search") {
            searchStaff(CURRENTSEARCHKEY, CURRENTROLE, ++CURRENTPAGEVALUE)
        }
    }
}

function previous() {

    if (CURRENTPAGEVALUE > 0) {
        if (CURRENTPAGEMODE == "all") {
            allStaff(CURRENTROLE, --CURRENTPAGEVALUE)
        } else if (CURRENTPAGEMODE == "search") {
            searchStaff(CURRENTSEARCHKEY, CURRENTROLE, --CURRENTPAGEVALUE)
        }
    }
}

function showRoleMenu() {
    $("#role-menu").toggle();
    $(".toggle-menu").not("#role-menu").hide()
}

function changeRole(element) {

    if (element.id == "role-admin") {
        CURRENTROLE = "Admin"
    } else if (element.id == "role-staff") {
        CURRENTROLE = "Staff"
    } else if (element.id == "role-both") {
        CURRENTROLE = "Any"
    }

    $("#role-menu p").removeClass("selected-role")
    $(`#${element.id}`).addClass("selected-role")

    if (CURRENTPAGEMODE == "all") {
        allStaff(CURRENTROLE, CURRENTPAGEVALUE);
    } else if (CURRENTPAGEMODE == "search") {
        searchStaff(CURRENTSEARCHKEY, CURRENTROLE, CURRENTPAGEVALUE);
    }

}

// appendStaffMembers is a function that appends set of staff members from the result container
function appendStaffMembers(tableBody, staffMembers, key) {

    staffMembers.forEach(
        (staffMember) => {

            let id = highlightElementEx($("<span></sapn>").text(staffMember.ID), key)
            let firstName = highlightElementRx($("<span></sapn>").
                css("text-transform", "capitalize").text(staffMember.FirstName), key)
            let lastName = $("<span></sapn>").css("text-transform", "capitalize").text(" " + staffMember.LastName)
            let viewIcon = $("<span></sapn>").addClass("material-icons").text("visibility")
            let email = highlightElementRx($("<span></sapn>").text(staffMember.Email), key)
            let phoneNumber = highlightElementEx(
                highlightElementEx($("<span></sapn>").text(staffMember.PhoneNumber), key)
                    .text(staffMember.PhoneNumber.replace("+251", "0")), key)

            let editIcon = $("<span></sapn>").addClass("material-icons").text("create")
            let deleteIcon = $("<span></sapn>").addClass("material-icons").text("delete")

            let row = $("<tr> </tr>")

            let col1 = $("<td> </td>").append(id)
            let col2 = $("<td> </td>").append(firstName, lastName)
            let col3 = $("<td> </td>").append(viewIcon)
            let col4 = $("<td> </td>").append(email)
            let col5 = $("<td> </td>").append(phoneNumber)
            let col6 = $("<td> </td>").text(staffMember.Role).
                css({ "text-transform": "uppercase", "font-size": "12px" })
            let col7 = $("<td> </td>").append(editIcon)
            let col8 = $("<td> </td>").append(deleteIcon)

            // Setting onclick functions
            viewIcon.attr({ "data-toggle": "modal", "data-target": "#staff-member-profile" })
            viewIcon.click(() =>
                viewStaffMember(staffMember.ID)
            )

            editIcon.click(() =>
                initEditStaff(staffMember.ID)
            )

            deleteIcon.attr({ "data-toggle": "modal", "data-target": "#delete-staff-member" })
            deleteIcon.click(() =>
                initDeleteStaff(staffMember.ID)
            )

            row.append(col1, col2, col3, col4, col5, col6, col7, col8)
            tableBody.append(row)

        }
    )
}

// viewAll is a function that is fired when view all icon is shown
function viewAll() {

    CURRENTPAGEVALUE = 0

    allStaff(CURRENTROLE, CURRENTPAGEVALUE)
}

function allStaff(role, page) {

    // Changing the current page mode
    CURRENTPAGEMODE = "all"
    let csrf = $("#csrf").val()

    // Loading table
    loadingTable()

    $.ajax({
        url: `../../admin/staffs/all/${csrf}`,
        data: { role, page },
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
        let staffMembers = resultContainer.Result;

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

        if (staffMembers.length == 0) {
            $("#search-noresult-div").show();
            $("#search-noresult-div-msg").text(`Sorry we couldn't find any staff member.`);
            return;
        } else {
            $("#search-noresult-div").hide();
        }

        appendStaffMembers(tableBody, staffMembers);
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

    searchStaff(CURRENTSEARCHKEY, CURRENTROLE, CURRENTPAGEVALUE)
}

function searchStaff(key, role, page) {

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
        url: `../../admin/staffs/search/${csrf}`,
        data: { key, role, page },
        type: "GET",
        success: successFunc,
        error: errorFunc,
    });

    function successFunc(result, status, xhr) {

        let resultContainer = JSON.parse(result);
        let currentPage = resultContainer.CurrentPage;
        let pageCount = resultContainer.PageCount;
        let staffMembers = resultContainer.Result;


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

        if (staffMembers.length == 0) {
            tableBody.empty();
            $("#search-noresult-div").show();
            let keySpan = $("<span></span>").css("color", "#006391").text(key + ".")
            $("#search-noresult-div-msg").text(`Sorry we couldn't find any matches for `)
            $("#search-noresult-div-msg").append(keySpan);
            return;
        } else {
            $("#search-noresult-div").hide();
        }

        appendStaffMembers(tableBody, staffMembers, key);
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
    $("#edit-section").hide()
    $("#profile-section").hide()

    section.show()
}

// addInternalOptions is a function that creates and add the needed internal option
function addInternalOptions() {

    let div1 = $("<div> </div>").addClass("tooltip-div")
    let div1Span1 = $("<span> </span>").addClass("material-icons").text("search").click(() => initSearchStaff())
    let div1Span2 = $("<span> </span>").addClass("tooltiptext").text("search staff").
        css({ "left": "50%", "width": "80px", "top": "25px", "margin-left": "-25px" })

    let div2 = $("<div> </div>").addClass("tooltip-div")
    let div2Span1 = $("<span> </span>").addClass("material-icons").text("person_add").click(() => initAddStaff())
    let div2Span2 = $("<span> </span>").addClass("tooltiptext").text("add staff").
        css({ "left": "50%", "width": "60px", "top": "25px", "margin-left": "-25px" })

    div1.append(div1Span1, div1Span2)
    div2.append(div2Span1, div2Span2)

    $("#internal-option-div").append(div1, div2)

}