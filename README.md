# gator
RSS feed aggregator in Go

## Installation

To use gator you'll need to have the Go toolchain and Postgres installed.
Then you can install gator with the command:
```
go install github.com/ercorn/gator
```

## Configuration

Create a .gatorconfig.json file in your home director with the structure:
```
{
    "db_url":"postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
    "current_user_name":""
}
```
Replace the username, password, host, etc. values with your database connection string.

## Usage

Create a new user:
```
gator register <name>
```

Login as a specific user:
```
gator login <name>
```

Add a feed:
```
gator addfeed <url>
```

Follow an added feed as the current user:
```
gator follow <url>
```

Run the feed aggregator:
```
gator agg <time_duration_string (ex. 1s, 1m, 1m30s, 1h)>
```

Browse the aggregate posts corresponding to the current users followed feeds (default limit = 2):
```
gator browse [limit]
```
