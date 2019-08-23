# Ossn Backend

## Intro

This project is build using the following:

* [gqlgen](https://gqlgen.com/) for the graphql server.
* [gorm](https://gorm.io/) for the ORM.
* [qor admin](https://github.com/qor/admin) for the admin panel.
* [Dep](https://golang.github.io/dep/) for the package manager
* [PostrgeSQL](https://www.postgresql.org/) for the database.
* [Redis](https://redis.io/) for session caching.

## Setting up a machine

1. Install [go](https://golang.org/doc/install)
2. Install [PostrgeSQL](https://www.postgresql.org/download/)
3. Install [redis](https://redis.io/download)
4. Install [Dep](https://golang.github.io/dep/docs/installation.html)
5. Clone project into `$GOPATH/src/github.com/ossn/ossn-backend` and go into that directory
6. Install dependencies by running `dep ensure`
7. Set the environment variables defined in the .env.example file to your machine
8. Build and run the application `go run .`

You can also use the docker file in this project to build the app but you still need to pass the environment variables to the container and connect the container to a PostgreSQL and to a Redis.

## Deployment

The master branch of this project is automatically deployed to heroku and is used as the production server.

The instructions on "Setting up a machine" can be used to deploy this app to an alternative environment.

## Development notes

* Whenever a model is update this server is responsible for rebuilding the frontend website, in order for the server to achieve this the environment variable `REBUILD_URL` should be set to the rebuild hook of the frontend
* Heroku has an issue with qor's static assets (the static assets of the admin panel), so whenever qor is update the static assets from inside the vendor folder should be placed inside the folder "app"
* Once you've updated the schema.graphql file don't forget to run `go generate ./...` in order for your changes to be reflected in the go code

## Project structure

* app            ->    Static assets of the admin panel (see development notes for more info)
* controllers    ->    The controller for the OpenID connect endpoints
* helpers        ->    Small functions that are used across the project
* main           ->    Main server configuration including routes and middleware declaration
* middlewares    ->    Middleware definition
* models         ->    Database models and database connection
* resolvers      ->    Graphql models, Graphql mutation and query resolvers
* scripts        ->    Allows to run generate command
* Procfile       ->    Heroku's run command
* gqlgen.yml     ->    Gqlgen configuration
* schema.graphql ->    Graphql definition
