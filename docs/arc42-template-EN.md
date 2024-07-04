# 

**About arc42**

arc42, the template for documentation of software and system
architecture.

Template Version 8.2 EN. (based upon AsciiDoc version), January 2023

Created, maintained and © by Dr. Peter Hruschka, Dr. Gernot Starke and
contributors. See <https://arc42.org>.

# Introduction and Goals {#section-introduction-and-goals}

## Requirements Overview {#_requirements_overview}
### Basic Usage
- Users manage tasks and pomodoros via CLI or a web interface.
- The system supports creating, updating, deleting, and retrieving tasks and pomodoros.

### General Functionality
- 2 modes of working:
    - Local:
        - Unix Socket API for local clients
    - Remote:
        -  REST API for web clients. 
- Separate storage layer for persistence

### Reporting and Output Requirements
- Console output for CLI interactions
- JSON and HTML responses for web interactions
- Detailed error logging and user-friendly error messages

## Quality Goals {#_quality_goals}
- Maintainability: This is a recreative project it should be plug-and-play and easy to maintain

# Architecture Constraints {#section-architecture-constraints}
- Golang for backend
- React for frontend, typescript and tailwind css

# System Scope and Context {#section-system-scope-and-context}

## Business Context {#_business_context}

**\<Diagram or Table>**

**\<optionally: Explanation of external domain interfaces>**

## Technical Context {#_technical_context}

![Deployment Diagram](docs/images/deployment-diagram.png)

- **Nodes:** Synology Cluster for storage, Jetson cluster for running 
- **Artifacts:** Docker containers, Kubernetes deployments.

# Solution Strategy {#section-solution-strategy}

- Implement server-side logic using Golang for high performance and concurrency. 
- Wrap implementation into Docker containers for consistency and ease of deployment.

- Frontend: Make the frontend statically 

- Deploy the application on a Kubernetes cluster
- Use the Synology NAS for storage for backups robustness.

Communication:
- Use interfaces to map behaviour and use the ports and adapters pattern to ensure switchability
- Ensure that REST and basic socket communication are in the available options   

Maintainability:
-  Modularize the codebase to ensure components can be developed, tested, and maintained independently.
- Leverage Kubernetes for automated deployments, scaling, and management of containerized applications.

# Building Block View {#section-building-block-view}
![Deployment Diagram](docs/images/building-blocks-diagram.png)

## Whitebox Overall System {#_whitebox_overall_system}
Sure, let's translate that structure into the context of your Pomodoro management system using a similar style.

## Whitebox Overall System {#_whitebox_overall_system}

### Rationale

We used functional decomposition to separate responsibilities:

- **Server** shall encapsulate task management logic and database interactions.
- **Client** shall handle user interactions, including both CLI and web interfaces.
- **Storage** shall manage the persistence of tasks and user data.

### Contained Blackboxes

#### Pomodoro Management System Building Blocks

| Building Block     | Description                                                                                             |
|--------------------|---------------------------------------------------------------------------------------------------------|
| **Server**         | Encapsulates backend logic, task management, database interactions, and authentication.                |
| **Client**         | Handles user interactions through the UI, manages local storage, and communicates with the server.     |
| **Storage**        | Provides persistent storage for tasks and user data using Synology NAS.                                |
| **Communication Protocol** | Manages the communication between client and server using REST API and Unix Socket API.         |

### Interfaces

#### Pomodoro Management System Internal Interfaces

| Interface                | Description                                                                                             |
|--------------------------|---------------------------------------------------------------------------------------------------------|
| **User Interaction**     | Users interact with the system via a web interface or CLI to manage tasks and pomodoros.                |
| **REST API**             | Web clients communicate with the server using REST API for task management and user interactions.       |
| **Unix Socket API**      | Local clients interact with the server using Unix Socket API for efficient local communication.         |
| **Database Access**      | The server accesses the database to store and retrieve tasks and user data.                             |
| **External Services**    | The server may need to access external services for additional functionality or data validation.         |


Motivation

:   *\<text explanation>*

#### Contained Building Blocks

The system should support users in managing their tasks and pomodoros efficiently and reliably. By decomposing the system into clear functional components, we can ensure maintainability, scalability, and clear separation of concerns.

### Contained Building Blocks

1. **Server**
    - **Purpose**: Handles all backend logic, including task management, database interactions, and user authentication.
    - **Interfaces**: REST API, Unix Socket API
    - **Quality Characteristics**: High availability, scalability, security.
    - **Directory/File Location**: `pkg/server`, `pkg/core`

2. **Client**
    - **Purpose**: Manages user interactions through a web interface and CLI, and handles local storage.
    - **Interfaces**: REST API, Unix Socket API
    - **Quality Characteristics**: User-friendly, responsive.
    - **Directory/File Location**: `ui/pomo-client`, `pkg/cli`

3. **Communication Protocol**
    - **Purpose**: Manages the communication between the client and server.
    - **Interfaces**: REST API, Unix Socket API
    - **Quality Characteristics**: Reliable, low latency.
    - **Directory/File Location**: `pkg/client`, `pkg/core`

4. **Storage**
    - **Purpose**: Manages persistent storage for tasks and user data.
    - **Interfaces**: Database APIs
    - **Quality Characteristics**: High performance, data integrity.
    - **Directory/File Location**: `pkg/store`

### Important Interfaces

1. **User Interaction**
    - **Description**: Users interact with the system via web interface or CLI to manage tasks and pomodoros.

2. **REST API**
    - **Description**: Web clients communicate with the server using REST API for task management and user interactions.

3. **Unix Socket API**
    - **Description**: Local clients interact with the server using Unix Socket API for efficient local communication.

4. **Database Access**
    - **Description**: The server accesses the database to store and retrieve tasks and user data.

5. **External Services**
    - **Description**: The server may need to access external services for additional functionality or data validation.

:   *\<Description of contained building block (black boxes)>*

Important Interfaces

:   *\<Description of important interfaces>*

### \<Name black box 1> {#__name_black_box_1}

*\<Purpose/Responsibility>*

*\<Interface(s)>*

*\<(Optional) Quality/Performance Characteristics>*

*\<(Optional) Directory/File Location>*

*\<(Optional) Fulfilled Requirements>*

*\<(optional) Open Issues/Problems/Risks>*

### \<Name black box 2> {#__name_black_box_2}

*\<black box template>*

### \<Name black box n> {#__name_black_box_n}

*\<black box template>*

### \<Name interface 1> {#__name_interface_1}

...

### \<Name interface m> {#__name_interface_m}

## Level 2 {#_level_2}

### White Box *\<building block 1>* {#_white_box_emphasis_building_block_1_emphasis}

*\<white box template>*

### White Box *\<building block 2>* {#_white_box_emphasis_building_block_2_emphasis}

*\<white box template>*

...

### White Box *\<building block m>* {#_white_box_emphasis_building_block_m_emphasis}

*\<white box template>*

## Level 3 {#_level_3}

### White Box \<\_building block x.1\_\> {#_white_box_building_block_x_1}

*\<white box template>*

### White Box \<\_building block x.2\_\> {#_white_box_building_block_x_2}

*\<white box template>*

### White Box \<\_building block y.1\_\> {#_white_box_building_block_y_1}

*\<white box template>*

# Runtime View {#section-runtime-view}

## \<Runtime Scenario 1> {#__runtime_scenario_1}

-   *\<insert runtime diagram or textual description of the scenario>*

-   *\<insert description of the notable aspects of the interactions
    between the building block instances depicted in this diagram.\>*

## \<Runtime Scenario 2> {#__runtime_scenario_2}

## ... {#_}

## \<Runtime Scenario n> {#__runtime_scenario_n}

# Deployment View {#section-deployment-view}

## Infrastructure Level 1 {#_infrastructure_level_1}

***\<Overview Diagram>***

Motivation

:   *\<explanation in text form>*

Quality and/or Performance Features

:   *\<explanation in text form>*

Mapping of Building Blocks to Infrastructure

:   *\<description of the mapping>*

## Infrastructure Level 2 {#_infrastructure_level_2}

### *\<Infrastructure Element 1>* {#__emphasis_infrastructure_element_1_emphasis}

*\<diagram + explanation>*

### *\<Infrastructure Element 2>* {#__emphasis_infrastructure_element_2_emphasis}

*\<diagram + explanation>*

...

### *\<Infrastructure Element n>* {#__emphasis_infrastructure_element_n_emphasis}

*\<diagram + explanation>*

# Cross-cutting Concepts {#section-concepts}

## *\<Concept 1>* {#__emphasis_concept_1_emphasis}

*\<explanation>*

## *\<Concept 2>* {#__emphasis_concept_2_emphasis}

*\<explanation>*

...

## *\<Concept n>* {#__emphasis_concept_n_emphasis}

*\<explanation>*

# Architecture Decisions {#section-design-decisions}

# Quality Requirements {#section-quality-scenarios}

## Quality Tree {#_quality_tree}

## Quality Scenarios {#_quality_scenarios}

# Risks and Technical Debts {#section-technical-risks}
# Glossary

| Term         | Definition                                   |
|--------------|----------------------------------------------|
| Pomodoro     | A time management method involving intervals |
| CLI          | Command Line Interface                       |
| REST API     | Representational State Transfer Application Programming Interface |
| gRPC         | A high-performance RPC framework             |
| Unix socket  | A communication endpoint for local inter-process communication |
| GraphQL      | A query language for APIs                    |

