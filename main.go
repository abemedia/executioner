package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "os/exec"
    "net/http"
    "github.com/julienschmidt/httprouter"
)

func auth(h httprouter.Handle) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
        if r.FormValue("api_key") == Config.Secret {
            h(w, r, p)
        } else {
            http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
        }
    }
}

func Index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
    json.NewEncoder(w).Encode(Config.Endpoints)
}

func Execute(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
    run := Config.CMD[p.ByName("slug")]
    cmd := exec.Command(run[0], run[1:]...)
    var stdout bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    err := cmd.Run()
    if err != nil {
        var e string
        if err.Error() != "exit status 1" {
            e = err.Error()
        }
        if stderr.String() != "" {
            if e != "" {
                e += ": "
            }
            e += stderr.String()
        }
        http.Error(w, e, http.StatusInternalServerError)
        log.Print(e)
    }

    if res := stdout.String(); res != "" {
        w.Write([]byte(res))
        log.Println(res)
    }
}

func init() {
    // swag
    fmt.Println(`
 _____                         _    _
|  ___|                       | |  (_)
| |__ __  __ ___   ___  _   _ | |_  _   ___   _ __    ___  _ __
|  __|\ \/ // _ \ / __|| | | || __|| | / _ \ | '_ \  / _ \| '__|
| |___ >  <|  __/| (__ | |_| || |_ | || (_) || | | ||  __/| |
\____//_/\_\\___| \___| \__,_| \__||_| \___/ |_| |_| \___||_|

================================================================

Launching Executioner...`)
}

func main() {
    // log to file
    f, err := os.OpenFile(Config.LogPath, os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
    if err != nil {
        fmt.Println("Error creating log file:", err)
        fmt.Println("WARNING: Error logs will only be displayed here!")
    } else {
        log.SetOutput(f)
    }
    defer f.Close()

    router := httprouter.New()
    router.GET("/", auth(Index))
    router.GET("/:slug", auth(Execute))

    host := fmt.Sprintf("%s:%d", Config.Host, Config.Port)
    log.Println("Executioner launched on", host)
    log.Fatal(http.ListenAndServe(host, router))
}