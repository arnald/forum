# Project Architecture

The application follows a **Clean Architecture** design.

---

## Design Principles

The architecture is structured around the following core principles:

1. **Separation of Concerns**  
   Each component is responsible for a single, well-defined aspect of the system. This improves maintainability and readability by reducing coupling between unrelated functionality.

2. **Dependency Inversion**  
   The design enforces a hierarchical dependency flow. Higher-level layers depend on abstractions (interfaces) defined by lower-level (inner) layers rather than concrete implementations. This allows for easier testing, flexibility, and substitution of infrastructure details without affecting core business logic.

---

## Layers Overview

The architecture is divided into **three main layers**, organized from the most fundamental (inner) to the most implementation-specific (outer):

---

### 1. Domain Layer

**Purpose:**  
Contains the core business entities and models that define the fundamental concepts of the system.

**Characteristics:**  
- Does *not* depend on any other layer.  
- Contains only pure business rules and data structures.  
- No external dependencies.

**Examples of Entities:**  
- `UserModel`  
- `PostModel`  
- Other domain-specific structs or classes.

---
k
### 2. Application Layer

**Purpose:**  
Defines all the use cases and application-specific operations. This is where business workflows are implemented by orchestrating the entities from the Domain layer.

**Characteristics:**  
- Depends only on the Domain layer.  
- Contains application logic and use case definitions.  
- Does not include infrastructure concerns such as frameworks or databases.

**Examples of Use Cases:**  
- Creating, updating, or deleting posts (CRUD operations)  
- User registration or authentication logic.

---

### 3. Infrastructure Layer

**Purpose:**  
Connects the application to external resources and delivery mechanisms, such as web servers, databases, file systems, or APIs.

**Characteristics:**  
- Depends on the Application layer to invoke business workflows.  
- Implements interfaces defined in inner layers to fulfill dependencies.  
- Contains adapters, frameworks, and runtime infrastructure.

**Examples of Responsibilities:**  
- HTTP handlers and request routing  
- Dependency injection setup  
- Database repositories  
- Server configuration.

---


