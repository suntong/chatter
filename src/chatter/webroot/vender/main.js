
var getUri = "ws://localhost:22222/p2pGet";
var addUri = "ws://localhost:22222/p2pAdd";
var addWebSocket = null;
var getWebSocket = null;
var output;
var linkToDocument;
var documentId = null;

function init() {
    output = document.getElementById("output");
    linkToDocument = document.getElementById("linkToDocument");
    document.getElementById("createPage").onclick = function(evt) {
        evt.preventDefault();
        submitPageUsingSocket(evt);
    };
    documentId = getParameterByName("id");
    createWebSocket();
}

function submitPageUsingSocket(evt) {
    var author = document.getElementById("inputAuthor").value;
    var message = document.getElementById("inputBody").value;
    var command = "SetDocument";


    var dataPacket = {};
    dataPacket["author"] = author;
    dataPacket["data"] = message;
    dataPacket["command"] = command;

    output.innerHTML = message;
    doSend(JSON.stringify(dataPacket));

}

function createWebSocket() {
    this["addWebSocket"] = new ReconnectingWebSocket(addUri, null, {debug: true, reconnectInterval: 500});
//    waitForSocketConnection(this["addWebSocket"], 3);
    this["getWebSocket"] = new ReconnectingWebSocket(getUri, null, {debug: true, reconnectInterval: 500});
//    waitForSocketConnection(this["getWebSocket"], 3);

    registerHandlers(this["addWebSocket"]);
    registerHandlers(this["getWebSocket"]);
}

function registerHandlers(webSocket) {

    webSocket.onopen = function(evt) {
        onOpen(evt)
        if(documentId == "") {
            $(".se-pre-con").fadeOut("slow");
        } else {
            getWebSocket.send(documentId);
        }
    };
    webSocket.onclose = function(evt) {
        $(".se-pre-con").fadeIn("fast");
        onClose(evt, webSocket)
    };
    webSocket.onmessage = function(evt) {
        onMessage(evt)
    };
    webSocket.onerror = function(evt) {
        onError(evt)
    };
}

function onOpen(evt) {
    log("CONNECTED");
//    $(".se-pre-con").fadeOut("fast");
}

function onClose(evt, webSocket) {
    log("DISCONNECTED");
//    $(".se-pre-con").fadeIn("fast");
}

function onMessage(evt) {
    var evtJson = JSON.parse(evt.data);
    var command = evtJson.command;
    var data = evtJson.data;

    if(command == "GetDocument") {
        log("Sending Document");
        payload = sendDocumentPayload();
        doSend(JSON.stringify(payload));
    } else if (command == "SetDocument") {
        $(".se-pre-con").fadeOut("slow");
        log("Setting Document")
        document.getElementById("submit-form").remove();
        document.getElementById("title").innerHTML = "Document...";
        output.innerHTML = data;
    } else if (command == "Join") {
        documentId = evtJson.documentId;
        log("Document Id is " + documentId);
        linkToDocument.innerHTML = '<a target="_blank" href=?id=' + documentId + ">Link to your document. Do not close this tab.</a>";
    } else {
        log("Unknown command");
        log(evt.data);
    }
}

function sendDocumentPayload() {
    var payload = {}
    payload["command"] = "SetDocument";
    payload["author"] = document.getElementById("inputAuthor").value;
    payload["data"] = document.getElementById("inputBody").value;
    payload["status"] = 200;
    if (documentId == "") {
        log("DOCUMENT ID is NULL")
    }
    payload["documentId"] = documentId;
    return payload;
}

function onError(evt) {
    log('<span style="color: red;">ERROR:</span> ' + evt.data);
}

function doSend(message) {
    addWebSocket.send(message);
}

function log(message) {
    console.log(message);
}

function getParameterByName(name) {
    name = name.replace(/[\[]/, "\\[").replace(/[\]]/, "\\]");
    var regex = new RegExp("[\\?&]" + name + "=([^&#]*)"),
    results = regex.exec(location.search);
    return results === null ? "" : decodeURIComponent(results[1].replace(/\+/g, " "));
}

window.addEventListener("load", init, false);