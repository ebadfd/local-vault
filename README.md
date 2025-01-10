# Local vault to store env

this tool will store the env values in a local vault. has a sqlite3 database which stores metadata about the project

basically allows to have a project eg:

slocoach. 

which can has multiple applications. eg:

- backend-api
- frontend

- initially user do need to configure the key he want to use, which is written in the config file

has 2 main functions. 


1. import existing env file. 

local-vault i --project slocoach --app backend-api --env dev --file .env

2. dump env values

local-vault d --project slocoach --app backend-api --env dev --file .env
> ask the user password? 


---

stores data like follows.

projects -> | id | project
apps -> | id | app | project_id | env | encrypted_file

---

file strucire. 

/local-env?
    - config.toml
    - data.sqlite3
    - data/
        - slocoach-prject-env.enc
