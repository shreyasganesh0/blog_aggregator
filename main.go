package main

import(
    "fmt"
    "os"
    "time"
    "context"
    "database/sql"
    "github.com/google/uuid"
    _ "github.com/lib/pq"
    "github.com/shreyasganesh0/config" 
    "github.com/shreyasganesh0/blog_aggregator/database" 
)

// type defs
type command struct{
    name string 
    args []string
}

type state struct{
    queries *database.Queries
    conf *config.Config
}

type commands struct{
    command_map map[string]func(*state, command) error
}

var cmds commands;

// extern funcs
func init(){
    
    command_map := make(map[string]func(*state, command) error);
    cmds.command_map = command_map;

}

func (c *commands) register(name string, f func(*state, command) error) error{
    c.command_map[name] = f;
    return nil;
}

func (c *commands) run(s *state, cmd command) error{
    command_func, exists := c.command_map[cmd.name];

    if !exists{
        return fmt.Errorf("Command doesnt exist");
    }

    command_func(s, cmd);
    return nil;
}

func handlerLogin(s *state, cmd command) error{
    
    if len(cmd.args) == 0{
        fmt.Printf("Arguments to login empty\n");
        os.Exit(1);
    }
    user_name, err := s.queries.CheckUser(context.Background(), cmd.args[0]);

    if err == sql.ErrNoRows {
        fmt.Printf("User doesnt exist\n");
        os.Exit(1);
    }
    s.conf.SetConfig(user_name);

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

func startUp(s *state) error{

    dbUrl := s.conf.Dburl;
    
    db, err := sql.Open("postgres", dbUrl);
    if err != nil {
        return err;
    }

    dbQueries := database.New(db);
    s.queries = dbQueries;

    // add new commands here
    if err := cmds.register("login", handlerLogin); err != nil{
    
        return err;
    }
    if err := cmds.register("register", handlerRegister); err != nil{
    
        return err;
    }
    if err := cmds.register("reset", handlerReset); err != nil{
    
        return err;
    }
    if err := cmds.register("users", handlerUsers); err != nil{
    
        return err;
    }
    if err := cmds.register("agg", handlerAggregate); err != nil{
    
        return err;
    }
    return nil;
    
}

//main
func main(){
    var s state;
    var c config.Config;
    var cmd command;


    s.conf = &c;
    args_cleaned := os.Args[1:];

    if err := s.conf.Read(); err != nil{ // reads state from conf file and loads it into s.conf
        fmt.Printf("%v", err);
    }
    if err := startUp(&s); err != nil{
        fmt.Printf("Startup error: %v", err);
    }
    
    cmd.name = args_cleaned[0];
    if cmd.name != "agg"{
        cmd.args = args_cleaned[1:];
    }
    

    if err := cmds.run(&s, cmd); err != nil{
        fmt.Printf("%v", err);
        os.Exit(1);
    }

}

