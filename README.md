# Gator

An RSS feed aggregator built in Go. Made as a guided project within [Boot Dev](https://www.boot.dev/courses/build-blog-aggregator-golang)

## To Run

Before you run make sure you have Go and Postgres installed

### Install

```bash
go install ...
```

### Config

Create `.gatorconfig.json` in your home directory `~/`
Insert this in the file and fill out the db_url with your information

```json
{
  "db_url": "postgres://username:@localhost:5432/database?sslmode=disable"
}
```

### Run

1. Start by creating a user by running:

```bash
gator register <username>
```

2. Add a feed:

```bash
gator addfeed <url>
```

3. Start the aggregator (use timestamps like 1s, 1m, 1h):

```bash
gator agg <time_between_requests>
```

4. View the posts (limit is optional):

```bash
gator browse [limit]
```

### All Commands

- `login <username>` - login to a user
- `register <username>` - create a user account
- `reset` - reset the entire database
- `users` - view all registered users
- `agg <time_between_requests>` - scrape from list of added feeds
- `addfeed <name> <url>` - add a feed to list of feeds to aggragate
- `feeds` - view a list of added feeds
- `follow <url>` - follow a feed
- `following` - view a list of followed feeds
- `unfollow <url>` - unfollow a feed
- `browse [limit]` - view aggragated posts

## Tools Used

- Goose
- SQLC
- Postgres
