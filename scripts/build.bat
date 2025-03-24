@echo off
echo Building AI-Powered Personal Knowledge Base...

echo Building frontend...
cd web
npm install
npm run build
cd ..

echo Building backend...
go build -o bin/pkb.exe cmd/server/main.go

echo Done! Run bin/pkb.exe to start the server.
