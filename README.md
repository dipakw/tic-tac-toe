# Tic-Tac-Toc

This is a simple networked played **Tic-Tac-Toc** game written in Go. There will be one central game server running on a TCP port with an HTTP and WebSocket API. Clients connect to the server and provide an web UI for the peers to play, and thier moves reflect to each other via WebSocket/SSE without a refresh.


## Usage

### Build binaries

Run the build script first to get binaries, it will build into the bin folder.

```bash
./build.sh
```

### Run the server

Go to the bin folder and run below. The default port for the server is 16000, so passing `-port` argrument is optional.

```bash
./server -port 16000
```

### Run clients

You should run multiple instances of clients, at least 2 to play with each other.

Running a client will start a service on a random available port, so they won't collide.

Both `-server-addr` and `-server-port` arguments are optional, the default values are `-server-addr 0.0.0.0` and `-server-port 16000`

```bash
./client -server-addr 127.0.0.1 -server-port 16000
```

It will print the link like.

```txt
Client will connect to the server: 127.0.0.1:16000
Play game at: http://localhost:58010
```

Now, click the `Play game at` link where you will get the interface to enter the peer ID and play with them.