package main

import (
	cyoa "choose-your-own-adventure-gophercies/cmd/cyoaweb/cyoa"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

var storyTmpl = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
	<title>Choose your own adventure</title>
</head>

<body>
	<section class="page">
		<h1>{{.Title}}</h1>
		{{range .Paragraphs}}
		<p>{{.}}</p>
		{{end}}
		<ul>
			{{range .Options}}
			<li><a href="/story/{{.Chapter}}">{{.Text}}</a></li>
			{{end}}
		</ul>
	</section>
	<style>
      body {
        font-family: helvetica, arial;
      }
      h1 {
        text-align:center;
        position:relative;
      }
      .page {
        width: 80%;
        max-width: 500px;
        margin: auto;
        margin-top: 40px;
        margin-bottom: 40px;
        padding: 80px;
        background: #FCF6FC;
        border: 1px solid #eee;
        box-shadow: 0 10px 6px -6px #797;
      }
      ul {
        border-top: 1px dotted #ccc;
        padding: 10px 0 0 0;
        -webkit-padding-start: 0;
      }
      li {
        padding-top: 10px;
      }
      a,
      a:visited {
        text-decoration: underline;
        color: #555;
      }
      a:active,
      a:hover {
        color: #222;
      }
      p {
        text-indent: 1em;
      }
    </style>
</body>

</html>`

func main() {
	port := flag.Int("port", 3000, "the port to start the CYOA web application on")
	filename := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	flag.Parse()
	fmt.Printf("Using the story in %s. \n", *filename)

	f, err := os.Open(*filename)

	if err != nil {
		panic(err)
	}

	story, err := cyoa.JsonStory(f)
	if err != nil {
		panic(err)
	}

	// tpl := template.Must(template.New("").Parse("Hello"))
	// h := cyoa.NewHandler(story, cyoa.WithTemplate(tpl))

	tpl := template.Must(template.New("").Parse(storyTmpl))

	h := cyoa.NewHandler(story, cyoa.WithTemplate(tpl), cyoa.WithPathFunc(pathFn))

	mux := http.NewServeMux()
	mux.Handle("/story/", h)

	fmt.Printf("Starting the server on port : %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))

}

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	// get sllce from index 1 to rest
	// "/intro" -> "intro"
	return path[len("/story/"):]
}
