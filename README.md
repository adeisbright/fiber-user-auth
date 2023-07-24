# fiber-user-auth

Fiber User Auth is an Authentication Service built using Go Fiber Framework.

It allows new users to create account, login , and access their profile.

To make this possible, user data is stored when they create an account, and we check against the database when they try to login.

Authentication and Authorization is enabled via the use of JWT .

Once a user logs, we blacklist their token and they cannot access any protected route till they request for a new token through a
fresh login.

## Table of Contents

1. Project Constraints
1. Requirements
1. Downloading the Project
1. Dependencies
1. Environemnt Settings
1. Starting the App
1. API Endpoints

### Project Constraints

The project adheres to the following constraints :

1. A user must provide their username, email , and password to signup

1. No user can view a protected route except they are logged in

### Project Requirements

1. ### Technologies

This project requires that you have mysql running either locally or remotely.

For this project, I made use of the foloowing packages (Standard , and External)

> fmt : For printing to the standard output

> log : Used for logging and modified to serve the purpose of my custom logging

> fiber : Web Server Framework

> Go Redis : For caching blacklisted token

> gorm : ORM for interacting with our database

> mysql : Mysql driver that is used by gorm

> Go Dotenv : For loading environment variables

> Bcrypt : For Password Encryptino and Password Comparision

> Jwt Go : For JSON Web Authentication Implementation

When you download and run this project using : go run main.go , the dependencies will be managed automatically for you

### Downloading the Project

Clone this project using the git clone command.
The stable branch is set to main and that is the branch you should run.
Other development branch as at the time of pushing the codes were short lived and
therefore deleted.

To get this project, do :

`git clone https://github.com/adeisbright/fiber-user-auth.git`

### Environment Settings

Provide appropriate values where needed as specified in env.sample file

Without values , the project will not run

### Starting the APP

After cloning the app, change into the project directory:

> $ cd <directory-name>

Example:

> $ cd fiber-user-auth

To run the app , do :

> $ go run main.go

### API Endpoints

You can clone an already prepared endpoints for testing from postman using
[Postman](https://documenter.getpostman.com/view/24003787/2s946mbWKX)
