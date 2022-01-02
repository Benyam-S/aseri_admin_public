$(document).ready(onload);

function onload() {
    // Keybindings
    $(document).on("keyup", (e) => {
        if (e.which == 13) {
            login($("#login-button")[0])
        }
    })
}

function login(btn) {

    let identifier = $('#identifier').val();
    let password = $('#password').val();

    if (!identifier || !password) {
        return;
    }

    loading(btn);

    $.ajax({
        url: "../../staff/login",
        data: { identifier, password },
        type: "POST",
        success: successFunc,
        error: errorFunc,
    });

    function successFunc(result, status, xhr) {

        unLoading(btn, "Log In");
        $(btn).prop("disabled", true) // This is used to prevent from double login

        $("#login-error").css("visibility", "hidden");
        if (location.href.search("/staff/login") >= 0) {
            window.location.replace("../../staff/dashboard")
        } else {
            window.location.replace(location.href)
        }
    }

    function errorFunc(xhr, status, error) {

        unLoading(btn, "Log In");
        $("#login-error").css("visibility", "visible");
    }

}