# ORM and Repository Integration 

This package relies on an ORM (Object-Relational Mapper) and the Repository interface to manage data interactions.

## ORM  
The ORM is used by ApiRouter to interact with an Repository. It acts as the underlying mechanism to connect and manipulate the database through the repository pattern.

## Repository   
Repository is an interface that defines all the necessary functions for the ApiRouter to perform CRUD (Create, Read, Update, Delete) operations. It acts as a contract that any repository implementation must fulfill, ensuring compatibility with the ApiRouter.


## GormRepository

Restman provides a built-in implementation called GormRepository. It adheres to the principle of separating entities (business logic) from models (database representation).  It is using Gorm as the ORM.


By using Restman with GormRepository, you can quickly set up a robust and consistent data layer for your application, following best practices for DDD and ORM usage.


GormRepository is a forked of https://github.com/Ompluscator/gorm-generics