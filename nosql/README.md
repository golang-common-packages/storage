# NoSQL

````go
package main

import (
    "fmt"
    "github.com/golang-common-packages/database/nosql"
)

func main() {
    mongo := nosql.New(nosql.MONGODB, &nosql.Config{MongoDB: nosql.MongoDB{
                        User:     "USERNAME",
                        Password: "PASSWORD",
                        Hosts:    "URI",
                        Options:  "MONGO_OPTIONS",
                        DB:       "DATABASE_NAME",
                    }})
    
    result, err := mongo.GetALL("DATABASE_NAME", "COLLECTION_NAME", "ID_AFTER", "LIMIT", "YOUR_MODEL_TYPE")
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }
    
    fmt.Print(result)
}
````