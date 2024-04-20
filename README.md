# gRPC_abstract_calculator 
#### (project for Yandex lyceum)

## How to connect to the repository and start the project
#### (I hope you download git, any IDE and etc.)
1. open the terminal and
write ```git clone https://github.com/Smile8MrBread/gRPC_abstract_calculator.git```
2. write ```cd gRPC_abstract_calculator```
3. write ```go mod download```, ```go mod tidy```
4. write in the first terminal ```go run ./app/cmd/server/auth  --config=./app/config/local.yaml``` to start auth service, <br>
in the second terminal ```go run ./app/cmd/server/orkestrator  --config=./app/config/local.yaml``` to start orkestrator service <br>
in the third terminal ```go run ./app/cmd/server/agent --config=./app/config/local.yaml``` to start agent service <br>
in the fourth terminal ```go run ./app/cmd/client/goClients to start client```
#### Nice! You launched the app! Now you can go to the <http://localhost:8080>, and test my project! (design is bad xD)


##### P .S: if you need to change app port from :8080 to any, went "app/cmd/client/goClients/main.go", scroll to the bottom and change port

###### I hope for a positive experience :)
