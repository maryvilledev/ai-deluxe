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
var verbose bool

func main() {
  initGlobals()
  startServer()
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
  verbosePtr := flag.Bool(
    "verbose",
    false,
    "use verbose logging")
  flag.Parse()
  dir = *dirPtr
  hostname = *hostnamePtr
  port = *portPtr
  verbose = *verbosePtr
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

// Write the image to the file system under the "/images" dir within
// the path specified by the global "dir" var.
func writeImageFile(imgFile multipart.File) (string, error) {
  os.Mkdir(dir + "/images", 0777)
  iconBytes, err := ioutil.ReadAll(imgFile)
  if err != nil {
    return "", err
  }
  imgFile.Close()
  imgType := getFormat(iconBytes)
  imgName := strconv.FormatInt(time.Now().Unix(), 10) + "." + imgType
  imgPath := dir + "/images/" + imgName
  err = ioutil.WriteFile(imgPath, iconBytes, 0777)
  if err != nil {
    return "", err
  }
  if verbose {
    fmt.Println("Wrote: "+imgPath)
  }
  return imgPath, nil
}

// Opens an editor to mark up the image. Blocks until the editor
// is closed.
func editImage(imgPath string) {
  if verbose {
    fmt.Println("Editing: "+imgPath)
  }
  cmd := exec.Command("open", "-W", imgPath)
  cmd.Run()
}

// Opens an image editor to mark up the image in the request.
// Returns URL to access the image after editing is completed.
func postHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("----- FIND: "+path.Base(r.URL.Path)+" -----")
  iconFile, _, err := r.FormFile("image")
  if err != nil {
    fmt.Println(err)
    http.Error(w, "Invalid request", http.StatusBadRequest)
    return
  }
  imgPath, err := writeImageFile(iconFile)
  if err != nil {
    fmt.Println(err)
    http.Error(w, "Internal server error", http.StatusInternalServerError)
    return
  }
  editImage(imgPath)
  if verbose {
    fmt.Println()
  }
  fmt.Fprintf(w, "http://"+hostname+":"+port+"/images/" + path.Base(imgPath))
}

// Serves the marked up images
func getHandler(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  resource := dir + r.URL.String()
  if verbose {
    fmt.Println("Serving: "+resource+"\n")
  }
  file, err := os.Open(resource)
  if err != nil {
      fmt.Println(err)
      http.Error(w, "Resource not found", http.StatusNotFound)
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
