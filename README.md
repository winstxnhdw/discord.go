# discord.go

Existing Discord API wrappers do not have support for Discord self-bots. This is because self-bots are against Discord's [Terms of Service](https://support.discord.com/hc/en-us/articles/115002192352). However, self-bots are still useful for automating tasks that would otherwise be tedious to do manually. This library aims to provide a simple, easy-to-use interface for interacting with the Discord API.

## Installation

```bash
go get -u github.com/winstxnhdw/discord.go
```

## Usage

Create a Discord client with your personal Discord token.

```go
client := discord.Create(os.Getenv("DISCORD_TOKEN"))
defer client.Dispose()
```

By default, every request is sent as soon as possible. You can set a delay between requests to avoid rate limiting.

```go
client := discord.Create(os.Getenv("DISCORD_TOKEN"), time.Second*2)
```

You can send a message to a channel. Every request returns a cancel function that can be used to cancel the request. This is useful for sending multiple messages without waiting for a response.

```go
response, cancelFunction, err := client.MessageChannel(
    os.Getenv("DISCORD_CHANNEL_OR_USER_ID"),
    "Hey, did you know that direct messages use the same API as channels?",
)
```
