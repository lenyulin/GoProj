@rm weDY

@set GOOS=linux

@set GOARCH=arm64

@go build -o WeDY ./

@docker rmi -f leiyulin/wedy:v0.0.1

@docker build -t leiyulin/wedy:v0.0.1


C:\Users\lyl69\GolandProjects\GoProj>cd ./wedy

C:\Users\lyl69\GolandProjects\GoProj\wedy>set GOOS=linux

C:\Users\lyl69\GolandProjects\GoProj\wedy>set GOARCH=arm64

C:\Users\lyl69\GolandProjects\GoProj\wedy>set CGO_ENABLED=0

C:\Users\lyl69\GolandProjects\GoProj\wedy>go build -o wedy ./