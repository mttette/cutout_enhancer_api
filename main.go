package main

import (
	"cutout_enhancer/handlers"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)



func main(){
    r := chi.NewRouter()

    r.Post("/enhance",handlers.PostEnhanceHandler)

    fmt.Println(`Listening on Port 6009 (u cant customiz it yet )
    U can do A post request on /enhance it must has a json payload like this :-
            {
                url: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSd4y7zVmHqMDDZPFYCAtIvlWWGYofVYEwNg4AyzbXsRg&s"
            }
    yeah .. thats it (0 - 0) `,
    )

    http.ListenAndServe(":6009",r)
}

