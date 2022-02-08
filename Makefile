NAME=game_of_life
SRC=main.go game_of_life.go

build:
	GOARCH=amd64 GOOS=linux   go build -o $(NAME)     $(SRC)
	GOARCH=amd64 GOOS=windows go build -o $(NAME).exe $(SRC)

clean:
	go clean
	rm -f $(NAME)
	rm -f $(NAME).exe
