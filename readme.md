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
| `GET /pwd` | - | `{"outcome": "...", "cwd": "path/cwd"}`|
| `POST /cd` | `{"path": "pathvalue"}` | `{"outcome": ".."}`|
| `POST /mkdir` | `{"path":"pathvalue"}` | `{"outcome": ".."}` |
| `POST /rmdir` | `{"path":"pathvalue"}`| `{"outcome": ".."}` |
| `POST /store` | `{"filename":"name", "file": {...} }`| `{"outcome": ".."}` |
| `GET /retrieve` | `{"filename":"name"}` | `{"outcome": "..", "file":{...}}` |

Actions `mkdir, rmdir, store, retrieve` have effect in the current working directory.
