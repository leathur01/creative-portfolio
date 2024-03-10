### Prerequisites  

Install Go: https://go.dev/doc/install  
Setup Go environment: [https://medium.com/@fsufitch/go-environment-setup-minus-the-insanity-b872f34351c8](https://learn.gopherguides.com/courses/preparing-your-environment-for-go-development/modules/setting-up-mac-linux/#slide-1)  
Install Revel Framework (bỏ bước Install Go và setup GOPATH trong hướng dẫn của Revel): https://revel.github.io/tutorial/gettingstarted.html  
Install Postgresql: https://www.postgresql.org/download/  
Install pgAdmin: https://www.pgadmin.org/download/  
Install migration tool: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate  

### Installation  

1. Clone the repo in the src directory of the GOPATH (GOPATH đã được set up ở bước setup Go environment), ví dụ cho GOPATH: /home/leathur/go  
   git clone https://github.com/leathur01/creative-portfolio.git  

2. Download all the dependencies of the project  
   cd lio  
   go mod download  

4. Change the database-server connection string to connect to your local database server in ./lio/app/init.go  
   connStr := "postgres://db-username:your-db-server-password@localhost/db-name"  

5. Run the migration scripts  
   migrate -path ./lio/migrations -database postgres://postgres:password@local
host/creative_portfolio up  

### Start web server  

   At the root folder of the repository, run this command  
   revel run -a lio/  
   Access the application on localhost:8080  
