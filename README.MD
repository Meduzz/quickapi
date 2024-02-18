## Quickapi

The idea is to create restlike apis quickly...ier.

As you can see in the example folder, it's focused around your go struct. The struct defines both the api and the domain (and validation). The lib then applies a thin layer to turn it into an api that persists data (via Gin and Gorm).

### The fork in the road

There's two options of how to consume the lib.

1. `Run` - Let the lib handle the Gin instance, and let you focus on providing the entitiy. Hence quick apis.
2. `For` - Provide the lib with a Gin instance, and only wire up the api for the provided entity. Also quick apis, but also flexibel.

Both method will require you to provide the gorm.DB connection. And `Run` uses `For` under the hood. Just look at the example already.