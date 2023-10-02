$(function () {
    $("#postButton").on("click", function() {sendNewResponse()})
    getResponses()

})

function upvoteButton(context) {
    let responseID = parseInt(context[0].id.split("Button")[1])
    $.ajax({
        type: "POST",
        url: window.location.origin + "/response/upvote",
        data: JSON.stringify({
            "responseID": responseID
        })
    })
}

function reportButton(context) {
    let responseID = parseInt(context[0].id.split("reportButton")[1])
    //TODO: Report modal warning yadda yadda
    $.ajax({
        type: "POST",
        url: window.location.origin + "/response/report",
        data: JSON.stringify({
            "responseID": responseID
        })
    })
}

function deleteButton(context) {
    let responseID = parseInt(context[0].id.split("deleteButton")[1])
    //TODO: Deletion modal warning yadda yadda
    $.ajax({
        type: "POST",
        url: window.location.origin + "/response/delete",
        data: JSON.stringify({
            "responseID": responseID
        }),
        success: function() {
            context[0].parent().parent().parent().parent().delete()
        }
    })
}

function sendNewResponse()  {
    let response = $("#responseValue").val()
    let username = $("#responseUser").val()

    let data = {
        "username": username,
        "response": response
    }

    $('#responseValue').val("")

    $.ajax({
        type: "POST",
        url: window.location.origin + "/response/add",
        data: JSON.stringify(data),
        success: function(res) {
            fillResponses(res)
        }
    })
}

function getResponses() {
    $.ajax({
        type: "GET",
        url: window.location.origin + "/response/get",
        success: function (res) {
            fillResponses(res)
        }
    })
}

function fillResponses(data) {

    const input = JSON.parse(data)
    let responseArray = input
    if (typeof input[1] === 'boolean') {
        responseArray = input[0]
    }
    console.log(input)
    for (const element of responseArray) {
        let responseID = element.responseID,
            response = element.response,
            username = element.username,
            author = element.author,
            currentUser = input[2] ?? "",
            isEvil = input[1] ?? false

        let responseCard = `<div id="response${responseID}" class="card border-secondary mb-3 d-inline-block" style="max-width: 20rem;">
            <div class="card-body text-primary">
                <h4 class="card-title">${response}</h4>
            </div>
            <div class="card-header container">
            <div class="row">
                <div class="col h-100 align-self-center" style="min-width: fit-content">
                    <a href="https://profiles.csh.rit.edu/user/{{.username}}">${author} (${username})</a>
                </div>
                <div class="col text-right w-100 pr-1">`
        if (username !== currentUser) {
            responseCard = responseCard.concat(`
                    <button id="upvoteButton${responseID}" type="button" class="btn btn-info align-self-center p-1 mx-2 upvote"
                        style="line-height: 8px"><span class="material-symbols-outlined align-self-center">exposure_plus_1</span>
                    </button>
                    <button id="reportButton${responseID}" type="button" class="btn btn-danger align-self-center p-1 mx-1 report"
                        style="line-height: 8px"><span class="material-symbols-outlined align-self-center">report</span>
                    </button>`)
        }
        if (username === currentUser || isEvil) {
            responseCard = responseCard.concat(`
                    <button id="deleteButton${responseID}" type="button" class="btn btn-danger align-self-center p-1 mx-1 delete"
                        style="line-height: 8px"><span class="material-symbols-outlined align-self-center">delete</span>
                    </button>`)

        }
        responseCard = responseCard.concat(`
                    </div>
                <div class="w-100"></div>
                </div>
            </div>
        </div>`)
        $(responseCard).appendTo("#responseHolder")
    }
    $(".upvote").on("click", function() {upvoteButton($(this))})
    $(".report").on("click", function() {reportButton($(this))})
    $(".delete").on("click", function() {deleteButton($(this))})
}
