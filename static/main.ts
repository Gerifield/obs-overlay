let socket = new WebSocket("ws://127.0.0.1:8080/websocket");
socket.onmessage = event => {
    let el = document.getElementById("field1");
    el.innerText = event.data
    el.style.visibility="visible";

    setTimeout(() => {
        el.style.visibility="hidden";
    }, 3000);
};