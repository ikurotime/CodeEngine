#!/bin/bash

set -e

echo "🚀 Setting up CodeEngine LeetCode Clone..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Build sandbox containers
echo "🐳 Building sandbox containers..."
docker build -t sandbox-python ./dockerfiles/python/
docker build -t sandbox-nodejs ./dockerfiles/nodejs/
docker build -t sandbox-java ./dockerfiles/java/

# Install Go dependencies
echo "📦 Installing Go dependencies..."
go mod tidy

# Copy environment config
if [ ! -f config.env ]; then
    cp config.env.example config.env
    echo "📝 Created config.env from template. Please update it with your settings."
fi

# Start PostgreSQL with Docker Compose
echo "🗄️ Starting PostgreSQL database..."
docker-compose up -d postgres

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
sleep 10

# Check if database is ready
until docker-compose exec postgres pg_isready -U codeengine; do
    echo "⏳ Waiting for PostgreSQL..."
    sleep 2
done

echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Update config.env with your preferred settings"
echo "2. Run 'go run ./cmd' to start the server"
echo "3. Visit http://localhost:8080/health to test"
echo ""
echo "Or use Docker Compose:"
echo "docker-compose up --build" 