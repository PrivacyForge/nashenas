# Nashenas
Nashenas is a anonymous messaging bot for Telegram users, designed for security purposes. It allows users to send encrypted messages to each other without revealing their identities. Nashenas uses hybrid encryption (RSA + AES) similar to TLS, ensuring that messages remain secure and protected against Man-in-the-Middle (MitM) attacks. The goal is to create a safe and enjoyable platform for anonymous communication. The chart below explains the core mechanism of Nashenas in the simplest way:

![chart](https://github.com/PrivacyForge/nashenas/blob/main/chart.png)

## Deployment
This section provides the necessary guidelines for deploying the project.

### Redis
We send message notifications using Redis, so you need to set up an instance.

### Front-end
Set env variables at `.env`.

install dependencies:
```sh
bun install
```
to build:
```sh
bun run build
```

### Back-end
Set env variables at `./server/.env`.

install dependencies:
```sh
go install
```
to run:
```sh
go run main.go
```

### Telegram bot
Set env variables at `./bot/.env`.

install dependencies:
```sh
bun install
```
to run:
```sh
bun run src/index.ts
```
