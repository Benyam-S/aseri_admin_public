// viewFeedback is a function that views selected feedback details
function viewFeedback(id) {

    let csrf = $("#csrf").val()

    $("#feedback-detail-container input").val("");
    $("#view-feedback-comment").val("")
    $("#feedback-detail-container .tooltiptext").remove();

    $.ajax({
        url: `../../staff/feedback/${id}/${csrf}`,
        type: "GET",
        success: successFunc,
        error: errorFunc,
    });

    function successFunc(result, status, xhr) {

        let feedback = JSON.parse(result);

        $("#feedback-detail-container").show();
        $("#view-error-div").hide();

        // Tooltip incase of visibility problem
        let feedbackID = feedback.ID ? $("<span></span>").text(feedback.ID).addClass("tooltiptext") : null;
        let userID = feedback.ID ? $("<span></span>").text(feedback.UserID).addClass("tooltiptext") : null;
        let createdAt = moment(feedback.CreatedAt).year() > 10 ?
            $("<span></span>").text(moment(feedback.CreatedAt).format("dddd, MMMM Do YYYY, h:mm:ss a")).
                addClass("tooltiptext") : null;


        $("#view-feedback-id").parent().addClass("tooltip-div").append(feedbackID)
        $("#view-feedback-created").parent().addClass("tooltip-div").append(createdAt)
        $("#view-feedback-user").parent().addClass("tooltip-div").append(userID)

        $("#view-feedback-id").val(feedback.ID)
        $("#view-feedback-created").val(
            moment(feedback.CreatedAt).year() > 10 ?
                moment(feedback.CreatedAt).format("dddd, MMMM Do YYYY, h:mm:ss a") : null)
        $("#view-feedback-user").val(feedback.UserID)
        $("#view-feedback-comment").val(feedback.Comment)

        // Marking feedback as seen
        if (!feedback.Seen){
            markFeedbackAsSeen(feedback.ID)
        }

    }

    function errorFunc(xhr, status, error) {

        $("#feedback-detail-container").hide();
        $("#view-section .modal-dialog").addClass("width-600")
        $("#view-error-div").show();

    }

}

// markFeedbackAsSeen is a function that marks a given feedback as seen
function markFeedbackAsSeen(id){

    let csrf = $("#csrf").val()
    $.ajax({
        url: `../../staff/feedback/${id}/${csrf}`,
        type: "PUT",
        success: successFunc,
    });

    function successFunc(result, status, xhr) {
        $("#unseen-" + id).remove()
    }
}