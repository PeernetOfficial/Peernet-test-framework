package UDT_benchmark

//
//import (
//    "fmt"
//    "github.com/PeernetOfficial/Peernet-test-framework"
//    "github.com/gorilla/mux"
//    "log"
//    "net/http"
//    "time"
//)
//
//func UDTBenchmark(config *testframework.Config) {
//    r := mux.NewRouter()
//
//    // Get Config information
//
//    srv := &http.Server{
//        Handler: r,
//        Addr:    config.MainServerAddress,
//        // Good practice: enforce timeouts for servers you create!
//        WriteTimeout: 15 * time.Second,
//        ReadTimeout:  15 * time.Second,
//    }
//
//    manager, err := config.RunManager()
//    if err != nil {
//        fmt.Println(err)
//    }
//
//    // Simulate file transfer
//    // The objective here is to simulate file transfer
//    // from nodes to multiple nodes in a network
//
//    manager.
//
//        // create files testing
//        Abstrations.Touch(manager[0].NodeConfig, "example.go")
//
//    fmt.Println(len(*manager))
//
//    // Lister for the main server
//    log.Fatal(srv.ListenAndServe())
//}
