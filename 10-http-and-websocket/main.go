package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"html/template"
	"log"
	"net/http"
)

const listenAddr = "localhost:4000"

var rootTemplate = template.Must(template.New("root").Parse(`
<!DOCTYPE html>
<html>
	<head>
		<title>Http & Websocket example</title>
		<meta charset="UTF-8"/>
	</head>
	<body>
		<script type="text/javascript">
			const websocket = new WebSocket("ws://{{.}}/socket");
			websocket.onmessage = onMessage;
			websocket.onclose = onClose;
            websocket.onopen = onOpen;
            
            function onMessage(msg) {
                console.log("Received: ", msg)
            }
            
            function onOpen(msg) {
                websocket.send('Hello!')
            }
            
            function onClose() {
                console.log("Websocket connection closed.")
            }
		</script>
	</body>
</html>
`))

func main() {
	http.HandleFunc("/", rootHandler)
	http.Handle("/socket", websocket.Handler(socketHandler))
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	rootTemplate.Execute(w, listenAddr)
}

func socketHandler(conn *websocket.Conn) {
	var s string
	fmt.Fscan(conn, &s)
	fmt.Println("Received: ", s)
	fmt.Fprint(conn, "How do you do?")
}
