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
)

func main() {
  startServer()
}

func getFormat(fileBytes []byte) (string) {
  bytes := fileBytes[0:3]
  if bytes[0] == 0x89 && bytes[1] == 0x50 && bytes[2] == 0x4E && bytes[3] == 0x47 { return "png" }
  if bytes[0] == 0xFF && bytes[1] == 0xD8 { return "jpg" }
  if bytes[0] == 0x47 && bytes[1] == 0x49 && bytes[2] == 0x46 && bytes[3] == 0x38 { return "gif" }
  if bytes[0] == 0x42 && bytes[1] == 0x4D { return "bmp" }
  return ""
}

func postHandler(w http.ResponseWriter, r *http.Request) {
  iconFile, _, err := r.FormFile("image")
  if err != nil {
    fmt.Println("Failed, woops")
    fmt.Println(err)
  }

  dir := "/Users/user/Desktop"
  os.Mkdir(dir + "/images", 0777)
  iconBytes, err := ioutil.ReadAll(iconFile)
  if err != nil {
    fmt.Println(err)
  }
  iconFile.Close()
  imgType := getFormat(iconBytes)
  imgName := strconv.FormatInt(time.Now().Unix(), 10) + "." + imgType
  imgPath := dir + "/images/" + imgName
  err = ioutil.WriteFile(imgPath, iconBytes, 0777)
  if err != nil {
    fmt.Println(err)
  }
  cmd := exec.Command("open", "-W", imgPath)
  cmd.Run()
  fmt.Fprintf(w, "http://localhost:8080/images/" + imgName)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("in getHandler...")
  r.ParseForm()

  dir := "/Users/user/Desktop"
  file, err := os.Open(dir + r.URL.String())
  fmt.Println("path: " + dir + r.URL.String())
  if err != nil {
      http.Error(w, err.Error(), http.StatusNotFound)
      return
  }

  modeTime := time.Now()
  http.ServeContent(w, r, dir+r.URL.String(), modeTime, file)
}

func startServer() {
  mux := http.NewServeMux()
  mux.HandleFunc("/test", postHandler)
  mux.HandleFunc("/images/", getHandler)
  http.ListenAndServe(":8080", cors.Default().Handler(mux))
}
