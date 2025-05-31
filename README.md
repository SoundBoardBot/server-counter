# Guild Counter Service

The service that keeps the count of guilds to be displayed on bot listings.

## How to run

### Production

1. Clone the repository (`git clone https://github.com/SoundBoardBot/server-counter.git`)
2. Run `docker build -t guild-counter-service .`
3. Run `docker run -d -e LOG_LEVEL=info guild-counter-service`

### Development

1. Clone the repository (`git clone https://github.com/SoundBoardBot/server-counter.git`)
2. Run `go mod download`
3. Run `go mod verify`
4. Run `go run main.go`
