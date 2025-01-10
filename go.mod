module github.com/shreyasganesh0/blog_aggregator

go 1.23.4

replace (
        github.com/shreyasganesh0/config v0.0.0 => ./internal/config
        github.com/shreyasganesh0/blog_aggregator/database v0.0.0 => ./internal/database
    )

require (
	github.com/google/uuid v1.6.0 // indirect
    github.com/shreyasganesh0/config v0.0.0
	github.com/lib/pq v1.10.9 // indirect
    github.com/shreyasganesh0/blog_aggregator/database v0.0.0
)
