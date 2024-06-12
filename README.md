# Backup Databases (2024-06-12)
# Releases
* **2024-06-12** Release version 1.0.0 support Microsoft SQL Server, PostgreSQL and MySQL

## Getting started
### Installing
To install project, you need to clone the project by following command below:

```bash
$ git clone https://<Username>:<Personal access token>@github.com/navamin11/Database-Backup.git backupdb
```

## Configure your databases
### Copy the template of environment to actual one (we'll use this)

```bash
$ cp config.json.example config.json
```

### Setup config.json
```bash
{
    "database": [
        {
            "driver":"<sqlserver|postgres|mysql>",
            "host": "<ip address>",
            "port": "<port>",
            "username": "<user>",
            "password": "<pass>",
            "schema": "<schema>",
            "dbName":"<database>",
            "tables": [
                {
                    "select": "*",
                    "tbName": "actor"               
                }
            ]
        }
    ]
}
```

## Running project
### Run this command to start application or build application

```bash
$ make run
```
or
```bash
$ make build
```

> [!NOTE]
> After build, Please check binary file and config.json to bin.

### If using docker to initialize surrounding program (postgres, mysql, myserver)

```bash
$ make rebuild-all
```

> [!NOTE]
> See detail to Makefile.