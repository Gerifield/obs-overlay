var socket = new WebSocket("ws://127.0.0.1:8080/websocket");
var queue = [];
var eventActive = false;
function queueProcessor() {
    if (!eventActive && queue.length > 0) {
        // We have an event, trigger it
        handleEvent(queue.shift());
    }
    setTimeout(queueProcessor, 1000);
}
queueProcessor();
function handleEvent(event) {
    eventActive = true;
    var textField1 = document.getElementById("textField1");
    textField1.innerHTML = event.data;
    textField1.classList.add("fadeIn", "animate__animated", "animate__bounce");
    setTimeout(function () {
        textField1.classList.remove("fadeIn", "animate__animated", "animate__bounce");
        // Wait until the animation end
        //textField1.addEventListener('animationend', () => {
        setTimeout(function () { eventActive = false; }, 1000);
        //});
    }, 3000);
}
socket.onmessage = function (event) {
    queue.push(event);
};
