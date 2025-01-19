# BLOG AGGREGATOR

This is a RSS blog aggregator written in go that uses postgres

## Requirements

Requires go and postgres as dependencies

1. Go version 1.23.5 is available to install [here](https://go.dev/doc/install)

2. Postgres version 15 can be installed [here](https://www.postgresql.org/download/)

## Usage

Commands available to user are

1.  login - login as existing user

2.  register name - register a user

3.  reset - reset the current state of the feeds and user database

4.  users - list users

5.  agg time (example go run . agg 1s) - aggregate posts from followed sites

6.  addfeed url - add feed to the list of available feeds

7.  feeds - list availabe feeds

8.  following - list feeds being followed by logged in users

9.  follow name - follow a feed by name

10.  unfollow name - unfollow a feed by name

11.  browse nooffeeds - list a number of posts from followed feeds

## Installation

Build the code from source
```
go install
go build
```



