#!/bin/bash

# Relif Platform BFF - Production Deployment Configuration
# This script helps set up environment variables for production deployment

echo "=== Relif Platform BFF Production Configuration ==="
echo ""

# Check if we're in production mode
if [ -z "$ENVIRONMENT" ]; then
    echo "Warning: ENVIRONMENT not set. Setting to 'production'"
    export ENVIRONMENT=production
fi

echo "Current Environment: $ENVIRONMENT"
echo ""

# Required environment variables for production
REQUIRED_VARS=(
    "ENVIRONMENT"
    "AWS_REGION" 
    "MONGO_URI"
    "TOKEN_SECRET"
)

# Optional variables with defaults
OPTIONAL_VARS=(
    "SERVER_PORT:8080"
    "ROUTER_CONTEXT:/api/v1"
    "FRONTEND_DOMAIN:localhost:3000"
    "EMAIL_DOMAIN:localhost"
    "MONGO_DATABASE:relif_prod"
    "S3_BUCKET_NAME:relif-prod-bucket"
    "COGNITO_CLIENT_ID:prod-client-id"
    "POOL_ID:prod-pool-id"
)

echo "Checking required variables..."
missing_vars=()

for var in "${REQUIRED_VARS[@]}"; do
    if [ -z "${!var}" ]; then
        missing_vars+=("$var")
        echo "❌ $var: Not set"
    else
        echo "✅ $var: Set"
    fi
done

echo ""
echo "Checking optional variables (with defaults)..."

for var_def in "${OPTIONAL_VARS[@]}"; do
    var_name=$(echo "$var_def" | cut -d: -f1)
    default_value=$(echo "$var_def" | cut -d: -f2)
    
    if [ -z "${!var_name}" ]; then
        echo "⚠️  $var_name: Not set (will use default: $default_value)"
        export "$var_name"="$default_value"
    else
        echo "✅ $var_name: Set to ${!var_name}"
    fi
done

echo ""

if [ ${#missing_vars[@]} -gt 0 ]; then
    echo "❌ Missing required variables: ${missing_vars[*]}"
    echo ""
    echo "Please set these environment variables before running the application:"
    echo ""
    for var in "${missing_vars[@]}"; do
        case $var in
            "AWS_REGION")
                echo "export AWS_REGION=us-east-1  # Your AWS region"
                ;;
            "MONGO_URI")
                echo "export MONGO_URI=mongodb://your-mongo-host:27017/relif_prod"
                ;;
            "TOKEN_SECRET")
                echo "export TOKEN_SECRET=your-super-secure-jwt-secret-key-here"
                ;;
            *)
                echo "export $var=your-value-here"
                ;;
        esac
    done
    echo ""
    echo "For AWS Secrets Manager (recommended for production):"
    echo "export SECRET_NAME=your-secret-name-in-aws-secrets-manager"
    echo ""
    exit 1
else
    echo "✅ All required variables are set!"
    echo ""
    echo "Configuration Summary:"
    echo "  Environment: $ENVIRONMENT"
    echo "  AWS Region: $AWS_REGION"
    echo "  Server Port: $SERVER_PORT"
    echo "  Database: $MONGO_DATABASE"
    echo "  Frontend Domain: $FRONTEND_DOMAIN"
    echo ""
    
    if [ -n "$SECRET_NAME" ]; then
        echo "  Using AWS Secrets Manager: $SECRET_NAME"
    else
        echo "  Using environment variables (consider AWS Secrets Manager for production)"
    fi
    
    echo ""
    echo "Ready to start the application!"
fi 