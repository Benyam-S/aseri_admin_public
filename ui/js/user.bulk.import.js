$(document).ready(function () {
  $(document).on('change', '.job-uploader', handleSelectedFile);
});

function selectFile(employerID) {
  $(`#job-uploader-${employerID}`).click()
}

function handleSelectedFile(evt) {
  let employerID = evt.target.id.substring("job-uploader-".length)
  var files = evt.target.files;
  var xl2json = new ExcelToJSON(employerID);
  xl2json.parseExcel(files[0]);

  // Clearing the input file so the same file can be uploaded again
  evt.target.value = ""
}

var ExcelToJSON = function (employerID) {

  this.parseExcel = function (file) {
    var reader = new FileReader();

    reader.onload = function (e) {
      try {
        var data = e.target.result;
        var workbook = XLSX.read(data, {
          type: 'binary'
        });

        let jobs = []
        workbook.SheetNames.forEach(function (sheetName) {
          jobs = XLSX.utils.sheet_to_row_object_array(workbook.Sheets[sheetName]);
        })

        let jobsMap = {}

        for (let index = 0; index < jobs.length; index++) {

          let jobFormData = {
            "Title": jobs[index]["Title"], "Description": jobs[index]["Description"],
            "EmployerID": employerID, "ContactType": jobs[index]["Contact Type"],
            "Type": jobs[index]["Type"], "Sector": jobs[index]["Sector"],
            "Experience": jobs[index]["Experience"], "EducationLevel": jobs[index]["Education Level"],
            "Gender": jobs[index]["Gender"]
          }

          jobsMap[index + 1] = jobFormData
        }

        let mapLength = 0
        for (var count in jobsMap) {
          mapLength++;
        }

        // If empty file is provided
        if (mapLength <= 0) {
          $("#upload-error-div-msg").
            text(`The provided file is empty or doesn't contain any jobs!`)
          $("#upload-job .modal-dialog").removeClass("modal-lg")
          hideAllExceptInModal("upload-error-div")
          return
        }

        let formData = { "jobs": JSON.stringify(jobsMap) }
        let csrf = $("#csrf").val()

        // Loading table
        loadingTable()

        $.ajax({
          url: `../../staff/jobs/${csrf}`,
          data: formData,
          type: "POST",
          success: successFunc,
          error: errorFunc,
        });


        function successFunc(result, status, xhr) {

          // Unloading table
          unLoadingTable()

          $("#upload-success-div-msg").
            text(`You have imported ${mapLength} ${mapLength > 1 ? "jobs" : "job"} from the provided file!`)
          hideAllExceptInModal("upload-success-div")
        }

        function errorFunc(xhr, status, error) {

          // Unloading table
          unLoadingTable()

          if (xhr.status == 400 || xhr.status == 500) {
            $("#upload-job .modal-dialog").addClass("modal-lg")
            $("#upload-error-detail").empty()
            hideAllExceptInModal("upload-error-detail")

            errMap = JSON.parse(xhr.responseText)
            displayAddMultipleError(errMap, xhr.status)
          } else {
            $("#upload-error-div-msg").
              text(`Unable to import jobs from the provided file!`)
            $("#upload-job .modal-dialog").removeClass("modal-lg")
            hideAllExceptInModal("upload-error-div")
          }

        }
      }
      catch (err) {

        // Unloading table
        unLoadingTable()

        $("#upload-error-div-msg").
          text(`Unable to import jobs from the provided file. Please check file type and format!`)
        $("#upload-job .modal-dialog").removeClass("modal-lg")
        hideAllExceptInModal("upload-error-div")
      }
    };

    reader.onerror = function (ex) {

      $("#upload-error-div-msg").
        text(`Unable to import jobs from the provided file. Please check file type and format!`)
      $("#upload-job .modal-dialog").removeClass("modal-lg")
      hideAllExceptInModal("upload-error-div")
    };

    reader.readAsBinaryString(file);
  };
};

function displayAddMultipleError(jobErrMap, statusCode) {

  if (statusCode == 400 && jobErrMap["0"] != undefined &&
    jobErrMap["0"]["error"] == "unable to parse jobs") {

    $("#upload-error-div-msg").
      text(`Unable to import jobs from the provided file. Please check the format of the file!`)
    $("#upload-job .modal-dialog").removeClass("modal-lg")
    hideAllExceptInModal("upload-error-div")
    return
  }

  if (statusCode == 400) {

    hideAllExceptInModal("upload-error-detail")

    for (let index in jobErrMap) {

      $("#upload-error-detail").append(
        ` <div class="card">
          <p>
            Row <span id="detail-row-index">${index}</span>:
            <span id="detail-row-msg" style="color: #e62600">Validation error has occurred! Please input valid data in the
            respective job entries.</span>
          </p>
  
          <div 
            class="col-md-3"
            data-toggle="collapse"
            data-target="#detail-row-collapsible-${index}"
            aria-expanded="false"
            aria-controls="collapseExample"
            onclick="rotateArrow('detail-row-arrow-${index}')"
          >
            <span
              class="material-icons mr-2"
              id="detail-row-arrow-${index}"
            >keyboard_arrow_right</span>
            <span>View Details</span>
          </div>
  
          <div class="collapse" id="detail-row-collapsible-${index}">
            <div class="card card-body" id="detail-row-container-${index}"></div>
          </div>
        </div>`
      )

      let errMap = jobErrMap[index]

      for (let prop in errMap) {
        $(`#detail-row-container-${index}`).
          append(`<p>
              <span class="mr-2 font-weight-bold">${getErrorName(prop)}:</span>
              <span>${errMap[prop]}.</span>
           </p>`)
      }
    }

  } else if (statusCode == 500) {

    hideAllExceptInModal("upload-error-detail")

    for (let index in jobErrMap) {

      $("#upload-error-detail").append(
        ` <div class="card">
          <p>
            Row <span id="detail-row-index">${index}</span>:
            <span id="detail-row-msg" style="color: #e62600"> Oops! unable to add the provided job. Please try again.</span>
          </p>
        </div>`
      )
    }

    $("#upload-error-detail").append(
      `<p class="text-center mr-5 ml-5 mt-3 d-flex">
        <span class="material-icons" style="color: #006391"> warning </span>
        <span> Note! All jobs except the above one's have been added to 
        the system. Please try again removing the added once. </span>
      </p>`
    )

  }

}

function hideAllExceptInModal(elementID) {
  $("#upload-error-detail").hide()
  $("#upload-success-div").hide()
  $("#upload-error-div").hide()

  $(`#${elementID}`).show()
  $("#upload-job").modal("show")
}

function rotateArrow(arrowID) {
  $(`#${arrowID}`).toggleClass("rotate-reverse")
}

function getErrorName(name) {
  switch (name) {
    case "title":
      return "Title"
    case "description":
      return "Description"
    case "employer_id":
      return "Employer ID"
    case "contact_type":
      return "Contact Type"
    case "type":
      return "Type"
    case "experience":
      return "Experience"
    case "education_level":
      return "Education Level"
    case "sector":
      return "Sector"
    case "gender":
      return "Gender"
  }
}