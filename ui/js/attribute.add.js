function initAddJobAttribute(category) {

    let loadingStatus = $(`#loading-${category}-status`).val();
    if (loadingStatus == "true") {
        return
    }

    // Removing any duplicate new row
    $(`#add-${category}-row`).remove()

    let row = $("<tr> </tr>").attr("id", `add-${category}-row`)
    let col1 = $("<td> </td>")
    let col2 = $("<td> </td>")
    let col3 = $("<td> </td>")

    col1.append(`
    <div>
        <input type="hidden" id="loading-${category}-status" />
        <span id="spinner-${category}-status" 
        class="spinner-border spinner-border-sm float-right" style="display:none" role="status">
        </span>
    </div>
    `)
    col2.append(`
    <div>
        <input class="form-control" type="text" id="add-${category}-name" />
        <p class="error" id="add-${category}-error"></p>
    </div>`)

    col3.append(`
    <div class="row m-0">
        <div style="margin-right: 10px;">
            <span style="color: green;" 
                id="done-${category}-btn"
                class="material-icons" 
                onclick="addJobAttribute('${category}')">done</span>
        </div>
        <div>
            <span style="color: red;"
                id="cancel-${category}-btn"
                class="material-icons"
                onclick="cancelNewRow('${category}')">clear</span>
        </div>
    </div>`)

    row.append(col1, col2, col3)

    if (category == 'sector') {
        let sectorTableBody = $("#table-sector-body")
        sectorTableBody.append(row)
    } else if (category == 'type') {
        let typeTableBody = $("#table-type-body")
        typeTableBody.append(row)
    } else if (category == 'level') {
        let levelTableBody = $("#table-level-body")
        levelTableBody.append(row)
    }

    var n = $(`#scrollable-${category}-section`).get(0).scrollHeight
    $(`#scrollable-${category}-section`).animate({ scrollTop: n }, 200);
}

function addJobAttribute(category) {

    $(`#add-${category}-error`).text("").css("visibility", "hidden")

    let csrf = $("#csrf").val()
    let name = $(`#add-${category}-name`).val();
    let loadingStatus = $(`#loading-${category}-status`).val();

    if (loadingStatus == "true") {
        return
    }

    let emptyRegx = new RegExp(/^\s*$/)
    let emptyName = emptyRegx.test(name)
    if (emptyName) {
        $(`#add-${category}-error`).text("Job attribute name can't be empty").css("visibility", "visible")
        return
    }

    // Disabling the buttons
    $(`#loading-${category}-status`).val("true");
    $(`#spinner-${category}-status`).show();
    $(`#add-${category}-row`).addClass("disabled-button")

    let jobAttributeCategory = ""
    switch (category) {
        case "sector":
            jobAttributeCategory = "job_sector"
            break;
        case "type":
            jobAttributeCategory = "job_type"
            break;
        case "level":
            jobAttributeCategory = "education_level"
            break;
    }

    $.ajax({
        url: `../../staff/attribute/${jobAttributeCategory}/${csrf}`,
        data: { name },
        type: "POST",
        success: successFunc,
        error: errorFunc,
    });

    function successFunc(result, status, xhr) {
        // Reloading the table content
        cancelNewRow(category)
        allAttributes()
    }

    function errorFunc(xhr, status, error) {

        if (xhr.status == 400) {
            errMap = JSON.parse(xhr.responseText)
            $(`#add-${category}-error`).text(errMap["error"]).css("visibility", "visible")
        } else {
            $(`#add-${category}-error`).text("Unable to add job attribute").css("visibility", "visible")
        }

        $(`#loading-${category}-status`).val("false");
        $(`#spinner-${category}-status`).hide();
        $(`#add-${category}-row`).removeClass("disabled-button")

    }
}

function cancelNewRow(category) {
    $(`#add-${category}-error`).text("").css("visibility", "hidden")
    $(`#add-${category}-row`).remove()
}