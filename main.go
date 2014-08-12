package main

import (
  "net/http"
  "os"
  "log"
  "fmt"
  "html/template"

  "github.com/gorilla/mux"
  "github.com/driftyco/go-utils"
  "github.com/moovweb/gosass"
)

func SassHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  log.Println("Vars", vars)

  _, ok := vars["version"]
  if !ok {
    goutils.Send400Json(w, "No version supplied")
    return
  }
  _, ok = vars["format"]
  if !ok {
    goutils.Send400Json(w, "No format supplied (ex, ionic.min.css)")
    return
  }

  path := "sass/nightly/ionic.scss"

  ctx := gosass.FileContext{
    Options: gosass.Options{
      OutputStyle: gosass.NESTED_STYLE,
      IncludePaths: make([]string, 0),
    },
    InputPath: path,
    OutputString: "",
    ErrorStatus: 0,
    ErrorMessage: "",
  }

  gosass.CompileFile(&ctx)

  if ctx.ErrorStatus != 0 {
    fmt.Fprintf(os.Stderr, "Build error %s\n", ctx.ErrorMessage)
    goutils.Send500Json(w, "Build error on server")
    return
  }

  goutils.JsonResponse(w, map[string]string{
    "status": "success",
    "sass": ctx.OutputString,
  }, 200)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
  notFoundTemplate := template.Must(template.ParseFiles("templates/404.html"))

  notFoundTemplate.Execute(w, nil)
}

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/{version}/{format}", SassHandler)

  r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
  http.Handle("/", r)

  //http.Handle("/", http.StripPrefix("/app/", http.FileServer(http.Dir("../app"))))
  //http.StripPrefix("/", http.FileServer(http.Dir("../app"))))

  port := "8080"
  if len(os.Args) > 1 {
    port = os.Args[1]
  }

  log.Println("Running on port", port)
  err := http.ListenAndServe(":" + port, nil)

  if err != nil {
    log.Fatalln("Unable to start server", err)
  }
}
