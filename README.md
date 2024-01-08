# TCP-server with DDOS protection based on [PoW](https://en.wikipedia.org/wiki/Proof_of_work) algorithm

## PoW Algorithm

**PoW** implemented with a challenge-response protocol:

1. The client creates a tcp connection. The server starts to listening messages from client.
2. The client sends a message with the *`RequestChallenge`* command: `1:\n`.
3. The server generates a new puzzle with a [hashcash algorithm](https://en.wikipedia.org/wiki/Hashcash), stores it in the cache with TTL and sends a message with the *`ResponseChallenge`* command and puzzle: `2:puzzle\n`.
4. The client receives a puzzle and tries to compute a hash with enough number of zero bits in the beginning. Then the client requests a resource sending a message with the *`RequestResource`* command and solved puzzle: `3:solved-puzzle\n`.
5. The server receives the solved puzzle, checks TTL. If the puzzle was solved correctly, the server sends a message with the *`ResponseResource`* command and any resource: `4:some-resource\n`.

#### Supported commands:

* `0` - Error (server -> client)
* `1` - RequestChallenge (client -> server)
* `2` - ResponseChallenge (server -> client)
* `3` - RequestResource (client -> server)
* `4` - ResponseResource (server -> client)

## Getting Started
* Configure the server and client with environment variables:
  ```sh
  # Cache
  export CACHE_TTL=10
  
  # Server
  export SERVER_PORT=8080
  export SERVER_HOST=127.0.0.1
  export SERVER_MAX_ATTEMPTS=1000000
  export SERVER_HASH_BITS=3
  export SERVER_HASH_TTL=300
  export SERVER_TIMEOUT=10

  # Client
  export CLIENT_PORT=8080
  export CLIENT_HOST=127.0.0.1
  export CLIENT_MAX_ATTEMPTS=1000000
  ```
   
* Server and client run with docker-compose:
  ```sh
  docker-compose up
  ```
  
* Makefile commands:
  ```sh
  # Run tests
  make test
  ```
  ```sh
  # Build server binary
  make server-build
  ```
  ```sh
  # Run server with building binary
  make server-run
  ```
  ```sh
  # Build client binary
  make client-build
  ```
  ```sh
  # Run client with building binary
  make client-run
  ```