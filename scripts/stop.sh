#!/bin/bash

# Script to stop the system
echo "Dừng hệ thống..."

# Dừng các container
docker-compose down

echo "Hệ thống đã được dừng thành công!"
