# Example service

Boilerplate for architecture of a service based on [go-kit](https://gokit.io) ([repo](https://github.com/go-kit/kit)).

**Service** - is where all of the business logic is implemented. A service usually glues together multiple endpoints. In Go kit, services are typically modeled as interfaces, and implementations of those interfaces contain the business logic. Go kit services should strive to abide [the Clean Architecture](https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html) or [the Hexagonal Architecture](http://alistair.cockburn.us/Hexagonal+architecture). That is, the business logic should have no knowledge of endpoint- or especially transport-domain concepts: your service shouldn’t know anything about HTTP headers, or gRPC error codes.

**Endpoints** is like an action/handler on a controller; it’s where safety and antifragile logic lives. If you implement two transports (HTTP and gRPC), you might have two methods of sending requests to the same endpoint.

**Transport** is bound to concrete transports like HTTP or gRPC. In a world where microservices may support one or more transports, this is very powerful; you can support a legacy HTTP API and a newer RPC service, all in a single microservice.

To get more details go to the official documentation - [go-kit - architecture and design](https://gokit.io/faq/#architecture-and-design).