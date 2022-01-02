$(document).ready(onload);

let botSubscribersAnimated = false;
let referralVisitsAnimated = false;
$(window).on("scroll", function () {
    if ($(window).scrollTop() + $(window).height() - 100 >=
        $("#bot-subscribers-count").offset().top && !botSubscribersAnimated) {

        let botSubscribersCount = $("#hidden-subscribers-count").val()

        if (botSubscribersCount != "" && botSubscribersCount != undefined && botSubscribersCount != null) {
            animateValue("bot-subscribers-count", 0, botSubscribersCount, 1000);
            botSubscribersAnimated = true;
        }
    }

    if ($(window).scrollTop() + $(window).height() - 100 >=
        $("#referral-link-visits").offset().top && !referralVisitsAnimated) {

        let referralLinkVisits = $("#hidden-referral-visits").val()

        if (referralLinkVisits != "" && referralLinkVisits != undefined && referralLinkVisits != null) {
            animateValue("referral-link-visits", 0, referralLinkVisits, 1000);
            referralVisitsAnimated = true;
        }
    }
})

function onload() {
    $("#dashboard-link").addClass("active-link")

    // Keybindings
    $(document).on("keyup", (e) => {
        if (e.altKey && e.which == 80) {
            initUserProfile()
        }
    })

    getInitData()
}

function getInitData() {

    let csrf = $("#csrf").val()

    $.ajax({
        url: `../../staff/dashboard/init/${csrf}`,
        type: "GET",
        success: successFunc,

    });

    let spinner1 = $("<span> </span>").addClass("spinner-border spinner-border-sm").
        attr({ "role": "status", "aria-hidden": true }).css({ "width": "50px", "height": "50px" })
    let spinner2 = $("<span> </span>").addClass("spinner-border spinner-border-sm").
        attr({ "role": "status", "aria-hidden": true }).css({ "width": "50px", "height": "50px" })
    $("#users-statistics").empty().append(spinner1).
        css({ "display": "flex", "align-items": "center", "justify-content": "center" });
    $("#jobs-statistics").empty().append(spinner2).
        css({ "display": "flex", "align-items": "center", "justify-content": "center" });

    function successFunc(result, status, xhr) {
        let resultContainer = JSON.parse(result);
        let jobsTotal = resultContainer["jobs_total"]
        let usersTotal = resultContainer["users_total"]

        $("#job-pending-count").text(numberWithCommas(jobsTotal.pending))
        $("#job-pending-container").
            attr({ "data-toggle": "tooltip", "data-placement": "top", "title": numberWithCommas(jobsTotal.pending) })
        $("#job-opened-count").text(numberWithCommas(jobsTotal.opened))
        $("#job-opened-container").
            attr({ "data-toggle": "tooltip", "data-placement": "top", "title": numberWithCommas(jobsTotal.opened) })
        $("#job-closed-count").text(numberWithCommas(jobsTotal.closed))
        $("#job-closed-container").
            attr({ "data-toggle": "tooltip", "data-placement": "top", "title": numberWithCommas(jobsTotal.closed) })
        $("#job-declined-count").text(numberWithCommas(jobsTotal.declined))
        $("#job-declined-container").
            attr({ "data-toggle": "tooltip", "data-placement": "top", "title": numberWithCommas(jobsTotal.declined) })

        $("#user-seeker-count").text(numberWithCommas(usersTotal["job_seeker"]))
        $("#user-seeker-container").
            attr({ "data-toggle": "tooltip", "data-placement": "top", "title": numberWithCommas(usersTotal["job_seeker"]) })
        $("#user-agent-count").text(numberWithCommas(usersTotal["agent"]))
        $("#user-agent-container").
            attr({ "data-toggle": "tooltip", "data-placement": "top", "title": numberWithCommas(usersTotal["agent"]) })
        $("#user-aseri-count").text(numberWithCommas(usersTotal["aseri"]))
        $("#user-aseri-container").
            attr({ "data-toggle": "tooltip", "data-placement": "top", "title": numberWithCommas(usersTotal["aseri"]) })

        $("#hidden-subscribers-count").val(resultContainer["subscribers_total"])
        $("#hidden-referral-visits").val(resultContainer["referral_count"])

        let users = resultContainer["all_users"]
        let jobs = resultContainer["all_jobs"]

        let userCountMap = getGraphData(users)
        userCountMap = new Map([...userCountMap.entries()].sort());
        let usersStartPoint = userCountMap.keys().next().value
        let usersData = Array.from(userCountMap.values())
        let usersDates = Array.from(userCountMap.keys())

        for (let i = 0; i < usersDates.length - 1; i++) {
            if (usersDates[i] + 86400000 < usersDates[i + 1]) {
                let tempDates = usersDates.copyWithin(0, usersDates.length)
                let tempData = usersData.copyWithin(0, usersData.length)

                usersDates = usersDates.slice(0, i + 1)
                usersDates.push(tempDates[i] + 86400000)
                usersDates = usersDates.concat(tempDates.slice(i + 1, tempDates.length))

                usersData = usersData.slice(0, i + 1)
                usersData.push(null)
                usersData = usersData.concat(tempData.slice(i + 1, tempData.length))

            }
        }

        let jobCountMap = getGraphData(jobs)
        jobCountMap = new Map([...jobCountMap.entries()].sort());
        let jobsStartPoint = jobCountMap.keys().next().value
        let jobsData = Array.from(jobCountMap.values())
        let jobsDates = Array.from(jobCountMap.keys())

        for (let i = 0; i < jobsDates.length - 1; i++) {
            if (jobsDates[i] + 86400000 < jobsDates[i + 1]) {
                let tempDates = jobsDates.copyWithin(0, jobsDates.length)
                let tempData = jobsData.copyWithin(0, jobsData.length)

                jobsDates = jobsDates.slice(0, i + 1)
                jobsDates.push(tempDates[i] + 86400000)
                jobsDates = jobsDates.concat(tempDates.slice(i + 1, tempDates.length))

                jobsData = jobsData.slice(0, i + 1)
                jobsData.push(null)
                jobsData = jobsData.concat(tempData.slice(i + 1, tempData.length))

            }
        }


        $("#users-statistics").empty()
        $("#jobs-statistics").empty()
        createGraph(usersData, usersStartPoint, 86400000, "Users Statistics", "User Count", "users-statistics")
        createGraph(jobsData, jobsStartPoint, 86400000, "Jobs Statistics", "Job Count", "jobs-statistics")
    }

}

function hideAllExcept(section) {

    $("#dashboard-section").hide()
    $("#profile-section").hide()

    section.show()
}

function createGraph(data, startPoint, interval, title, yTitle, target) {

    // Create a timer
    var start = +new Date();

    // Create the chart
    let chart = Highcharts.stockChart(target, {
        chart: {
            events: {
                load: function () {
                    if (!window.TestController) {
                        this.setTitle(null, {
                            text: 'Built chart in ' + (new Date() - start) + 'ms'
                        });
                    }
                }
            },
            zoomType: 'x'
        },

        rangeSelector: {

            buttons: [{
                type: 'day',
                count: 3,
                text: '3d'
            }, {
                type: 'week',
                count: 1,
                text: '1w'
            }, {
                type: 'month',
                count: 1,
                text: '1m'
            }, {
                type: 'month',
                count: 6,
                text: '6m'
            }, {
                type: 'year',
                count: 1,
                text: '1y'
            }, {
                type: 'all',
                text: 'All'
            }],
            selected: 3
        },

        yAxis: {
            title: {
                text: yTitle
            }
        },

        title: {
            text: title
        },

        series: [{
            name: 'Count',
            data: data,
            pointStart: startPoint,
            pointInterval: interval,
        }]

    });
}

function getGraphData(list) {

    let elementCountMap = new Map()

    if (list != null && list != undefined) {

        list.forEach(element => {
            let newDate = new Date(element.CreatedAt).toDateString()
            newDate = new Date(newDate).getTime()

            let isNew = true

            elementCountMap.forEach((count, date) => {
                if (newDate == date) {
                    elementCountMap.set(newDate, count + 1)
                    isNew = false
                    return
                }
            });

            if (isNew) {
                elementCountMap.set(newDate, 1)
            }
        });
    }

    return elementCountMap
}

function animateValue(id, start, end, duration) {

    start = parseInt(start)
    end = parseInt(end)

    if (start == NaN) {
        start = 0
    }

    if (end == NaN) {
        end = 0
    }

    if (start === end) return;
    var range = end - start;
    var current = start;
    var increment = end > start ? 1 : -1;
    var stepTime = Math.abs(Math.floor(duration / range));
    var obj = document.getElementById(id);
    var timer = setInterval(function () {
        current += increment;
        obj.innerHTML = numberWithCommas(current);
        if (current == end) {
            clearInterval(timer);
        }
    }, stepTime);
}

function numberWithCommas(x) {
    return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}