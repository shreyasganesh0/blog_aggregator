package main

import(
    "fmt"
    "github.com/shreyasganesh0/config" 
    "os"
)

type command struct{
    name string 
    args []string
}

type state struct{
    conf *config.Config
}

type commands struct{
    command_map map[string]func(*state, command) error
}

var cmds commands;

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

func main(){
    var args_cleaned []string
    var s state;
    var c config.Config;

    s.conf = &c;

    args_cleaned = os.Args[1:];


    if err := s.conf.Read(); err != nil{
        fmt.Printf("%v", err);
    }

    var cmd command
    
    cmd.name = args_cleaned[0];
    cmd.args = args_cleaned[1:];
    
    if len(cmd.args) == 0{
        fmt.Printf("Arguments to login empty");
        os.Exit(1);
    }




    if err := cmds.register("login", handlerLogin); err != nil{
    
        fmt.Printf("%v", err);
        os.Exit(1);
    }

    if err := cmds.run(&s, cmd); err != nil{
        fmt.Printf("%v", err);
        os.Exit(1);
    }

}

func handlerLogin(s *state, cmd command) error{
    

    s.conf.SetConfig(cmd.args[0]);

    fmt.Printf("User Config set");
    
    return nil;
}
