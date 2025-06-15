# Relif Platform BFF - Production Deployment Guide

## Issues Fixed

✅ **Missing Environment Configuration** - Added default values and fallback handling  
✅ **AWS Configuration Failures** - Made AWS services optional with graceful degradation  
✅ **Build Compilation Errors** - Fixed case handler dependencies  
✅ **Production Configuration** - Added comprehensive configuration management  

## Quick Fix for Production Login Issue

### The Root Cause
Your production deployment was failing because of missing environment variables, specifically `AWS_REGION`, which prevented the application from starting.

### Solution Applied
1. **Made AWS configuration optional** - App now starts even without full AWS setup
2. **Added default values** - All configuration has sensible defaults
3. **Improved error handling** - Better logging and fallback mechanisms
4. **Environment validation** - Clear error messages for missing config

## Production Environment Setup

### Option 1: Environment Variables (Recommended for most deployments)

Set these environment variables in your production deployment:

```bash
# Required
export ENVIRONMENT=production
export AWS_REGION=us-east-1  # Your AWS region
export MONGO_URI=mongodb://your-production-mongo-host:27017/relif_prod
export TOKEN_SECRET=your-super-secure-jwt-secret-key-here

# Optional (with defaults)
export SERVER_PORT=8080
export ROUTER_CONTEXT=/api/v1
export FRONTEND_DOMAIN=your-production-domain.com
export EMAIL_DOMAIN=your-email-domain.com
export MONGO_DATABASE=relif_prod
export S3_BUCKET_NAME=your-s3-bucket
export COGNITO_CLIENT_ID=your-cognito-client-id
export POOL_ID=your-cognito-pool-id
```

### Option 2: AWS Secrets Manager (Enterprise/High-Security)

1. Set minimal environment variables:
```bash
export ENVIRONMENT=production
export AWS_REGION=us-east-1
export SECRET_NAME=relif-platform-bff-config
```

2. Create AWS Secrets Manager secret:
```json
{
  "FRONTEND_DOMAIN": "your-production-domain.com",
  "EMAIL_DOMAIN": "your-email-domain.com",
  "TOKEN_SECRET": "your-production-secret-key",
  "MONGO_URI": "your-mongo-connection-string",
  "MONGO_DATABASE": "production_db",
  "S3_BUCKET_NAME": "your-s3-bucket",
  "COGNITO_CLIENT_ID": "your-client-id",
  "POOL_ID": "your-pool-id"
}
```

## Deployment Platform Examples

### AWS ECS/Fargate
```json
{
  "environment": [
    {"name": "ENVIRONMENT", "value": "production"},
    {"name": "AWS_REGION", "value": "us-east-1"},
    {"name": "MONGO_URI", "value": "mongodb://your-mongo:27017/relif_prod"},
    {"name": "TOKEN_SECRET", "value": "your-secure-secret"}
  ]
}
```

### Docker Compose
```yaml
services:
  relif-bff:
    image: your-relif-bff:latest
    environment:
      - ENVIRONMENT=production
      - AWS_REGION=us-east-1
      - MONGO_URI=mongodb://mongo:27017/relif_prod
      - TOKEN_SECRET=your-secure-secret
    ports:
      - "8080:8080"
```

### Kubernetes
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: relif-bff
spec:
  template:
    spec:
      containers:
      - name: relif-bff
        env:
        - name: ENVIRONMENT
          value: "production"
        - name: AWS_REGION
          value: "us-east-1"
        - name: MONGO_URI
          valueFrom:
            secretKeyRef:
              name: relif-secrets
              key: mongo-uri
        - name: TOKEN_SECRET
          valueFrom:
            secretKeyRef:
              name: relif-secrets
              key: token-secret
```

## Pre-Deployment Validation

Use the provided configuration checker:

```bash
# Make executable
chmod +x scripts/deploy-config.sh

# Check configuration
./scripts/deploy-config.sh
```

This will validate all required variables and show exactly what's missing.

## Testing After Deployment

1. **Health Check**
```bash
curl http://your-app-url/health
# Should return: "Healthy"
```

2. **API Availability**
```bash
curl http://your-app-url/api/v1/auth/sign-in \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'
# Should return 400 (Bad Request) for invalid credentials, not connection error
```

## Security Considerations

### Production Security Checklist

- [ ] `TOKEN_SECRET` is at least 32 characters long and random
- [ ] MongoDB connection uses authentication
- [ ] HTTPS is enabled (handled by load balancer/reverse proxy)
- [ ] `FRONTEND_DOMAIN` is set to your actual domain (for CORS)
- [ ] Environment variables are stored securely (not in code)
- [ ] AWS IAM permissions are minimal (if using AWS services)

### Secure Token Generation
```bash
# Generate a secure token secret
openssl rand -hex 32
```

## Troubleshooting

### Common Issues

1. **"missing AWS_REGION env variable"**
   - Solution: Set `AWS_REGION=us-east-1` (or your region)

2. **"could not initialize mongo client"**
   - Check `MONGO_URI` is correct and MongoDB is accessible
   - Verify network connectivity to database

3. **"configuration validation failed"**
   - Run `./scripts/deploy-config.sh` to see what's missing
   - Ensure `TOKEN_SECRET` is set for production

4. **CORS errors from frontend**
   - Set `FRONTEND_DOMAIN` to your frontend URL
   - Check that protocol (http/https) matches

### Application Logs
The application now provides detailed logging:
- Configuration loading status
- AWS service availability
- Database connection status
- Service initialization progress

## What Changed in the Code

1. **Enhanced Configuration (`settings/settings.go`)**
   - Added default values for all settings
   - Graceful fallback from Secrets Manager to environment variables
   - Better validation and error messages

2. **Resilient AWS Setup (`settings/aws.go`)**
   - Non-blocking AWS configuration
   - Continues without AWS if not available
   - Clear warnings when services are unavailable

3. **Improved Initialization (`cmd/main.go`)**
   - Better error handling in startup sequence
   - Detailed logging of configuration status
   - Graceful degradation when services are unavailable

4. **Deployment Tooling (`scripts/deploy-config.sh`)**
   - Pre-deployment configuration validation
   - Clear guidance on missing variables
   - Environment-specific recommendations

## Next Steps

1. Set the required environment variables in your production environment
2. Run the configuration checker to validate setup
3. Deploy the updated application
4. Verify authentication endpoints are working
5. Test frontend login functionality

The login issue should now be resolved! 🎉 