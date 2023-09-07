package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"html/template"
	"io"
	"log"
	"markov"
	"net"
	"net/http"
	"time"
)

const listenAddr = "localhost:4000"

var chain = markov.NewChain(2) // 2-word prefixes

type socket struct {
	io.Reader
	io.Writer
	done chan bool
}

func (s socket) Close() error {
	s.done <- true
	return nil
}

// Bot runs an io.ReadWriteCloser that responds to
// each incoming write with a generated sentence.
func Bot() io.ReadWriteCloser {
	r, out := io.Pipe()
	return bot{r, out}
}

type bot struct {
	io.ReadCloser
	out io.Writer
}

func (b bot) Write(p []byte) (int, error) {
	go b.speak()
	return len(p), nil
}

func (b bot) speak() {
	time.Sleep(time.Second)
	msg := chain.Generate(10) // at most 10 words
	b.out.Write([]byte(msg))
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
	s := socket{conn, conn, make(chan bool)}
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
	// The bot should jump in if a real partner doesn't join.
	// To do this, we add a case to the select that triggers after 5 seconds,
	// starting a chat between the user's socket and a bot.
	case <-time.After(5 * time.Second):
		chat(Bot(), c)
	}
}

// The chat function remains untouched.
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

// TCP and HTTP at the same time.
func main() {
	go netListen()
	http.HandleFunc("/", rootHandler)
	http.Handle("/socket", websocket.Handler(socketHandler))
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func netListen() {
	l, err := net.Listen("tcp", "localhost:4001")
	if err != nil {
		log.Fatal(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go match(c)
	}
}
