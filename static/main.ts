let socket = new WebSocket("ws://127.0.0.1:8080/websocket");
let queue = [];
let eventActive = false;

function queueProcessor() {
    if(!eventActive && queue.length > 0) {
        // We have an event, trigger it
        handleEvent(queue.shift());
    }
    setTimeout(queueProcessor, 1000);
}
queueProcessor();

function handleEvent(event) {
    eventActive = true;
    let textField1 = document.getElementById("textField1");
    textField1.innerHTML = event.data;
    textField1.classList.add("fadeIn", "animate__animated", "animate__bounce");

    setTimeout(() => {
        textField1.classList.remove("fadeIn", "animate__animated", "animate__bounce");
        // Wait until the animation end
        //textField1.addEventListener('animationend', () => {
            setTimeout(() => {eventActive = false;}, 1000 );

        //});

    }, 3000);
}

socket.onmessage = event => {
    queue.push(event);
};