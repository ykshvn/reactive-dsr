# DSR Simulation - Dynamic Source Routing Protocol Visualizer

A interactive web application designed to visualize and demonstrate the working principles of the **(DSR)** reactive routing protocol used in Mobile Ad-hoc Networks (MANETs).

Built as a pet project to deeply understand reactive routing protocols, microservices architecture, and real-time visualization.

---

## Project Goal

The main objective of this project is to study and demonstrate how the **DSR** protocol works through:

- Interactive network topology visualization
- Step-by-step simulation of Route Request (RREQ) and Route Reply (RREP) packets
- Clear demonstration of source routing, route caching, and flooding mechanisms

---

## Key Features

- **Interactive Graph Visualization** - Up to 50 nodes with force-directed layout
- **Step-by-step Simulation** - Manual and automatic mode for RREQ/RREP propagation
- **Realistic DSR Implementation**:
  - Route Discovery process
  - Route Caching at each node
  - Duplicate detection using Request ID
  - Loop prevention
- **Real-time WebSocket Updates** - Live visualization of packet propagation
- **Microservices Architecture** with API Gateway
- **Customizable Network Topology** - Random connected graphs without bridges
- **Route Cache Inspection** - View cached routes for each node

---

## Architecture

### Backend (Go)

- **API Gateway** - Entry point, routing, middleware, rate limiting
- **Simulation Service** - Core service responsible for:
  - Graph generation and validation
  - Node simulation (each node runs in its own goroutine)
  - DSR protocol logic (RREQ/RREP processing)
  - Event streaming via WebSocket

### Frontend (Angular)

- Built with **Angular 19+** and **TypeScript**
- Graph visualization using **Cytoscape.js**
- Real-time updates through WebSocket
- Intuitive controls for simulation (Play, Pause, Step, Reset)

---

## Tech Stack

### Backend

- **Go 1.24+**
- **Chi** - HTTP router
- **Gorilla WebSocket**
- **Zap** - Structured logging
- **Viper** - Configuration management
- **Docker & Docker Compose**

### Frontend

- **Angular 19**
- **TypeScript**
- **Cytoscape.js** - Network visualization
- **RxJS** - WebSocket handling
- **Angular Material** - UI components

### Architecture & DevOps

- Microservices architecture
- API Gateway pattern
- Event-driven communication (in-memory message bus)
- Docker Compose
- REST + WebSocket API
