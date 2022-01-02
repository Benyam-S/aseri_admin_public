
let mediaQuery1 = window.matchMedia("(max-width: 950px)")
let mediaQuery2 = window.matchMedia("(min-width: 950px)")

mediaQuery1.addListener(widthMax950)
mediaQuery2.addListener(widthMin950)

$(document).ready(onload);

function onload() {

    // First check profile pic of the user
    let userProfilePic = $("#hidden-profile-pic").val();
    if (userProfilePic) {
        $("#profile-pic").css("background-image", `url(../../assets/profilepics/${userProfilePic})`)
    }

    // Checking saff member role
    let role = $("#hidden-role").val()
    if (role == "Admin") {
        $("#staff-link").show();
    } else {
        $("#staff-link").hide();
    }

    // Setting uri
    addBaseURI()

    // Reset page state
    let pageState = $("#page-state").val()
    switch (pageState) {
        case '2':
            collapse()
            break
        case '1':
            hideNav()
            break
    }

}

function widthMax950(x) {
    if (x.matches) { // If media query matches
        closeNav()
    }
}

function widthMin950(x) {
    if (x.matches) { // If media query matches
        $("#top-logo-div img").toggle(
            $("#navbar-left").attr("class").includes("invisible-flag")
            || $("#navbar-left").attr("class").includes("collapsed-flag"));
    }
}

function collapse() {

    $("#navbar-left").toggleClass("collapsed-nav");
    $("#navbar-left").toggleClass("collapsed-flag");
    $("#top-logo-div img").toggle(
        $("#navbar-left").attr("class").includes("invisible-flag")
        || $("#navbar-left").attr("class").includes("collapsed-flag"));
    $("#left-logo-div img").toggle(100);
    $("#navbar-left-title div").toggle(100);
    $("#navbar-items-conatiner").toggleClass("collapsed-items-container");
    $("#collapsable-div").toggleClass("collapsed-coll-div");
    $("main").toggleClass("wide-main");
    $("#navbar-top").toggleClass("wide-nav")
}

function closeNav() {

    $("#navbar-left").addClass("hidden-nav");
    $("#navbar-left").addClass("invisible-flag");
    $("main").addClass("wider-main");
    $("#navbar-top").addClass("wider-nav")
    $("#hideNav-txt").text("show menu");

}

function hideNav() {

    $("#navbar-left").toggleClass("hidden-nav");
    $("#navbar-left").toggleClass("invisible-flag");
    $("#top-logo-div img").toggle(
        $("#navbar-left").attr("class").includes("invisible-flag")
        || $("#navbar-left").attr("class").includes("collapsed-flag"));
    $("#top-logo-div div").toggle(
        $("#navbar-left").attr("class").includes("invisible-flag")
        || $("#navbar-left").attr("class").includes("collapsed-flag"));
    $("main").toggleClass("wider-main");
    $("#navbar-top").toggleClass("wider-nav")

    let hideNavTxt = $("#hideNav-txt").text();
    if (hideNavTxt == "hide menu") {
        $("#hideNav-txt").text("show menu");
    } else if (hideNavTxt == "show menu") {
        $("#hideNav-txt").text("hide menu");
    }
}

function addBaseURI(extra) {

    let uri = $("#hidden-uri").val()
    let uris = uri.split("/")
    let base = "";

    for (let x = 0; x < uris.length; x++) {

        base += "/" + uris[x];
        if (x == 0) {
            continue;
        }

        let element1 = $("<div> </div>").addClass("internal").text("/");
        let element2 = $("<a> </a>").addClass("top-nav-link").click(() => { navigate(`../..${base}`) })
            .text(uris[x]).css("text-transform", "capitalize")
        $("#uri-div").append(element1, element2)

    }

    if (extra) {
        let element1 = $("<div> </div>").addClass("internal").text("/");
        let element2 = $("<div> </div>").addClass("top-nav-link internal").text(extra)
        $("#uri-div").append(element1, element2)
    }
}

function removeBaseURI() {

    $("#uri-div").find(".top-nav-link").remove()
    $("#uri-div").find(".internal").remove()

}

function navigate(href, section, location) {

    let state = $("#navbar-left").attr("class")
    let displaySection = ""

    if (state == undefined || state == null) {
        state = 0
    } else if (state.includes("invisible-flag")) {
        state = 1
    } else if (state.includes("collapsed-flag")) {
        state = 2
    } else {
        state = 0
    }
    if (section != undefined && section != null) {
        displaySection = "&display-section=" + section
    }

    // If it needs to be opend in a new tab
    if (location == "_blank") {
        window.open(`${href}?page-state=${state}${displaySection}`, '_blank')
        return
    }

    document.location.replace(`${href}?page-state=${state}${displaySection}`)
}

function showHeaderMoreMenu(event) {

    $("#header-more-menu").toggle()
    $(".toggle-menu").not("#header-more-menu").hide()
}

function logout() {

    $.ajax({
        url: "../../staff/logout",
        type: "GET",
        success: successFunc,
    });

    function successFunc(result, status, xhr) {
        window.location.replace("../../staff/login")
    }
}