# Foreign Currency Service
This APIs to be used by front-end engineers to develop an application that store and display foreign exchange rate for currencies on a daily basis.

## Quick Start
### 1. Go and Glide Installation
##### A. Go
1. Download Go latest version SDK from https://golang.org/dl/
2. Extract the file using command :
    ````
    tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
    ````
3. Edit the bin in `$HOME/.profile` or `./bashrc` using text editor. then add
    ````
    export PATH=$PATH:/usr/local/go/bin
    ````
4. Create GO workspace folder, inside the folder create the following folder :
    ````
     > src
     > pkg
     > bin
    ````
5. Set GOPATH value to path of our workspace folder.

    ````
    export GOPATH=$PATH_TO_YOUR_FOLDER/workspace
    ````
6. Check the configuration in system by running `go env`.

##### B. Go Dep
Go Dep is the package manager that we used on this project. Please follow the instruction here https://github.com/golang/dep

### 3. Prepare Database
    This service use mysql database. 
    
    1. You must install mysql database.
    
    2. Export this database to you database server, in file ./foreign_currency.sql
    
### 4. Setup configuration file on this service in configs/app.yaml
````
# Database connection and credential configuration
database_sql_slave:
  db_name: 'foreign_currency'
  host: '<your host db>'
  port: 3306
  user: '<your user db>'
  password: '<your password user>'
  charset: 'utf8'

# Database connection and credential configuration
database_sql_master:
  db_name: 'foreign_currency'
  host: '<your host db>'
  port: 3306
  user: '<your user db>'
  password: '<your password user>'
  charset: 'utf8'
````

### 4. Build
To build :

    $ [ -d $GOPATH/src/github/vaksi/foreign_currency ] || mkdir -p $GOPATH/src/github/vaksi/foreign_currency
    $ cd $GOPATH/src/github/vaksi/foreign_currency
    $ git clone http://github.com/vaksi/foreign_currency
    $ cd $GOPATH/src/github.com/vaksi/foreign_currency
    $
    $ dep ensure
    $ go build -race

### 5. Run
To run :

    $ ./foreign_currency serve

## Documentation

### Unit Test and Mock
To run unit test for the project :

    $ go test $(go list ./... | grep -v /vendor/) -cover

For mocking we use https://github.com/vektra/mockery. Generate mock for example:
    
    $ cd internal/exchangeRate
    $ mockery --all

For coverage we can use go convey feature. Please see the following https://github.com/smartystreets/goconvey. To view the coverage

    $ cd /{path_to_folder}/foreign_currency
    $ goconvey

### Run With Docker
To run application with docker: 

    * Create Docker File
    * Build Dokerfile
        $ docker build -t foreign_currency:v1 .
    * Run Docker
        $ docker run -d -p 8081:8081 foreign_currency:v1 serve
        
### Code Style

In golang they already define the convention style https://golang.org/doc/effective_go.html and use https://github.com/alecthomas/gometalinter for linter for tracking efficient my code.

### API Documentation

You can see API Documentation in file `api/swagger.yaml` and run using swagger hub or editor.
