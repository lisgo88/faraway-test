# TCP-server with DDOS protection based on [PoW](https://en.wikipedia.org/wiki/Proof_of_work) algorithm

## PoW Algorithm

**PoW** implemented with a challenge-response protocol:

1. The client creates a tcp connection. The server generates a new puzzle with a [hashcash algorithm](https://en.wikipedia.org/wiki/Hashcash), and sends a message with the *`challenge`*.
2. The client receives a puzzle and tries to compute a hash with enough number of zero bits in the beginning and try to get resources.

**Protocol:**
The text-based protocol was changed to binary-based protocol for better performance.

## Getting Started
* Configure the server and client with environment variables:
  ```sh 
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