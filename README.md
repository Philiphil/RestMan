# RESTMAN

Restman takes Golang Structs and create REST routes.  
Inspired by Symfony and Api Platform.  
Built on top of Gin.  

Restman can be used with any ORM as long as it is provided an implementation of its Repository Interface.  
It come with its own GORM based Implementation, compatible with Entity/Model separation but also a more straighforward approach.  

## Features  
Fully working structure to REST route generation using GIN  
Out of the box GORM based ORM  
[WIP] Firewall implementation allowing to filter who can access/edit which data  
Symfony Serializer allowing serialization groups to control which property are allowed to be readed or wrote using the generated route  


TODO  
batch operation are down for now (except Get-List)  
entity.ID would be so nice if it was UUID compatible somehow  
JSONLD collection's metadata are ALMOST relyable (see tests)  
Creating an Cache Interface somehow would be really nice  
More configuration option, for pagination, by default enable, forced or disabled, max Item per page ...   
Serializer is not as performant as it could be  
Somehow hooks could be nice ??  
The fireWall is in WIP state  