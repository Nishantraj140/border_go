run:
	./update &>> output.txt &

printlog:
	tail -f -n100 output.txt

buildQuery:
	GOOS=linux go build -o build/query cmd/query/main.go

zipQuery:
	cd build; query.zip query

buildMain:
	GOOS=linux go build -o build/main cmd/esInsert/main.go

zipMain:
	cd build; zip main.zip main