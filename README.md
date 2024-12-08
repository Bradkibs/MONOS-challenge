# Backend Developer trial task

`We are building a Mobile App that has a business directory. The Mobile app allows vendors to register their businesses and users to view these businesses. Before a business is listed to the directory the vendor has to pay for space to list their products.` 
## Technical Design Document
### Introduction
`Purpose`: Outline the high-level architecture for the mobile app business directory.

`Scope`: Focus on backend design, payment system, and integration with third-party services.

### System Architecture Overview:
Mobile App Backend Services using (Golang Fiber for RESTapi)

Database (Postgres)

Caching (Redis)

Load Balancer (Nginx)

Security (Tailscale Firewall)

Monitoring logs (Elastic Stack)

Containerization (Docker)

Hosting (AWS EC2 & AWS S3)

### Core Modules
`Authentication Service`: Third-party provider (e.g., OAuth)

`Business Service`: CRUD operations for businesses and branches

`Subscription Service`: Manage subscription tiers and pricing

`Payment Service`: Integration with payment gateway

`Notification Service`: Send reminders and alerts

`Admin Dashboard`: Manage subscriptions and monitor payments

### Database Design
`Users`: User/vendor details

`Businesses`: Vendor business information

`Branches`: Business branches

`Subscriptions`: Subscription details

`Payments`: Payment transactions

`Invoices`: Invoice details
### Testing Strategy
`Unit Testing`: Validate individual components

`Integration Testing`: Ensure module interactions
