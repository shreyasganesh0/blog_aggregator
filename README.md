# BLOG AGGREGATOR

This is a RSS blog aggregator written in go that uses postgres

## Requirements

Requires go and postgres as dependencies

1. Go version 1.23.5 is available to install [here](https://go.dev/doc/install)

2. Postgres version 15 can be installed [here](https://www.postgresql.org/download/)

## Usage

Commands available to user are

1. go run . login - login as existing user

2. go run . register name - register a user

3. go run . reset - reset the current state of the feeds and user database

4. go run . users - list users

5. go run . agg time (example go run . agg 1s) - aggregate posts from followed sites

6. go run . addfeed url - add feed to the list of available feeds

7. go run . feeds - list availabe feeds

8. go run . following - list feeds being followed by logged in users

9. go run . follow name - follow a feed by name

10. go run . unfollow name - unfollow a feed by name

11. go run . browse nooffeeds - list a number of posts from followed feeds

