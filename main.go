package main

import (
    "database/sql"
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strconv"
)

type routeHandler func(http.ResponseWriter, *http.Request) (interface{}, *httpError)

func (handlerFunction routeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.Header().Set("X-Frame-Options", "sameorigin")
    w.Header().Set("Cache-Control", "no-cache")
    // Note that the type of httpErr is httpError, not httpError
    if responseMap, httpErr := handlerFunction(w, r); httpErr != nil {
        log.Println(httpErr.Error)
        jsonErrorMessage, marshalErr := json.Marshal(map[string]string{"httpError": httpErr.Message})
        if marshalErr != nil {
            http.Error(w, "{\"httpError\": \"InternalServerError\"}", http.StatusInternalServerError)
            return
        }
        respondError(w, string(jsonErrorMessage), httpErr.Code)
    } else {
        jsonResponse, marshalErr := json.Marshal(responseMap)
        if marshalErr != nil {
            http.Error(w, "{\"httpError\": \"InternalServerError\"}", http.StatusInternalServerError)
            return
        }
        _, err := w.Write(jsonResponse)
        if err != nil {
            respondError(w, "{\"httpError\": \"InternalServerError\"}", http.StatusInternalServerError)
            return
        }
    }
}

func respondError(w http.ResponseWriter, error string, code int) {
    w.WriteHeader(code)
    _, _ = fmt.Fprintln(w, error)
}
func handleGo(w http.ResponseWriter, r *http.Request) (interface{}, *httpError) {
    switch r.Method {
    case "GET":
        key := r.URL.Query().Get("key")
        if len(key) == 0 {
            return nil, &httpError{Error: nil, Message: "requires key", Code: http.StatusNotFound}
        }

        db, err := getDb()
        if err != nil {
            return nil, &httpError{
                Error:   err,
                Message: "Failed to open database.",
                Code:    http.StatusInternalServerError,
            }
        }
        link, err := getLinkByKey(db, key)
        if err != nil {
            if err == sql.ErrNoRows {
                return nil, &httpError{Error: nil, Message: "Not found", Code: http.StatusNotFound}
            }
            return nil, &httpError{
                Error:   err,
                Message: "Failed to open database.",
                Code:    http.StatusInternalServerError,
            }
        }
        http.Redirect(w, r, link.url, http.StatusMovedPermanently)
        return "", nil
    default:
        return nil, &httpError{Error: nil, Message: "Method not allowed.", Code: http.StatusMethodNotAllowed}
    }
}

func handleLinks(w http.ResponseWriter, r *http.Request) (interface{}, *httpError) {
    switch r.Method {
    case "GET":
        db, err := getDb()
        if err != nil {
            return nil, &httpError{
                Error:   err,
                Message: "Failed to open database.",
                Code:    http.StatusInternalServerError,
            }
        }
        links, err := getLinks(db, -1, -1)
        if err != nil {
            return nil, &httpError{
                Error:   err,
                Message: "Failed to open database.",
                Code:    http.StatusInternalServerError,
            }
        }
        return links, nil
    case "POST":
        db, err := getDb()
        if err != nil {
            return nil, &httpError{
                Error:   err,
                Message: "Failed to open database.",
                Code:    http.StatusInternalServerError,
            }
        }

        config, err := getConfig()
        if err != nil {
            return nil, &httpError{
                Error:   err,
                Message: "Failed to fetch config",
                Code:    http.StatusInternalServerError,
            }
        }

        type linkParams struct {
            Key  string
            Url  string
        }

        var link linkParams
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            return nil, &httpError{
                Error:   err,
                Message: "Invalid parameters",
                Code:    http.StatusBadRequest,
            }
        }
        err = json.Unmarshal(body, &link)
        if err != nil {
            return nil, &httpError{
                Error:   err,
                Message: "Invalid parameters",
                Code:    http.StatusBadRequest,
            }
        }

        username, password, ok := r.BasicAuth()
        if username != config.Application.Username || password != config.Application.Password || !ok {
           return nil, &httpError{
               Error:   nil,
               Message: "Wrong credentials",
               Code:    http.StatusForbidden,
           }
        }

        _, err = addLink(db, link.Key, link.Url)
        if err != nil {
            return nil, &httpError{
                Error:   err,
                Message: "Invalid parameters",
                Code:    http.StatusBadRequest,
            }
        }

        return map[string]string{}, nil
    default:
        return nil, &httpError{Error: nil, Message: "Method not allowed.", Code: http.StatusMethodNotAllowed}
    }
}

func main() {
    const defaultConfigPath = "settings/development.yml"
    const defaultHostname = "localhost:8080"
    const maxSteps = 100
    const minSteps = -100

    c := flag.String("c", defaultConfigPath, "Config file path")
    flag.Parse()

    err := loadConfig(*c)
    if err != nil {
        log.Fatal(err)
    }

    switch flag.Arg(0) {
    case "migrate":
        switch flag.Arg(1) {
        case "gen":
            if len(flag.Arg(2)) == 0 {
                log.Fatal("requires revision name")
            }
            revisions, err := generateRevision(flag.Arg(2))
            if err != nil {
                log.Fatal(err)
            }
            fmt.Println("Generated", revisions)
        case "up":
            err := migrateUp()
            if err != nil {
                log.Fatal(err)
            }
            fmt.Println("Migrated: head")
        case "down":
            err := migrateDown()
            if err != nil {
                log.Fatal(err)
            }
            fmt.Println("Migrated: base")
        case "version":
            version, err := getVersion()
            if err != nil {
                log.Fatal(err)
            }
            fmt.Println("Version:", version)
        default:
            steps, err := strconv.ParseInt(flag.Arg(1), 10, 0)
            if err != nil {
                log.Fatal(fmt.Sprintf("Invaild migration command: %s", flag.Arg(1)))
            }
            if steps < minSteps || steps > maxSteps {
                log.Fatal("requires -100 <= steps <= 100")
            }
            err = migrateBySteps(int(steps))
            if err != nil {
                log.Fatal(err)
            }
            fmt.Println("Migrated:", steps, "steps")
        }
    case "run":
        hostname := defaultHostname
        if len(flag.Arg(1)) > 0 {
            hostname = flag.Arg(1)
        }
        logFile, err := getLogger(fmt.Sprintf("logs/custom-go-webserver.%s.log", config.Env))
        if err != nil {
            log.Fatal(err)
        }
        defer logFile.Close()

        db, err := getDb()
        if err != nil {
            log.Fatal(err)
        }
        defer db.Close()

        http.Handle("/", routeHandler(func(http.ResponseWriter, *http.Request) (interface{}, *httpError) {
            return nil, &httpError{Error: nil, Message: "Page not found", Code: http.StatusNotFound}
        }))
        http.Handle("/links", routeHandler(handleLinks))
        http.Handle("/go", routeHandler(handleGo))

        fmt.Println("Listening on:", hostname)
        log.Fatal(http.ListenAndServe(hostname, nil))
    default:
        log.Fatal("Invaild command:", flag.Arg(0))
    }
}
