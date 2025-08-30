# 🌟 Light Control System

## 📖 Overview
The **Light Control System** is an IoT-based solution for managing and controlling lighting infrastructure across multiple cities and zones.  
It allows administrators to create and manage **cities**, **zones**, and **luminaires**, send **commands** to control lights, and schedule those commands through **Dkron** for automated execution.  
Commands are delivered to luminaires in real-time using the **MQTT** protocol.

---

## ✨ Features
- **Cities & Zones Management**  
  Organize lighting infrastructure by geographic hierarchy (City → Zone → Luminaire).

- **Luminaire Management**  
  Add, update, and manage lighting devices.

- **Command Execution**  
  Send ON/OFF, brightness, or custom control commands to luminaires.

- **Scheduling with Dkron**  
  Schedule recurring or one-time jobs to control lights at specific times.

- **Real-time Control via MQTT**  
  Execute commands instantly over the lightweight MQTT protocol.

---

## 🛠 Architecture Overview
```plaintext
[API Server]
   ├── Manages cities, zones, luminaires, and commands
   ├── Schedules jobs via Dkron API
   └── Publishes control messages to MQTT broker
        ↓
[Dkron Scheduler] -- Triggers jobs --> command Exec API --> [MQTT Broker] -- Sends commands --> [Luminaires]
