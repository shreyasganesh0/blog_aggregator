package main

import(
    "fmt"
    "github.com/shreyasganesh0/config" 
)

type Commands struct{
    name string 
    commnand_handler func() error
}
command_map = map[string]Command{
    "login": {name: "login",
             command_handler: Login(),
         },
    "register": {name: "register",
                command_handler: Register(),
            },
    "users": {name: "users",
            command_handler: GetUser(),
        },
    }
func main(){
    var conf config.Config
    
    if err := conf.Read(); err != nil {
        fmt.Printf("%s",err);
    }
    
    if err := conf.SetConfig("shreyas"); err != nil {
        fmt.Printf("%s",err);
    }

    if err := conf.Read(); err != nil {
        fmt.Printf("%s",err);
    }
    
}
