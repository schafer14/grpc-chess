# Chess gRPC

## Goals

- Interface between chess engines and chess applications. 
- Learning gRPC/Knative

## Todo

- Think about technologies to use
  - Auth & Auth - Let's use auth0.com
  - Go
- Define the service (protobuf)
  - Firguring how it works
  - Doing (.proto)
- Build app
- Metrics & Logging
- Add authentication
- Deploy app

## What it should do

- There should be a home page
  - Watch featured game
  - Search for a game
  - Puzzels
  - Friends search

### Creating game process

- Player 1 search for opponent 
  - Include ratings filter, time control
- Player 2 is also searching for a game and criteria overlap
- Player 2 sees this game request and accepts it
- Player 1 should confirm 
- Both players should receive a game id

### Caveats 

- Software can act on either players behalf in the the above flow
- Either player may not be a human
  - But each agent should know the agent type of their opponent
  
## First Step

- Facilitate engine vs engine game
- 