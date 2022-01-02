function initDeleteAttribute(id, category) {
    $("#delete-attribute-id").text(id);
    $("#delete-attribute-category").text(category);
    $("#delete-error-div").hide();
    $("#delete-attribute .modal-footer").show();
    $("#delete-attribute-container").show();
}

function deleteAttribute(btn) {

    let csrf = $("#csrf").val()
    let id = $("#delete-attribute-id").text();
    let category = $("#delete-attribute-category").text();

    loading(btn);

    let jobAttributeCategory = ""
    switch (category) {
        case "job sectors":
            jobAttributeCategory = "job_sector"
            break;
        case "job types":
            jobAttributeCategory = "job_type"
            break;
        case "education levels":
            jobAttributeCategory = "education_level"
            break;
    }

    $.ajax({
        url: `../../staff/attribute/${jobAttributeCategory}/${id}/${csrf}`,
        type: "DELETE",
        success: successFunc,
        error: errorFunc,
    });

    function successFunc(result, status, xhr) {
        $('#delete-attribute').modal('hide');
        unLoading(btn, "Delete");

        // Reloading the table content
        allAttributes()
    }

    function errorFunc(xhr, status, error) {

        unLoading(btn, "Delete");

        $("#delete-attribute-container").hide();
        $("#delete-attribute .modal-footer").hide();
        $("#delete-error-div-msg").text("Unable to delete the selected job attribute.")
        $("#delete-error-div").show();

    }
}