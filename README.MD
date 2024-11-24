# Quickapi

The idea is to create restlike apis quickly...ier.

As you can see in the example folder, it's focused around your go struct. The struct defines both the api and the domain (and validation). 

You then wrap your struct in an `model.Entity` by calling the `model.NewEntity`-function.

The lib then applies a thin layer to turn the `model.Entity` into an api that persists data (via Gin|Rpc and Gorm).

### The fork in the road

There's a couple of options of how to consume the lib.

1. `GinStarter` - Returns a cobra command that starts a gin server and boots up all provided entities.
2. `http.For` - Provide the lib with a Gin instance which it uses to wire up the api for the provided entity. Also quick apis, but also flexibel.
3. `RpcStarter` -  Returns a cobra command that starts a rpc server and boots up all provided entities.
4. `rpc.For` - Provide the lib with a rpc server which it uses to wire up the api for the provided entity.

Both method will require you to provide the gorm.DB connection. And `GinStarter` uses `For` under the hood. Just look at the example already.

## Entity

The entity abstraction, allows us to tie it to a path in the api and also tie any filters to that entity.

Entities are created through `model.NewEntity[<struct>](name:string, preload:model.PreloadDelegate, [...filters:*model.NamedFilter])` function.

- `name` will be used for the api-path.
- `preload` is a function that returns the preload config to use (if the correct preload paramters was provided)
- `filters` can be triggered when the api is called via query-params.

See next section on how to create filters.

### Preload

Preload is a Gorm api to load "child" collections. Like in the `example/`, to load the `pets`-collection when fetching a `person` preload can be used. The Gorm preload can take a condition, like `alive = ?`. Preload is exposed as public api, like `where` and `sort`.

To avoid sql injection, the conditions used in preload is predefined. Multiple conditions can be grouped together, which can be both a blessing and a curse. At the same time, multiple of these groups of preloads can be called in the same api-call.

In each api call, the group can take one parameter. And if there's no preload data in the call, no preloading will be done.

### Filters

While filters can look awefully involved, they are quite simple but yet powerful. The biggest drawback of them is that you have to opt in to them when you call the api.

In the example-folder I use filters to preload the one-to-many `Pet`-collection when loading anything from the `Person`-collection.

To create a filter use the `model.NewFilter(name:string, handler:model..Scope)`-function. It takes a name (which is the  name of the query-param where we will look for data). The second param is a `model.Scope`-function (`func(map[string]string) model.Hook`). This is a higher order function that accepts a `map[string]string` and returns a function `model.Hook` which in turn is defined as `func(*gorm.DB) *gorm.DB`. The returned function is used as a scope inside the query to Gorm.

Lets define a filter:
```
    // define the handler function
    func preloadPets() model.Scope {
        return func(m map[string]string) model.Hook {
            return func(query *gorm.DB) *gorm.DB {
                // ... use the data in the map m to change the query
                return query
            }
        }
    }

    // create the filter
    filter := model.NewFilter("pets", preloadPets())

    // use filter when creating entities
    person := model..NewEntity[Person]("person", filter)
    ...
```

And now activate the filter in an api-call:
```
    GET /person/?pets[param]=value
    ...
```

Note that it is not enough to include `?pets` (neither is `?pets[]` but `?pets[_]` will work).

## Known issues

 * one-to-many collections are created and updated from the owner collection, but not removed (on certain gorm drivers and setups, ie read up in gorm how to setup your relations correctly). The work around for the moment is to also add an api for the child collection. (which also is the only way to get the child collection migrated (setup in db))

 ## TODO
 
* json entities
* meta endpoint for each entity
* select fields in queries?
> not very important, would also require rpc to have a full api.
