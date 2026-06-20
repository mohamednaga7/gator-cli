# Gator

Gator is a CLI-based RSS feed aggregator.

## Prerequisites

To run Gator, you need to have the following installed on your system:

- **Go**: [Download and install Go](https://go.dev/doc/install) (version 1.26 or higher recommended).
- **PostgreSQL**: [Download and install PostgreSQL](https://www.postgresql.org/download/).

## Installation

You can install the Gator CLI using the `go install` command:

```bash
go install github.com/mohamednaga7/gator-cli@latest
```

Make sure your `GOBIN` directory (typically `~/go/bin`) is in your system's `PATH`.

## Configuration

Gator requires a configuration file located at `~/.gatorconfig.json`. This file stores your database connection string and the current active user.

Example structure of `~/.gatorconfig.json`:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": "your_username"
}
```

Replace `username`, `password`, and `gator` with your PostgreSQL credentials and database name.

## Commands

Here are some of the commands you can run with Gator:

- `register <username>`: Register a new user.
- `login <username>`: Login as an existing user.
- `addfeed <name> <url>`: Add a new RSS feed to follow (requires login).
- `feeds`: List all registered feeds.
- `follow <url>`: Follow an existing feed (requires login).
- `following`: List feeds you are currently following (requires login).
- `agg <time_between_reqs>`: Start the aggregator to fetch posts from feeds (e.g., `agg 1m`).
- `browse [limit]`: Browse posts from feeds you follow (requires login).
- `users`: List all registered users.
