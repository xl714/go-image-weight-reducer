# go image weight reducer

Execute the following command to run the code:

```
go mod init github.com/xl714/go-image-weight-reducer
 
go mod tidy

go run main.go

rm _img_1_resized.jpg _img_2_resized.png; go run main.go


```

to import a new external module

```
go get github.com/nfnt/resize
```

Formatting: For better readability, use the go fmt command to automatically format your code:

```
go fmt main.go
```

Building a Binary: To create a standalone executable file, use the go build command:

```
# linux
go build -o build/GoImageWeightReducerBin main.go

# Windows
go build -o build/GoImageWeightReducer.exe main.go
```

Test:

```
# linux
rm _img_1_resized.jpg; rm _img_2_resized.png ; ./build/GoImageWeightReducerBin  --image-max-weight=1.1

# Windows
rm _img_1_resized.jpg; rm _img_2_resized.png ; ./build/GoImageWeightReducer.exe  --image-max-weight=1.1
```
 

This will create an executable file named hello in your project directory. You can run it directly without the go run command.

Using Modules: For larger projects, consider using Go modules for dependency management.
