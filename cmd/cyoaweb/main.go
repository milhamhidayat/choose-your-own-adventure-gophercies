package main

import (
	"bufio"
	cyoa "choose-your-own-adventure-gophercies/cmd/cyoaweb/cyoa"
	m "choose-your-own-adventure-gophercies/cmd/cyoaweb/models"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// default story template
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

type StoryCli struct {
	Story m.Story
}

func receiveAnswer(options []m.Option) int {

	optionLen := len(options)

	fmt.Println("############################")
	fmt.Println("Choose your answer : ")
	fmt.Println("############################")

	for i, o := range options {
		fmt.Printf("%v. [%v]\n", i+1, o.Chapter)
		fmt.Printf("%v\n", o.Text)
	}

	fmt.Println("Your answer : ")

	reader := bufio.NewReader(os.Stdin)

	answer, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Sorry, can't read your input")
		receiveAnswer(options)
	}

	i, err := strconv.Atoi(strings.TrimSpace(strings.ToLower(strings.Replace(answer, "\n", "", -1))))

	if err != nil {
		fmt.Println("Sorry, your answer must be a number")
		receiveAnswer(options)
	}

	newAnswer := i - 1

	if newAnswer > optionLen || newAnswer < 0 {
		fmt.Println("Sorry, your answer is not in options")
		receiveAnswer(options)
	}

	return newAnswer

}

func showStory(s *StoryCli, chapter string) {
	news := s.Story[chapter]

	fmt.Println("============================")
	fmt.Println(news.Title)
	fmt.Println("============================")

	for _, p := range news.Paragraphs {
		fmt.Printf("%v\n", p)
	}

	if chapter == "home" {
		os.Exit(0)
	}

	answer := receiveAnswer(news.Options)

	next := news.Options[answer].Chapter
	showStory(s, next)
}

func main() {
	port := flag.Int("port", 3000, "the port to start the CYOA web application on")
	isCli := flag.Bool("cli", false, "Choose CLI or not")
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

	if *isCli {
		res := &StoryCli{Story: story}
		showStory(res, "intro")
	} else {
		tpl := template.Must(template.New("").Parse(storyTmpl))

		h := cyoa.NewHandler(story, cyoa.WithTemplate(tpl), cyoa.WithPathFunc(pathFn))

		mux := http.NewServeMux()
		mux.Handle("/story/", h)
		fmt.Printf("Starting the server on port : %d\n", *port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
	}

}

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}
