echo "Run gotest"
gotest . -v --coverprofile cover.out.tmp

echo "Exclude files from coverage"
cat cover.out.tmp | grep -v -e "main.go" -e "security.go" -e "docs/*" -e "claims/*" -e "models/*" -e "file_manager.go" > cover.out

echo "Calc coverage"
go tool cover -func cover.out