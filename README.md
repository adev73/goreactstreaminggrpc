# Go + React Streaming gRPC example

This sample program shows, using a simple greet service, how to implement a server streaming gRPC
call using buf.build's connect ecosystem.

# Pre-requisites

- Go v1.18 or newer
- Node 16+ with npm (or better still, pnpm)
- Curl
- A web browser
- Linux (or mac? or Windows maybe?) -- I've only tested it on Linux.

# Build and start the front-end

From the repo root directory:

```
$ cd react-app
$ pnpm i
$ pnpm start
```

# Run the backend server

From the repo root directory, run `go mod tidy` or `go get -u ./...` to grab the dependencies; then:

```
$ go run cmd/server/main.go
```

You'll see a `Starting server...` message, then `There are 0 sessions currently active` repeated every 10s.

There's also a very basic Go client (`go run cmd/client/main.go`) which slowly asks for 4 greetings before quitting.

# Connect the front end

On the web page, click Connect. The server should output `Session <uuid> started.`, followed by `There are 1 sessions currently active.`

Type a name in the "greet someone box" and click "Greet them" - you'll almost instantly see the greeting

# Use the same session ID in Curl

For the purposes of this demo, we're using a simple session ID to control who sees what events. We can prove that the streaming
is working from any client, by making a request via cURL. Assume a session ID of `7f14798f-4bb0-4db0-acea-a9db84bd2884`, make the
following cURL request:

```
curl \
    --header "Content-Type: application/json" \
    --data '{"sessionId":"7f14798f-4bb0-4db0-acea-a9db84bd2884","name":"Curly"}' \
    http://localhost:8080/greet.v1.GreetService/Greet

```
You'll see the output `{"confirmed":true}` in your console window; and on the web page, `Hello, Curly! It's <date> <time>`

Click "Disconnect" to shut down the session cleanly. Or, just stop the pnpm process; which will simply drop the connection dead. 
Eventually, the server will notice it's PINGs aren't working, and will shut down the session.

# So, why?

This is a self-training project, working towards a multi-client event sourcing system - where an event store will service subscriptions as
server streaming gRPC requests. The GO side of stuff is (for me) pretty easy; getting the Typescript client to behave was, well, harder...

But it makes a neat little demo :)

# Modifying the protobuf definitions

If you want to change the protobufs, you'll need the [buf.build CLI](https://buf.build/product/cli/) which - at the time of writing, at least -
is free to use for local code generation.

Simply modify the proto, then type `buf generate`, and buf will re-build your Go and Typescript libraries.
