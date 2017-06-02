package main

import (
  "fmt"
  "net/http"
  "os/exec"
  "github.com/rs/cors"
  "io/ioutil"
  "os"
  "time"
  "strconv"
  "flag"
  "mime/multipart"
  "path"
)

// Global configurable vars, set by command line flags in initGlobals()
var dir string
var hostname string
var port string

func main() {
  initGlobals()
  startServer()
}

// Helper func to reduce error handling verbosity
func handleErr(err error) {
  if err != nil {
    fmt.Println(err)
  }
}

// Initializes the global config vars
func initGlobals() {
  dirPtr := flag.String(
    "img-dir",
    "/Users/user/Desktop",
    "directory to save and serve images from")
  hostnamePtr := flag.String(
    "hostname",
    "localhost",
    "the host being served on")
  portPtr := flag.String(
    "port",
    "8080",
    "port to listen on")
  flag.Parse()
  dir = *dirPtr
  hostname = *hostnamePtr
  port = *portPtr
}

// Given a raw image, will return the image type as a string.
// Recognized types are png, jpg, gif, and bmp.
// Returns empty string for unrecognized image types.
func getFormat(fileBytes []byte) (string) {
  bytes := fileBytes[0:4]
  if bytes[0] == 0x89 && bytes[1] == 0x50 && bytes[2] == 0x4E && bytes[3] == 0x47 { return "png" }
  if bytes[0] == 0xFF && bytes[1] == 0xD8 { return "jpg" }
  if bytes[0] == 0x47 && bytes[1] == 0x49 && bytes[2] == 0x46 && bytes[3] == 0x38 { return "gif" }
  if bytes[0] == 0x42 && bytes[1] == 0x4D { return "bmp" }
  return ""
}

func writeImageFile(imgFile multipart.File) string {
  os.Mkdir(dir + "/images", 0777)
  iconBytes, err := ioutil.ReadAll(imgFile)
  handleErr(err)
  imgFile.Close()
  imgType := getFormat(iconBytes)
  imgName := strconv.FormatInt(time.Now().Unix(), 10) + "." + imgType
  imgPath := dir + "/images/" + imgName
  err = ioutil.WriteFile(imgPath, iconBytes, 0777)
  handleErr(err)
  return imgPath
}

// Opens an image editor to mark up the image in the request.
// Returns URL to access the image after editing is completed.
func postHandler(w http.ResponseWriter, r *http.Request) {
  iconFile, _, err := r.FormFile("image")
  handleErr(err)
  imgPath := writeImageFile(iconFile)
  fmt.Println("Wrote: "+imgPath)

  cmd := exec.Command("open", "-W", imgPath)
  cmd.Run()
  fmt.Fprintf(w, "http://"+hostname+":"+port+"/images/" + path.Base(imgPath))
}

// Serves the marked up images
func getHandler(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  resource := dir + r.URL.String()
  fmt.Println("Serving: "+resource)

  file, err := os.Open(resource)
  if err != nil {
      fmt.Println(err)
      http.Error(w, err.Error(), http.StatusNotFound)
      return
  }

  modeTime := time.Now()
  http.ServeContent(w, r, resource, modeTime, file)
}

func startServer() {
  mux := http.NewServeMux()
  mux.HandleFunc("/find/", postHandler)
  mux.HandleFunc("/images/", getHandler)
  http.ListenAndServe(":"+port, cors.Default().Handler(mux))
}
