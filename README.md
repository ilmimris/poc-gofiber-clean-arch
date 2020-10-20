# poc-gofiber-clean-arch

## Description
This is a POC of implementation Clean Architecture in Go (Golang) projects using Gofiber.io framework.

Rule of Clean Architecture by Uncle Bob
 * Independent of Frameworks. The architecture does not depend on the existence of some library of feature laden software. This allows you to use such frameworks as tools, rather than having to cram your system into their limited constraints.
 * Testable. The business rules can be tested without the UI, Database, Web Server, or any other external element.
 * Independent of UI. The UI can change easily, without changing the rest of the system. A Web UI could be replaced with a console UI, for example, without changing the business rules.
 * Independent of Database. You can swap out Oracle or SQL Server, for Mongo, BigTable, CouchDB, or something else. Your business rules are not bound to the database.
 * Independent of any external agency. In fact your business rules simply don’t know anything at all about the outside world.

More at https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html

This project has  4 Domain layer :
 * Models Layer
 * Repository Layer
 * Usecase Layer  
 * Delivery Layer

 #### The diagram:
 ![golang clean architecture](https://github.com/bxcodec/go-clean-arch/raw/master/clean-arch.png)

#### The folder structure:
```bash
.
├── api
│   └── app.go
│
├── pkg
│   ├── common
│   │   ├── delivery
│   │   └── repository
│   │       └── helper.go
│   │
│   ├── domain
│   │   ├── author.go
│   │   ├── post.go
│   │   ├── errors.go
│   │   └── mocks
│   │       ├── AuthorRepository.go
│   │       ├── AuthorUsecase.go
│   │       ├── PostRepository.go
│   │       └── PostUsecase.go
│   │
│   ├── author
│   │   └── repository
│   │       └── psql
│   │           ├── psql_repository.go
│   │           └── psql_repository_test.go
│   │
│   └── post
│       ├── delivery
│       │   └── rest
│       │       ├── post_rest.go
│       │       └── post_rest_test.go
│       ├── repository
│       │   └── psql
│       │       ├── psql_repository.go
│       │       └── psql_repository_test.go
│       └── usecase
│           ├── post_usecase.go
│           └── post_usecase_test.go

```

#### The explanation for this folder structure:
- `api` folder contains all main app code, such as setup connection to db and rest.
- `pkg` folder contains all the modules, such as:
    - `common` module (helper, middleware, etc.)
    - `domain` module, where the domain or entity define as well as the interface (port) for repository and usecase contract 
    - `author` module, where the repository, usecase, and delivery of author defined
    - `post` module, where the repository, usecase, and delivery of post defined

> Author, post, and other module could be tested separately


### How To Run This Project
> Make Sure you have run the post_mysql.sql in your mysql

or
> Make Sure you have run the post_psql.sql in your postgres


Since the project already use Go Module, I recommend to put the source code in any folder but GOPATH.

#### Run the Testing

```bash
$ make test
```

#### Run the Applications
Here is the steps to run it with `docker-compose`

```bash
#move to directory
$ cd workspace

# Clone into YOUR $GOPATH/src
$ git clone https://github.com/ilmimris/poc-gofiber-clean-arch.git

#move to project
$ cd poc-gofiber-clean-arch

# Build the docker image first
$ make docker

# Run the application
$ make run

# Run the application using mysql db
$ make run-mysql

# Run the application using postgres db
$ make run-postgres

# check if the containers are running
$ docker ps

# Execute the call
$ curl localhost:8080/post

# Stop
$ make stop
```


### Tools Used:
In this project, I use some tools listed below. But you can use any simmilar library that have the same purposes. But, well, different library will have different implementation type. Just be creative and use anything that you really need. 

- All libraries listed in [`go.mod`](https://github.com/ilmimris/poc-gofiber-clean-arch/blob/master/go.mod) 
- ["github.com/vektra/mockery".](https://github.com/vektra/mockery) To Generate Mocks for testing needs.