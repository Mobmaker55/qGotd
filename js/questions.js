$(function () {
    $("#newQuestionBtn").on("click", function() {toggleForm()})

})

function toggleForm() {
    $("#newQuestionBtn").hide("fast", "linear")
    $("#newQuestionForm").show("fast", "linear")
}

dateTime = new Date();
date = dateTime.toISOString().split("T")[0]