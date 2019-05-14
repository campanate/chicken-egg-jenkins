## Chicken Egg Jenkins

Do you know when you use jenkins to build everything, and then came the question: how am I going to build jenkins? The solution is simple: Golang.

This projects create a entire jenkins infrastructure with one command.


Just type:

```` 
go build main.go
mv chicken-jenkins /usr/local/bin/chicken-jenkins
chicken-jenkins create
````

## Observations
* Still developing
* Only works on aws
* Feel free to help!!!