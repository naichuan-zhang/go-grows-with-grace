const sock = new WebSocket("ws://localhost:4000/");

sock.onmessage = function (ev) {
    console.log("Received from server: ", ev.data)
}
sock.onopen = function (ev) {
    sock.send("Hello!\n")
}
