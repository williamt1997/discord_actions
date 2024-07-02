## Prerequisites

Before running ensure go is installed:\
Built on:
- Go go1.22.3
- Postgresql DB
- Redis

## Installation

1. Clone the repository or download the zip file directly.
2. Create a .env file that contains variables need to run app (see .env.example)

## Usage

To run the app locally follow these steps:

1. Open a terminal or command prompt.
2. npm install
3. Run the following command:

```
   export ENV=development
   go run main.go
```

To run the app in docker, follow these steps:

1. docker build -t discordbot .
2. docker run -d --name {name} --env-file {env file} discordbot 
