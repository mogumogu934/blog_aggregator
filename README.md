# blog_aggregator

Required:

PostgreSQL v16.8+
```
sudo apt update
sudo apt install postgresql postgresql-contrib
```

Golang v1.23+
```
curl -sS https://webi.sh/golang | sh
```


Install the blog_aggregator CLI
```
go install github.com/mogumogu934/blog_aggregator@latest
```

Before running the blog_aggregator CLI, the config file must be set up.
Create a file named ".gatorconfig.json" in your home directory with the following structure:

```
{
"db_url":"postgres://username:password@hostname:port/database_name?sslmode=disable",
"current_user_name":"your username"
}
```

Replace the values with your actual information:
```
db_url: Your PostgreSQL connection string
username: Your PostgreSQL username (default is often "postgres")
password: Your PostgreSQL password
hostname: Where your PostgreSQL server is running (usually "localhost")
port: The port PostgreSQL is running on (default is 5432)
database_name: The name of your database (e.g., "gator")
current_user_name: Your preferred username for the blog_aggregator
```

Example:
```
{
"db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
"current_user_name": "mogumogu934"
}
```

Go to the application's folder. The binary may be compiled:
```
go build
```

Or installed:
```
go install
```

Start the CLI:
```
./blog_aggregator"
```

Here is a list of available commands:
```
register <username>: Registers username and sets current user to username
login <username>: Sets current user to username
users: Prints a list of all registered users
addfeed <name> <url>: Adds a feed
feeds: Prints a list of all added feeds
follow <url>: Follows a feed
following: Prints a list of all followed feeds
unfollow <url>: Unfollows a feed
agg <duration>: Scrapes new blogs from all followed feeds every <duration>
browse <num>: Prints a list of <num> most recent posts from scraped feeds
```

Recommended order of commands for new users:
```
register <username>
addfeed <name> <url>
// Add more feeds
agg <duration>
browse <num>
```