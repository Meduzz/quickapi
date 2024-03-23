## Quickapi

The idea is to create restlike apis quickly...ier.

As you can see in the example folder, it's focused around your go struct. The struct defines both the api and the domain (and validation). 

You then wrap your struct in an `Entity` by calling the `NewEntity`-function.

The lib then applies a thin layer to turn the `Entity` into an api that persists data (via Gin and Gorm).

### The fork in the road

There's two options of how to consume the lib.

1. `Run` - Let the lib handle the Gin instance, and let you focus on providing the entitiy. Hence quick apis.
2. `For` - Provide the lib with a Gin instance, and only wire up the api for the provided entity. Also quick apis, but also flexibel.

Both method will require you to provide the gorm.DB connection. And `Run` uses `For` under the hood. Just look at the example already.

## Entity

The entity abstraction, allows us to tie it to a path in the api and also tie any filters to that entity.

Entities are created through `NewEntity[<struct>](name:string, [...filters:*NamedFilter])` function.

- `name` will be used for the api-path.
- `filters` can be triggered when the api is called via query-params.

See next section on how to create filters.

## Filters

While filters can look awefully involved, they are quite simple but yet powerful. The biggest drawback of them is that you have to opt in to them when you call the api.

In the example-folder I use filters to preload the one-to-many `Pet`-collection when loading anything from the `Person`-collection.

To create a filter use the `NewFilter(name:string, handler:quickapi.Scope)`-function. It takes a name (which is the  name of the query-param where we will look for data). The second param is a `Scope`-function (`func(map[string]string) func(*gorm.DB) *gorm.DB`). This is a higher order function that accepts a `map[string]string` and returns a function `func(*gorm.DB) *gorm.DB`. The returned function is used as a scope in Gorm.

Lets define a filter:
```
    // define the handler function
    func preloadPets() quickapi.Scope {
        return func(m map[string]string) func(*gorm.DB) *gorm.DB {
            return func(query *gorm.DB) *gorm.DB {
                // ... use the data in the map m to change the query
                return query
            }
        }
    }

    // create the filter
    filter := quickapi.NewFilter("pets", preloadPets())

    // use filter when creating entities
    person := quickapi.NewEntity[Person]("", filter)
    ...
```

And now activate the filter in an api-call:
```
    GET /person/?pets[param]=value
    ...
```

Note that it is not enough with `?pets` (neither is `?pets[]` but `?pets[_]` will work).