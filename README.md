# Gator
## An RSS Feed aggregator
### From a boot.dev guided project

## Requirements
Requires Postgres and Go installed

## Installation instructions
To install run:
`go install github.com/davidw1457/gator`
Use goose to run the schema updates in sql/schema/*

## Config file
Create a config file named .gatorconfig.json in the root of your home
directory. The file has the following format:
`{
    "db_url":[database connection string],
    "current_user_name":[blank at start, this will be set by the application]
}`

## Available commands
Register a new user and switch to them with:
`gator register <username>`

Switch to a registered user with:
`gator login <username>`

Show all users with
`gator users`

Add and follow feed for current user:
`gator addfeed <name> <url>`

Follow feed added by other user:
`gator follow <url>`

Unfollow feed:
`gator unfollow <url>`

Show all feeds:
`gator feeds`

Show all feeds followed by current user:
`gator following`

Update all feeds:
`gator agg <time interval>

Show posts from followed feeds:
`gator browse [number of post to show. Default = 2]`
