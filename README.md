# db-struct-generator project
Generate structs from database.  
This program connects to a database and retrieves the informations about a table to generate a go struct corresponding to the table.

## usage
```
db-struct-generator generate <user> <host> <port> <table>
```


## Dependencies
- make 
- go
- pandoc for user manual generation

## make commands
- `make all` : builds for Windows, linux generic and Ubuntu/debian (deb), builds the manuals and put everything in the `publish` folder
- `make` : builds for your current platform


## Folder structure
### configs
config file
### internal
Code that won’t be reusable in other projects
### pkg 
Code that can be reused in other projects like log system
### assets
Files used in the program