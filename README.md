Gator is a CLI tool that allows users to:

Add RSS feeds from across the internet to be collected
Store the collected posts in a PostgreSQL database
Follow and unfollow RSS feeds that other users have added
View summaries of the aggregated posts in the terminal, with a link to the full post

This program requires that you have Postgres and Go installed

You can then use "go install gator" to install the CLI

At the root of your home directory, you will need to create a config file titled ".gatorconfig.json"
This is the file that the program will use to track the current logged in user

A few of the commands you will use are

gator login <NAME> - this will set you as the current user in the config file effectively "logging you in"

gator addfeed <FEED NAME> <FEED URL> - this command allows you to add a new feed to the tool and sets you as following that feed

gator unfollow <FEED URL> - this command allows you to unfollow a feed by URL

gator agg <TIME BETWEEN REQUESTS> - this commmand aggregates feeds and saves them to the database so they can be browsed using the brows command