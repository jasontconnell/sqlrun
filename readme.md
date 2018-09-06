_Usage_

sqlrun -dir "" -p comma,separated,prefixes

_Description_

Recursively search "dir" for sql files and run them against the database pointed to by the connection string in config.json

_config.json_

`{
    "connectionString": "connection string"
}`

_Notes_

For instance, if you wanted to make sure all table definition files run first, specify -p tbl    Multiple priorities can be set, otherwise it will order them by how it finds them in the directory structure.

Enjoy!