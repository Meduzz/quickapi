# Quickapi

The idea is to create restlike apis quickly...ier.

As you can see in the example folder, it's focused around your go struct. The struct defines both the api and the domain (and validation). 

Youe struct must implement the `model.Entity` interface, it can also optionally implement the `model.PreloadSupport` to preload child collection and `model.ScopeSupport` to support more fancy filtering.

The lib then turn your structs into an api that persists data (via Gin and Gorm).

### The fork in the road

There's a couple of options of how to consume the lib.

1. `GinStarter` - Returns a cobra command that starts a gin server and boots up all provided entities.
2. `http.For` - Provide the lib with a Gin instance which it uses to wire up the api for the provided entity. Also quick apis, but also flexibel.

Both method will require you to provide the gorm.DB connection. And `GinStarter` uses `For` under the hood. Just look at the example already.

## Entity

The entity abstraction, forces you to provide some metadata and utility for your entity.

*Flavours*

* Theres normal mode, which means your struct represents your entire table.
* And there's JSON mode, where your struct represents the data field and quickapi owns the table (ie your struct is stored as json).

### Where (built in)

The Where api allows you to send simple equality filters, ie: `GET /entity/?where[field]=value`.

### Preload (opt in)

Preload is a Gorm api to load "child" collections. Like in the `example/normal`, to load the `pets`-collection when fetching a `person` preload can be used. The Gorm preload can take a condition, like `alive = ?`. Preload is exposed as public api, like `where` and `sort`.

To avoid sql injection, the conditions used in preload is predefined. Multiple conditions can be grouped together, which can be both a blessing and a curse. At the same time, multiple of these groups of preloads can be called in the same api-call.

In each api call, the group can take one parameter. And if there's no preload data in the call, no preloading will be done.

### Scopes (opt in)

While the where api will take you far. The Scopes api (via `model.ScopeSupport`) gives you full control over what is queried and how.

When you add ScopeSupport to a struct, quickapi will attempt to apply all the named filter on Search queries. Note that you must opt in to use a filter when executing the query to the api.

Ex: 
Named filter = `olderThan`
Scope is defined like this:
```go
func(queryMap map[string]string) model.Hook {
    return func(db *gorm.DB) (*gorm.DB) {
        return db.Where("age > ?", queryMap["age"])
    }
}
```

Trigger usage of filter in query like so:
`GET /entity/?olderThan[age]=18`

Scopes also give you a good toolbox if you opt to use json tables. Becuase they are clonky to query with where api alone.

### Sort (built in)

The sort api is very similar to the Where api. Ie: `GET /entity/?sort[field]=asc|desc`. Keep in mind that it can have a performance impact on big tables and most likely wont work on json data.

## Known issues

 * one-to-many *
 Collections defined as one-to-many are created and updated from the owner collection, but not removed (on certain gorm drivers and setups, ie read up in gorm how to setup your relations correctly). The work around for the moment is to also add an api for the child collection. (which also is the only way to get the child collection migrated (setup in db))

 * json entity *
 The update method here will replace the entire document with what you send. If you only intend to update individual fields, use patch.

 ## TODO
 
* select fields in queries?
> not very important, would also require rpc to have a full api.
* api + interfaces for authorization
> prolly part of a larger infrastructure change, like activating engine features.
* batch apis. So a way to create/update/delete many. And a way to fetch many by their ids.
> prolly also part of a larger infrastructure change.