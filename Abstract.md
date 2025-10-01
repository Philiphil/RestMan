Abstract.

Api should provide an easy way, server and client wise to access subresources

as a list should be paginated to respect system limitation
An object aggregating several lists should provide access to those aggregated lists (subresources).

the trad way is
/api/user/{id}/items

so, how to think this though

Ideally I want to specify items in the user router and not create an item router

the thing is, the router is heavily dependent on generics and I'm not sure there's a serverside easy way to handle this without having to create a new router explicitly

items could be think as a simple router in itself with some kind of user_id=filter


something along the lines of this syntax looks like the most desirable server side syntax

	api := router.NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[User](getDB())),
		route.DefaultApiRoutes(),
        route.Subresource(user.items),
	)

From a client perspective ...
Should I hide subresources from the /user/{id} if route.Subresource(user.items), exists ?
Not sure, i should look at what competitors are doing.

______
Technical
Now, this highlights a new underlying problem, I dont want router.NewApiRouter to be too huge
A natural solution to this problem would be givin a factory

factory := ApiRouterFactory(OrmBuilder, CustomSetOfRoute, CustomSetOfConfiguration)
ApiRouter := ApiRouterFactory.Create[Type](route.DefaultApiRoutes(),
    route.Subresource(user.items),
    SetOfConfiguration
)
Something along these lines, one thing, CustomSetOfRoute might not be a good idea
from personal experience, I want the route.DefaultApiRoutes() for most items except 2 
 ________

	api := router.NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[User](getDB())),
		route.DefaultApiRoutes(),
        route.Subresource(user.items),
	)
    The thing is with this sub resource, 
    How to add configuration to this sub resource ...
        route.Subresource(user.items, ListOfConfiguration...),

    How to give an ORM to this sub resource.
        Either I can use the "root level issue" orm, clone it and change the T 
            but I dont know if that's feasible and that's not a good idea
            entity linked though relational logic should be in the same db 99% of the time but there's this 1%
        Either, I have to pass it some way or another ...

        This should be the definitive syntax of a subroute item

        route.Subresource[Item](
            *orm.NewORM(gormrepository.NewRepository[Item](getDB())),
            route.DefaultApiRoutes(),
            //nested subrsource
            route.Subresource[Item](
                *orm.NewORM(gormrepository.NewRepository[Item](getDB())),
                route.DefaultApiRoutes(),
                SetOfConfiguration,
            )
            SetOfConfiguration,
        )

        I wonder if I can do something like this
        api.AddSubRoute[Item]() instead.
        I should, maybe not the generic syntax, but it will make the nested subresource complex.
        also
        *orm.NewORM(gormrepository.NewRepository[Item](getDB())) should really be
        OrmFactory[Item].Create
        __________

