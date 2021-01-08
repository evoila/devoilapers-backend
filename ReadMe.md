# Operator Automation Backend (OPA-Backend)
Operator Automation golang Backend

![Go](https://github.com/evoila/devoilapers-backend/workflows/Go/badge.svg)

## Project structure
The project structure will probably change during the project.
The current structure is based on https://github.com/golang-standards/project-layout

- /api: Directory for swagger files (auto generated)
- /build: Directory for docker files
- /cmd: Directory for executables
- /configs: Directory for all config files
- /pkg: Directory for libraries
- /scripts: Directory for necessary scripts (i.e. build scripts)
- /test: Directory for tests
- /.github: Directory for Github actions

## Script overview
- InstallSwaggerGenerator 
    - Software to generate a swagger document out of the webserver definition and specific comments
    - Execute to install swaggo 
    - Details: https://github.com/swaggo/swag
- GenerateSwaggerDoc
    - Used to generate a swagger document out of the webserver definition and specific comments
    - Execute in `<ProjectRoot>` to generate the swagger documentation  
    - Details: https://github.com/swaggo/swag
    

## Getting started
Ensure you have Go **1.15.5** installed.

### GoLand configurations 
1. Config to generate the swagger documentation (Windows guide)
    - Ensure you have executed the InstallSwaggerGenerator-Script
    - Ensure you have the "Batch script support" plugin installed
    - Add a new "Batch" configuration (Windows guide):
    - Script: `<ProjectRoot>/scripts/GenerateSwaggerDoc.bat`
    - Working directory: `<ProjectRoot>`
2. Config to build and run the Webservice 
    - Add a new "Go build" configuration
    - Kind: Directory
    - Directory: set directory to `<ProjectRoot>\cmd\service`
    - Working directory: `<ProjectRoot>`
    - Add the first script as a precondition
4. Download dependencies. GoLand should offer you to sync dependencies. If not you can try to execute `go get ./...` in `<ProjectRoot>`
5. Run your Go build config. The webserver should start. You should be able to navigate to http://127.0.0.1:8080/swagger/index.html

### Navigate to swagger
If the webserver has been started the swagger page is available at http://127.0.0.1:8080/swagger/index.html.

## Conventions
Follow the common Go conventions. 

### Branch naming conventions
Use underscores to replace spaces or special characters

| Abbreviation                        | Description            |
| ------------                        | -----------            |
| feat/<ticket_or_issue_reference>    | Feature                |       
| bug/<ticket_or_issue_reference>     | Bug fix                |
| org/<ticket_or_issue_reference>     | Organizational         |
| junk/<any_title>                    | Experiment-branch      |
| release/<release_info>              | Stable releases        |

