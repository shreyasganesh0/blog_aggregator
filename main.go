package main

import(
    "fmt"
    "os"
    "context"
    "database/sql"
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

func middlewareLogic(handler func(s *state, cmd command, user database.CheckUserRow) error) func(s *state, cmd command) error{
    return func(s *state, cmd command) error{
                user_name := s.conf.CurrentUserName;
                var user database.CheckUserRow
                user, err := s.queries.CheckUser(context.Background(), user_name); // this is a check for safety
                if err != nil {
                    return fmt.Errorf("User does not exist\n");
                }

                return handler(s, cmd, user); 
            }
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
    if err := cmds.register("addfeed", middlewareLogic(handlerAddFeed)); err != nil{
    
        return err;
    }
    if err := cmds.register("feeds", handlerFeed); err != nil{
    
        return err;
    }
    if err := cmds.register("following", middlewareLogic(handlerFollowing)); err != nil{
    
        return err;
    }
    if err := cmds.register("follow", middlewareLogic(handlerFollowEntry)); err != nil{
    
        return err;
    }
    if err := cmds.register("unfollow", middlewareLogic(handlerDeleteFeedFollow)); err != nil{
    
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
    if cmd.name != "agg" || cmd.name == "following"{
        cmd.args = args_cleaned[1:];
    }
    

    if err := cmds.run(&s, cmd); err != nil{
        fmt.Printf("%v", err);
        os.Exit(1);
    }

}

