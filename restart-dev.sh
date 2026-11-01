#!/bin/bash

# Stop and remove existing containers
docker compose down

# Rebuild the backend image
docker compose build backend

# Start all services
docker compose up -d

# Show logs
docker compose logs -f backend
