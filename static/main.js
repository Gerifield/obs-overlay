var socket = new WebSocket("ws://127.0.0.1:8080/websocket");
socket.onmessage = function (event) {
    var el = document.getElementById("field1");
    el.innerText = event.data;
    el.style.visibility = "visible";
    setTimeout(function () {
        el.style.visibility = "hidden";
    }, 3000);
};
