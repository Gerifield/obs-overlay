let socket = new WebSocket("ws://127.0.0.1:8080/websocket");

socket.onmessage = event => {
    let textField1 = document.getElementById("textField1");
    textField1.innerHTML = event.data;
    textField1.classList.add("fadeIn", "animate__animated", "animate__bounce");

    setTimeout(() => {
        textField1.classList.remove("fadeIn", "animate__animated", "animate__bounce");
    }, 3000);
};