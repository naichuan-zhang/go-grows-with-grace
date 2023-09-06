package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"html/template"
	"io"
	"log"
	"net/http"
)

const listenAddr = "localhost:4000"

type socket struct {
	conn *websocket.Conn
	done chan bool
}

func (s socket) Read(b []byte) (int, error) {
	return s.conn.Read(b)
}

func (s socket) Write(b []byte) (int, error) {
	return s.conn.Write(b)
}

func (s socket) Close() error {
	s.done <- true
	return nil
}

var rootTemplate = template.Must(template.New("root").Parse(`
<!DOCTYPE html>
<html>
	<head>
		<title>Websocket Chat - Golang</title>
		<meta charset="UTF-8"/>
		<script>
			var input, output, websocket;

			function showMessage(msg) {
                var p = document.createElement("p");
                p.innerHTML = msg;
                output.appendChild(p);
			}
		
			function onMessage(e) {
                showMessage(e.data);
			}
            
            function onClose() {
                showMessage("Connection closed.");
            }
            
            function onKeyUp(e) {
                if (e.keyCode === 13) {
                    sendMessage();
                }
            }
            
            function sendMessage() {
                var msg = input.value;
                input.value = "";
                websocket.send(msg + "\n");
                showMessage(msg);
            }
            
            function init() {
                input = document.getElementById("input");
                input.addEventListener("keyup", onKeyUp, false);
                output = document.getElementById("output");
                websocket = new WebSocket("ws://{{.}}/socket");
                websocket.onmessage = onMessage;
                websocket.onclose = onClose;
            }
            
            window.addEventListener("load", init);
		</script>
	</head>
	<body>
		Say: <input id="input" type="text"/>
		<div id="output"></div>
	</body>
</html>
`))

func rootHandler(w http.ResponseWriter, r *http.Request) {
	rootTemplate.Execute(w, listenAddr)
}

func socketHandler(conn *websocket.Conn) {
	s := socket{conn: conn, done: make(chan bool)}
	go match(s)
	<-s.done
}

var partner = make(chan io.ReadWriteCloser)

func match(c io.ReadWriteCloser) {
	fmt.Fprintln(c, "Waiting for a partner...")
	select {
	case partner <- c:
		// handled by the other goroutine
	case p := <-partner:
		chat(p, c)
	}
}

func chat(a, b io.ReadWriteCloser) {
	fmt.Fprintln(a, "Found one! Say hi.")
	fmt.Fprintln(b, "Found one! Say hi.")
	errc := make(chan error, 1)
	go cp(a, b, errc)
	go cp(b, a, errc)
	if err := <-errc; err != nil {
		log.Println(err)
	}
	a.Close()
	b.Close()
	go io.Copy(b, a)
	io.Copy(a, b)
}

func cp(w io.Writer, r io.Reader, errc chan<- error) {
	_, err := io.Copy(w, r)
	errc <- err
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.Handle("/socket", websocket.Handler(socketHandler))
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
