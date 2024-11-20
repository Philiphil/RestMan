# RESTMAN

Restman takes Golang Structs and create REST routes.  
Inspired by Symfony and Api Platform.  
Built on top of Gin.

Restman can be used with any database as long as you implement the builtin repository interface
It come with its own GORM based implementation, compatible with Entity/Model separation but also a more straighforward approach.  

## Features  
Fully working structure to REST automated route generation using GIN, recursion and generics
Out of the box GORM based ORM  
Firewall implementation allowing to filter who can access/edit which data  
Serialization groups to control which property are allowed to be readed or wrote using the generated route  


## TODO, Ideas for myself and for random contributors
Pagination Filters
GraphQL-like PageInfo Object / after, before, first, last, pageof
groupS override parameter
entity.ID UUID compatiblility
InMemory cache with default redis integration
Mongo default Repository
Fix XML serialialization
fix CSV serialialization 
Check current golang json serialization speed
check force lowercase for json ? (golang default serializer is like the only thing in the world who does nt force lowercase)
check messagepack
Serializer could be refactord
Somehow hooks could be nice ??  (meh)
The fireWall could have a builtin requireOwnership
subressource pagination
performance evaluation could be nice (is there even a place for paralelilsm somewhere ??)