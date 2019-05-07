.PHONY: dev
dev: lecli.exe
	cp lecli.exe /mnt/c/Users/Robbie/Desktop/

lecli.exe: *.go
	GOOS=windows go build -o lecli.exe
