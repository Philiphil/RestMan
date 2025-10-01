# Security   

Restman provides a mechanism to verify if a user has the appropriate read or write permissions for a specific object. Below is an overview of how this system works:

## Overview  

To enforce permissions, an object should implement one or both of the following interfaces:

    WritingRights
    ReadingRights

These interfaces can be fulfilled by implementing the functions GetReadingRights() and/or GetWritingRights(). These functions should return a user-defined implementation of an AuthorizationFunction, which takes two parameters:

    An User: representing the user attempting to access the object.
    An Entity: representing the object being accessed.

The function should return a boolean value indicating whether the User is permitted to perform the specified operation on the Entity.
Usage with ApiRouter

An ApiRouter accepts a list of firewalls via the AddFireWalls method. Firewalls should implement a GetUser method, which retrieves an User or an error using the Gin request object. The ApiRouter uses these firewalls to fetch the User, then applies the appropriate WritingRights and/or ReadingRights checks to determine whether the User has the required access permissions.

Restman ensures that all requests are validated for appropriate user permissions, maintaining secure and controlled access to your application’s resources.

Note: GetList cannot rely on ReadingRights as the security checks are post fetch.
