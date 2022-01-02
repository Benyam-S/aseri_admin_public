// Toggle menu function
$(document).ready(
    () => {
        $(document).on("click", (event) => {
            var $trigger = $(".toggle-menu-parent");
            if ($trigger !== event.target && !$trigger.has(event.target).length) {
                $(".toggle-menu").hide();
            }
        });

        $("*").on("scroll", (event) => {
            $(".toggle-menu").hide();
        });

        $("form").on("submit", (event) => {
            event.preventDefault()
        })
    }
)

// toFLUpperCase is a function that change the provided text to text with first letter in upper case
function toFLUpperCase(txt) {

    if (txt != undefined && txt != null && txt.length > 0) {
        txt = txt.substring(0, 1).toUpperCase() + txt.substring(1).toLowerCase();
    }

    return txt;
}

// loading is a function that adds a loading animation to a given element
function loading(element) {

    let spinner = $("<span> </span>").addClass("spinner-border spinner-border-sm").
        attr({ "role": "status", "aria-hidden": true })

    $(element).prop("disabled", true).text("").append(spinner);

}

// unLoading is a function that removes a loading animation from a given element
function unLoading(element, name) {
    $(element).find(".spinner-border").remove()
    $(element).prop("disabled", false).text(name);
}

function loadButtonW(btn) {

    let spinner = $("<span> </span>").addClass("spinner-border spinner-border-sm").
        attr({ "role": "status", "aria-hidden": true })
    $(btn).empty();
    $(btn).prop("disabled", true).append(spinner);
}

function unLoadingButtonW(btn, name, iconClass, iconName) {

    let btnIcon = $("<span> </span>").addClass(iconClass).text(iconName)
    $(btn).find(".spinner-border").remove()
    $(btn).prop("disabled", false)
    $(btn).append(name ? $("<span> </span>").text(name) : null, btnIcon)
}

function loadingTable() {
    $('#search-loading-div').show()
    $('#search-result-section').addClass("disabled-button")
    $('#pagination-div').addClass("disabled-button")
}

function unLoadingTable() {
    $('#search-loading-div').hide()
    $('#search-result-section').removeClass("disabled-button")
    $('#pagination-div').removeClass("disabled-button")
}

function passwordVisibility(element) {

    let className = $(element).attr("class")
    let password = $(element).parent().find("input")

    if (className.search("icon-visible") != -1) {
        $(element).removeClass("icon-visible");
        $(element).addClass("icon-invisible");
        password.prop("type", "text");

    } else if (className.search("icon-invisible") != -1) {
        $(element).removeClass("icon-invisible");
        $(element).addClass("icon-visible");
        password.prop("type", "password");

    }
}

function validateStaffProfile(firstName, lastName, email, phoneNumber) {

    let errMap = {}

    let matchFirstName = firstName.search(/^[a-zA-Z]\w*$/)
    let matchLastName = lastName.search(/^\w*$/)
    let matchEmail = email.search(/^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/)
    let matchPhoneNumber = phoneNumber.search(/^(\+\d{11,12})$|(0\d{9})$/)


    if (matchFirstName == -1) {
        errMap["first_name"] = "firstname should only contain alpha numerical values and have at least one character"
    }
    if (matchLastName == -1) {
        errMap["last_name"] = "lastname should only contain alpha numerical values"
    }
    if (matchEmail == -1) {
        errMap["email"] = "invalid email address used"
    }

    if (matchPhoneNumber == -1) {
        errMap["phone_number"] = "phonenumber should be +XXXXXXXXXXXX or 0XXXXXXXXX formate"
    }

    return errMap
}

function validateJobEntries(postType, title, description, employer, contactType,
    jobType, educationLevel, workExperience, jobSector, gender, contactInfo, link, dueDate) {

    let errMap = {}

    let emptyRegx = new RegExp(/^\s*$/)

    let emptyTitle = emptyRegx.test(title)
    let emptyDescription = emptyRegx.test(description)
    let emptyEmployer = emptyRegx.test(employer)
    let unSelectedContactType = emptyRegx.test(contactType)
    let unSelectedJobType = emptyRegx.test(jobType)
    let unSelectedEducationLevel = emptyRegx.test(educationLevel)
    let emptyWorkExperience = emptyRegx.test(workExperience)
    let unSelectedJobSector = emptyRegx.test(jobSector)
    let unSelectedGender = emptyRegx.test(gender)
    let validGender = (new RegExp(/^[F|M|B]$/)).test(gender)
    let emptyContactInfo = emptyRegx.test(contactInfo)

    if (emptyTitle) {
        errMap["title"] = "job title can't be empty"
    } else if (title.length > 300) {
        errMap["title"] = "job title can't exceed 300 characters"
    }

    if (emptyDescription) {
        errMap["description"] = "job description can't be empty"
    } else if (description.length > 20000) {
        errMap["description"] = "job description can't exceed 20,000 characters"
    }

    if (emptyEmployer) {
        errMap["employer"] = "employer must be specified"
    }

    if (unSelectedContactType && postType == "User") {
        errMap["contact_type"] = "contact type must be specified"
    }

    if (emptyContactInfo && postType == "Internal") {
        errMap["contact_info"] = "contact info must be specified"
    }

    if (link != "" && !isValidURL(link) && postType == "External") {
        errMap["link"] = "invalid link has bee used"
    } else if (link.length > 1000 && postType == "External") {
        errMap["link"] = "link can not exceed 1000 characters"
    }

    if (emptyWorkExperience || workExperience == undefined || workExperience == null) {
        errMap["experience"] = "work experience must be specified"
    }

    if (unSelectedJobSector || jobSector == undefined || jobSector == null) {
        errMap["sector"] = "job sector must be specified"
    }

    if (postType == "User" || postType == "Internal") {

        if (unSelectedJobType || jobType == undefined || jobType == null) {
            errMap["type"] = "job type must be specified"
        }

        if (unSelectedEducationLevel || educationLevel == undefined || educationLevel == null) {
            errMap["education_level"] = "education level must be specified"
        }

        if (unSelectedGender || gender == undefined || gender == null) {
            errMap["gender"] = "gender preference must be specified"
        } else if (!validGender) {
            errMap["gender"] = "invalid gender value used"
        }

        if (dueDate != "" && dueDate != null) {
            try {
                let date = new Date(dueDate)
                let now = new Date(Date.now())
                if (date <= now.setHours(now.getHours() + 3)) {
                    errMap["due_date"] = "due date must exceed the current time at least by 3 hours"
                }
            } catch {
                errMap["due_date"] = "invalid date time formate used"
            }
        }
    }


    return errMap
}

function verifyPassword(password, vPassword) {

    let errMap = {}

    let matchPassword = password.search(/^[a-zA-Z0-9\._\-&!?=#]{8}[a-zA-Z0-9\._\-&!?=#]*$/)

    if (password.length < 8) {
        errMap["password"] = "password should contain at least 8 characters"

    } else if (matchPassword == -1) {
        errMap["password"] = "invalid characters used in password"

    }

    if (password != vPassword) {
        errMap["vPassword"] = "password does not match"
    }

    return errMap

}

function highlightElementRx(element, key) {

    let highlighted = $("<span> </span>")
    let normal = $("<span> </span>")
    let txt = element.text()
    if (key != undefined && !key.match(/^\s*$/)) {

        let keyReg = new RegExp(`^${escapeRex(key)}`, 'i');
        if (keyReg.test(txt)) {

            let subString = txt.replace(keyReg, "")
            let matches = txt.match(keyReg)
            highlighted.text(matches[0]).css({ "background-color": "yellow" })
            normal.text(subString)

            element.text("").empty()
            element.append(highlighted, normal)
        }
    }

    return element
}

function highlightElementEx(element, key) {

    let txt = element.text()
    if ((key != undefined && !key.match(/^\s*$/)) &&
        (txt != undefined && !txt.match(/^\s*$/))) {
        txt = txt.toLowerCase()
        key = key.toLowerCase()

        if (txt == key) {
            element.css({ "background-color": "yellow" })
        }
    }

    return element
}

function reduceTo(text, len, addDot) {
    let reducedText = text

    if (text.length > len) {
        reducedText = reducedText.substring(0, len) + (addDot ? "..." : "")
    }

    return reducedText
}

function escapeRex(value) {

    let escapedValue = value

    escapedValue = escapedValue.replace(/\[/, "\\\[")
    escapedValue = escapedValue.replace(/\\/, "\\\\")
    escapedValue = escapedValue.replace(/\^/, "\\\^")
    escapedValue = escapedValue.replace(/\$/, "\\\$")
    escapedValue = escapedValue.replace(/\./, "\\\.")
    escapedValue = escapedValue.replace(/\|/, "\\\|")
    escapedValue = escapedValue.replace(/\?/, "\\\?")
    escapedValue = escapedValue.replace(/\*/, "\\\*")
    escapedValue = escapedValue.replace(/\+/, "\\\+")
    escapedValue = escapedValue.replace(/\(/, "\\\(")
    escapedValue = escapedValue.replace(/\)/, "\\\)")

    return escapedValue

}

function rotate(element) {
    $(element).toggleClass("rotate")

}

function errorToMessage(error, extra) {
    let output = error;
    if (error) {

        let firstLetter = error.substr(0, 1)
        output = firstLetter.toUpperCase() + error.substring(1, error.length)

        if (extra) {
            output += extra
        } else {
            output += "."
        }
    }

    return output
}

function editUploadPhoto() {
    $("#edit-uploader").click()
}

function editReloadPhoto() {

    cleanFrame()
    $("#edit-upload-photo").show()

    let profilePic = $("#edit-hidden-profile").val()

    if (profilePic) {
        $("#edit-remove-photo").show()
        $("#edit-profile-pic").css("background-image", `url(../../assets/profilepics/${profilePic})`)
    } else {
        $("#edit-remove-photo").hide()
        $("#edit-profile-pic").css("background-image", `url(../../ui/images/user.svg)`)
    }
}

function cleanFrame() {

    $("#edit-uploading-photo").hide()
    $("#edit-remove-photo").hide()
    $("#edit-reload-photo").hide()
    $("#edit-photo-error").text("").css("visibility", "hidden");
    $("#edit-profile-pic > div").css("visibility", "hidden")
}

function setStatus(element) {
    if (element.text() == "P") {
        element.css({
            "background-color": "orange"
        }).text("Pending")
    } else if (element.text() == "O") {
        element.css({
            "background-color": "#0073aa"
        }).text("Opened")
    } else if (element.text() == "C") {
        element.css({
            "background-color": "#1ab31a"
        }).text("Closed")
    } else if (element.text() == "D") {
        element.css({
            "background-color": "tomato"
        }).text("Declined")
    }

    return element
}

function setCategory(element) {
    if (element.text() == "Aseri") {
        element.text("አሰሪ")
    } else if (element.text() == "Agent") {
        element.text("Agent")
    } else if (element.text() == "JobSeeker") {
        element.text("Job Seeker")
    }

    return element
}

function searchKeyBind(event) {
    var key = event.key;
    if (key == "Enter") {
        $("#search-button").click()
    }
}

// onEnter is a function that is fired when enter key is pressed on the element
function onEnter(event) {
    var key = event.key;
    if (key == "Enter") {
        $(event.target).click()
    }
}

function onSelectOption(selectorID, value) {

    $(`#${selectorID} option`).each(function (index, option) {
        $(option).prop("selected", false)
    });

    $(`#${selectorID} option[value="${value}"]`).prop("selected", true)
}

function isValidURL(url) {

    var pattern = /((?:[a-z][\w-]+:(?:\/{1,3}|[a-z0-9%])|www\d{0,3}[.]|[a-z0-9.\-]+[.][a-z]{2,4}\/)(?:[^\s()<>]+|\(([^\s()<>]+|(\([^\s()<>]+\)))*\))+(?:\(([^\s()<>]+|(\([^\s()<>]+\)))*\)|[^\s`!()\[\]{};:'".,<>?«»“”‘’]))/i
    return !!pattern.test(url);
}

// reConstructor just splits and reconstruct a string
function reConstructor(txt) {
    if (txt) {
        var splited = txt.split(",")
        var reConstructed = splited.join(", ")
        return reConstructed
    }

    return txt
}