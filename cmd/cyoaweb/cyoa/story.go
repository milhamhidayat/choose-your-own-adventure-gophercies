package cyoa

import (
	m "choose-your-own-adventure-gophercies/cmd/cyoaweb/models"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

type HandlerOption func(h *handler)

type HandlerOpts struct {
	*template.Template
	ParseFunc func(r *http.Request) string
}

type handler struct {
	s      m.Story
	t      *template.Template
	pathFn func(r *http.Request) string
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

// WithTemplate - Higher Order Function
// return func (type), assign h.t (handler struct) with t (template) parameter
func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

// WithPathFunc - Higher order function
// accept -> function that return string
// return -> func (type), assign h.pathFn(handler struct) with fn(function) parameter
func WithPathFunc(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = fn
	}
}

// NewHandler - return http handler
func NewHandler(s m.Story, opts ...HandlerOption) http.Handler {
	// set handler struct with value
	// - story (gopher.json file)
	// - opts ..HandlerOption -> func to set template, func to set path function
	h := handler{s, tpl, defaultPathFn}
	// loop function argument
	// call the function,
	// set template with template argument to handle struct
	// set pathFn with function argument to handle struct
	for _, opt := range opts {
		opt(&h)
	}
	// return handler
	return h
}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	// get sllce from index 1 to rest
	// "/intro" -> "intro"
	return path[1:]
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// will get path /story/{path}
	path := h.pathFn(r)

	// get story data from handler.s (Story)
	if chapter, ok := h.s[path]; ok {
		// write data to template
		err := h.t.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusBadRequest)
		}
		return
	}
	http.Error(w, "Chapter not found", http.StatusNotFound)
}

// JsonStory - read and get json data from file
func JsonStory(r io.Reader) (m.Story, error) {
	d := json.NewDecoder(r)
	var story m.Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

var defaultHandlerTmpl = `
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
			<li><a href="/{{.Chapter}}">{{.Text}}</a></li>
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
