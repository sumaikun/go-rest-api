# go-rest-api

This api was made with golang and the packages for jwt token generator and jwt token middleware


Main instructions:

1. get http server     
    go get -u github.com/gorilla/mux
2. get jwt token generator
    go get -u github.com/dgrijalva/jwt-go
3. get jwt middleware 
    go get -u github.com/auth0/go-jwt-middleware
4. generate server
    go build in project root folder
5. Add toml env variables
    go get -u github.com/BurntSushi/toml

    Must be a necessary a file with a name config.toml and 2 const defined
    port = "8090"
    jwtKey = "anykeystring"
