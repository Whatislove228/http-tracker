# Uptime Checker

A small Go CLI application that checks whether websites are reachable.

This is a learning pet project for practicing basic Go concepts:

- command-line arguments
- reading files
- working with slices
- error handling
- HTTP requests
- goroutines
- channels
- formatted output

## How It Works

The program reads a text file with URLs, checks each URL with an HTTP `GET` request, and prints the result.

Each result includes:

- URL
- status: `UP` or `DOWN`
- HTTP status code
- response time
- error message, if the request failed

HTTP status codes from `200` to `399` are treated as `UP`.
All other status codes are treated as `DOWN`.

## Project Structure

```txt
cmd/uptime-checker/main.go
sites.txt
README.md
