package config

import (
    "encoding/json"
    "os"
    "io"
    "fmt"
)

var path string

func init(){
    home_dir, err := os.UserHomeDir();
    if err != nil {
        fmt.Printf("home dir err");
    }
    path = home_dir+"/.gatorconfig.json"
}    

type Config struct{
    Dburl string `json:"db_url"`
    CurrentUserName string `json:"current_user_name"`
}


func (c *Config) Read() error{
    file, err := os.Open(path);
    if err != nil{
        return err;
    }
    defer file.Close();
    
    text, err2:= io.ReadAll(file);
    if err2 != nil {
        return err2;
    }

    if err := json.Unmarshal(text, c); err != nil{
        return err;
    }

    fmt.Printf("%v", string(text));
    return nil;
}

func (c *Config) SetConfig(current_user_name string) error{

    c.CurrentUserName = current_user_name;

    config_marshall, err := json.MarshalIndent(c, "", " ");
    if err != nil{
        return err;
    }
    err2 := os.WriteFile(path, config_marshall, 0644);
    if err2 != nil {
        return err;
    }

    return nil;
}




    


