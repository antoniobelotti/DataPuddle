# Abstract Interface

The data puddle interface is defined as such 

    pwd(): string, outcome
    cd(path string): outcome
    mkdir(path string): outcome
    rmdir(path string): outcome
    store(file jsonfile, filename string): outcome
    retrieve(filename): jsonfile | outcome

outcome is a message which can be "error ..." or "ok"

# Exposed WebApi Interface

| REQUEST | BODY | RESPONSE |
|---|---|---|
| `GET /sessionkey` | - | `{"outcome": "...", "key": "..."}`| 
| `GET /pwd?key={key}` | - | `{"outcome": "...", "cwd": "path/cwd"}`|
| `GET /cd?key={key}&path={path}` | - | `{"outcome": ".."}`|
| `GET /mkdir?key={key}&path={path}` | - | `{"outcome": ".."}` |
| `GET /rmdir?key={key}&path={path}` | - | `{"outcome": ".."}` |
| `POST /store?key={key}&filename={name}` | jsonfile content `{...}`| `{"outcome": ".."}` |
| `GET /retrieve?key={key}&filename={name}` | - | `{"outcome": "..", "file":{...}}` |

Each client has to request a session key and include it in every subsequent request. In this way the server can track
the current working directory for multiple connections. 

Actions `mkdir, rmdir, store, retrieve` have effect in the current working directory.
