#!/bin/bash

# Set the external environment variables
export REGISTRY="$1"
export REPOSITORY="$2"

# Create target folder & mv docker-compose file
mkdir -p /home/ec2-user/shopping-mall-go &&
mv /home/ec2-user/docker-compose-production.yml /home/ec2-user/shopping-mall-go/docker-compose.yml

# Change directory to shopping-mall-go
cd /home/ec2-user/shopping-mall-go

# Authenticate to AWS ECR
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin 617893088694.dkr.ecr.ap-northeast-1.amazonaws.com

# Remove local Docker image
docker rmi $REGISTRY/$REPOSITORY:latest

# Stop and restart Docker Compose services
docker-compose down
docker-compose up -d

# Remove the script after it has finished executing
sudo rm -f /home/ec2-user/$0