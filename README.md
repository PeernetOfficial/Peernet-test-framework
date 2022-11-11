# Peernet Test Framework

## Objective:
- To generate a swarm of Peernet nodes
- Setup custom Root and regular nodes
- Simulate discovery and file transfer

The following library spawns peernet nodes automatically and runs 
each Peernet instance as a go routine. 

The setup is intended to be as simple as possible.
```
Go build .
```
Run:
```
./Peernet-test-framework
```

Extend to your Go project (Instructions coming soon)
```
// import Test framework 

// Control config through function calls 

// Generate network and return Configs as a list 

// Manipulate each Peernet instance -> Reference WebAPI 
(Using as a pointers as they can be updated)
// -> Soon replaced as abstractions
```