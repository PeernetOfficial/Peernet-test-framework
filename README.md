# Peernet Test Framework

## Objective:
- To generate a swarm of Peernet nodes
- Setup custom Root and regular nodes
- Simulate discovery and file transfer

The following library spawns peernet nodes automatically and runs 
each Peernet instance as a go routine. The library also 
automatically generates root nodes and regular nodes and also 
automatically points the regular nodes to root node. 

The setup is intended to be as simple as possible.
```
Go build .
```
Run:
```
./Peernet-test-framework
```

Extend to your Go project (Sample Program to spawn Peernet nodes based on default settings)
```go
import (
    "github.com/PeernetOfficial/Peernet-test-framework"
    "github.com/gorilla/mux"
)

func main() {

r := mux.NewRouter()

// Get Config information
config, err := testframework.ConfigInit()
if err != nil {
  fmt.Println(err)
}

srv := &http.Server{
  Handler: r,
  Addr:    config.MainServerAddress,
  // Good practice: enforce timeouts for servers you create!
  WriteTimeout: 15 * time.Second,
  ReadTimeout:  15 * time.Second,
}

manager, err := config.RunManager()
if err != nil {
   fmt.Println(err)
}

fmt.Println(len(*manager))

// Lister for the main server
log.Fatal(srv.ListenAndServe())
}
```