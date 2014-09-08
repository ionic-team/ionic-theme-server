package main

import (
  "net/http"
  "os"
  "log"
  "fmt"
  "bytes"
  "errors"
  "strings"
  _"io/ioutil"
  "net/url"
  "html/template"
  _"encoding/base64"

  "github.com/gorilla/mux"
  "github.com/driftyco/go-utils"
  "github.com/driftyco/gosass"
)

func MakeVariableString(variables url.Values) string {

  var buffer bytes.Buffer
  //var decodedKey string
  //var decodedValue string

  for key, value := range variables {
    //decodedKey, _ = url.QueryUnescape(key)
    //decodedValue, _ = url.QueryUnescape(value[0])
    buffer.WriteString(key + ": " + value[0] + ";\n")
  }

  log.Println("Got String:\n", buffer.String())

  return buffer.String()
}

func RawSassBuilder(version string, variables url.Values) (string, error) {
  variableString := MakeVariableString(variables)

  str := variableString + "\n@import \"ionic\";"

  return str, nil
}

func CssBuilder(version string, variables url.Values) (string, error) {
  variableString := MakeVariableString(variables)

  str := variableString + "\n@import \"ionic\";"

  ctx := gosass.Context{
    Options: gosass.Options{
      SourceComments: false,
      OutputStyle: gosass.NESTED_STYLE,
      IncludePaths: []string{"sass/nightly"},
    },
    SourceString: str,
    OutputString: "",
    ErrorStatus: 0,
    ErrorMessage: "",
  }

  gosass.Compile(&ctx)

  if ctx.ErrorStatus != 0 {
    return "", errors.New(ctx.ErrorMessage)
  }

  return ctx.OutputString, nil
}

/**
 * Compile a new Ionic Sass file from the given URL values of the form $variable=HEX (no #)
 */
func Compile(version string, format string, variables url.Values) (string, error) {
  switch format {
  case "scss":
    return RawSassBuilder(version, variables)
  case "css":
    return CssBuilder(version, variables)
  }
  return "", nil
}

func GetFormat(filename string) string {
  parts := strings.Split(filename, ".")
  ext := parts[len(parts)-1]
  return ext
}

func SassHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")

  vars := mux.Vars(r)

  log.Println("Vars", vars)

  version, ok := vars["version"]
  if !ok {
    goutils.Send400Json(w, "No version supplied")
    return
  }
  filename, ok := vars["filename"]
  if !ok {
    goutils.Send400Json(w, "No filename supplied (ex, ionic.min.css)")
    return
  }

  format := GetFormat(filename)

  sass, err := Compile(version, format, r.URL.Query())

  if err != nil {
    fmt.Fprintf(os.Stderr, "Build error %s\n", err)
    goutils.Send500Json(w, err.Error())
    return
  }

  w.Header().Set("Content-Type", "text/css")
  w.Write([]byte(sass))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
  notFoundTemplate := template.Must(template.ParseFiles("templates/404.html"))

  notFoundTemplate.Execute(w, nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
  indexTemplate := template.Must(template.ParseFiles("templates/index.html"))

  indexTemplate.Execute(w, nil)

}

func main() {
  r := mux.NewRouter()
  //r.HandleFunc("/", HomeHandler)
  r.HandleFunc("/api/sass/{version}/{filename}", SassHandler)
  r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))//http.StripPrefix("/", http.FileServer(http.Dir("./static"))))

  r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
  http.Handle("/", r)

  //http.StripPrefix("/", http.FileServer(http.Dir("../app"))))

  port := "8081"
  if len(os.Args) > 1 {
    port = os.Args[1]
  }

  log.Println("Running on port", port)
  err := http.ListenAndServe(":" + port, nil)

  if err != nil {
    log.Fatalln("Unable to start server", err)
  }
}
