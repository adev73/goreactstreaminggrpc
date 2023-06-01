import React, { useRef, useState } from 'react';
import logo from './logo.svg';
import './App.css';

// Buf stuff
import { createPromiseClient } from '@bufbuild/connect'
import { createConnectTransport } from '@bufbuild/connect-web'
import { GreetService } from './gen/greet/v1/greet_connect'
import { GreetingsRequest } from './gen/greet/v1/greet_pb'



function App() {
  const [responses, setResponse] = useState<string[]>([])
  const [greet, setGreet] = useState<string>("")
  const [sessionId, setSessionId] = useState<string>("")

  const refSessionId = useRef<string>()
  refSessionId.current = sessionId

  const client = createPromiseClient(
    GreetService,
    createConnectTransport({
      baseUrl: "http://localhost:8080"
    })
  )

  const handleGreetChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setGreet(event.target.value)
  }

  const handleGreet = () => {
    if (sessionId === "") {
      console.log("No active session...")
    } else {
      console.log(`Greeting on session ${sessionId}`)
      client.greet({
        sessionId: sessionId,
        name: greet
      })
    }
  }
  const handleConnect = async () => {
    if (sessionId !== "") {
      console.log(`Already connected to session ${sessionId}!`)
      return
    }

    console.log("Connecting...")
    const request = new GreetingsRequest({})
    for await (const response of client.greetings(request)) {      
      if (response.greeting === "") {
        if (response.endSession) {
          console.log(`Session ${refSessionId.current} ended.`)
          break
        }
        if (refSessionId.current !== "") {
          console.log(`PINGed by session ${refSessionId.current}`)
          continue
        }
        // If we get here, we're starting a new session.
        setSessionId(response.sessionId)
        console.log(`New session ${response.sessionId} connected`)
      } else {        
        console.log(`Received greeting: ${JSON.stringify(response)}`)
        setResponse((resp) => [
          ...resp,
          response.greeting,
        ])
      }
    }
    console.log("Exited await for loop.")
    setSessionId("")
    setResponse([])
}

  const handleDisconnect = () => {
    console.log("Disconnecting...")
    client.greet({
      sessionId: sessionId,
      endSession: true,
    })
  }

  return (
    <div>
      <div>
        <h1>Greeter</h1>
        <h4>Part of learnings for the Event Source POC</h4>
      </div>
      <hr />
      <div>
        <h5>Greetings to:</h5>
        <div>
          {responses.map((greetings, i) => {
            return (
              <div key={`greetings${i}`}>
                <p>{greetings}</p>
              </div>
            )
          })}
        </div>
        <hr />
        <div>
          <h5>Greet Someone</h5>
          <input type='input' value={greet} onChange={handleGreetChange}></input>
          <input type='button' value="Greet them" onClick={handleGreet}></input>
        </div>
        <hr />
        <div>
          <h5>Control</h5>
          <input type='button' value="Connect" onClick={handleConnect}></input>
          <input type='button' value="Disconnect" onClick={handleDisconnect}></input>
        </div>
      </div>
    </div>
  );
}

export default App;
