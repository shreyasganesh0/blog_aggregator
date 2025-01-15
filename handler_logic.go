package main

import(
    "fmt"
    "os"
    "time"
    "context"
    "database/sql"
    "github.com/google/uuid"
    _ "github.com/lib/pq"
    "github.com/shreyasganesh0/blog_aggregator/database" 
)

func handlerLogin(s *state, cmd command) error{
    
    if len(cmd.args) == 0{
        fmt.Printf("Arguments to login empty\n");
        os.Exit(1);
    }
    user, err := s.queries.CheckUser(context.Background(), cmd.args[0]);

    if err == sql.ErrNoRows {
        fmt.Printf("User doesnt exist\n");
        os.Exit(1);
    }
    s.conf.SetConfig(user.Name);

    fmt.Printf("User Config set\n");
    
    return nil;
}

func handlerRegister(s *state, cmd command) error{

    if len(cmd.args) == 0{
        fmt.Printf("Arguments to register empty\n");
        os.Exit(1);
    }
    _, err := s.queries.CheckUser(context.Background(), cmd.args[0]);

    if err != sql.ErrNoRows {
        fmt.Printf("User exists, Login\n%v", err);
        os.Exit(1);
    }

    query_args := database.CreateUserParams{
        ID: uuid.New(),        
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Name: cmd.args[0],
    }
    s.queries.CreateUser(context.Background(), query_args);

    s.conf.SetConfig(cmd.args[0]);

    fmt.Printf("User: %v added to config file and db", cmd.args[0]);

    return nil;
}

func handlerReset(s *state, cmd command) error{
    fmt.Printf("Resetting Database tables");
    err := s.queries.DeleteAllUsers(context.Background());
    if err != nil {
        return err;
    }
    
    return nil;
}

func handlerUsers(s *state, cmd command) error{
    users, err := s.queries.GetUsers(context.Background());
    if err != nil{
        return err;
    }

    for _,user := range users {
        fmt.Printf("%s", user);
        if (user == s.conf.CurrentUserName){
            fmt.Printf(" (current)");
        }
        fmt.Printf("\n");
    }
    return nil;
}

func handlerAggregate(s *state, cmd command) error{

    url := "https://www.wagslane.dev/index.xml"
    rss_feed, err := fetchFeed(context.Background(), url);

    if err != nil{
        return err;
    }

    fmt.Printf("%v", *rss_feed);
    return nil;
}

func handlerAddFeed(s *state, cmd command, user database.CheckUserRow) error{
    if len(cmd.args) < 2 {
        fmt.Printf("Not Enough Args\n");
        os.Exit(1);
    }

    feed_name := cmd.args[0];
    feed_url := cmd.args[1];

    feed_id := uuid.New();        
    query_args := database.CreateFeedParams{
        ID: feed_id,        
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Name: feed_name,
        Url: feed_url,
        UserID: user.ID, 
    };
    _, err2 := s.queries.CreateFeed(context.Background(), query_args);

    if err2 != nil{
        return err2;
    }
    query_args_feed_follows := database.CreateFeedFollowParams{
        ID: uuid.New(),        
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        UserID: user.ID, 
        FeedID: feed_id, 
    }
    fmt.Printf("Feed added to user\n");

    _, err4 := s.queries.CreateFeedFollow(context.Background(), query_args_feed_follows);
    if err4 != nil{
        return err4;
    }

    fmt.Printf("Feed added to feed_follows\n");

    feed_fields, err3 := s.queries.FetchUserFeed(context.Background(), user.ID);
    if err3 != nil{
        return err3;
    }

    fmt.Printf("%v", feed_fields);

    return nil;
}

func handlerFeed(s *state, cmd command) error{

    feed_values, err := s.queries.FetchEntireFeed(context.Background());

    if err != nil{
        return err;
    }

    fmt.Printf("%v\n", feed_values);
    
    return nil;
}

func handlerFollowEntry(s *state, cmd command, user database.CheckUserRow) error{
    feed_url := cmd.args[0];

    feed_id, err := s.queries.FeedByUrl(context.Background(), feed_url);
    if err != nil{
        return err;
    }

    query_args := database.CreateFeedFollowParams{
        ID: uuid.New(),        
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        UserID: user.ID, 
        FeedID: feed_id, 
    };

    ret_v, err2 := s.queries.CreateFeedFollow(context.Background(), query_args);
    if err2 != nil{
        return err2;
    }

    fmt.Printf("username : %s, feedname : %s\n", ret_v[0].Name, ret_v[0].Name_2);
    return nil;
}

func handlerFollowing(s *state, cmd command, user database.CheckUserRow) error{

    feeds, err := s.queries.FeedFollowByUser(context.Background(), s.conf.CurrentUserName);
    if err != nil{
        return err;
    }

    fmt.Printf("username : %s, feeds : \t\n%s\n", s.conf.CurrentUserName, feeds);
    return nil;
}

