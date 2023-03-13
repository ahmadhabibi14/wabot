if [ -e main ]
then
   ./main
else
   go build cmd/main.go
fi

./main