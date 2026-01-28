# Sign-server Deployment Guide

## Step-by-Step Instructions

### 1、Copy the Environment Configuration File
Copy the .env.example file in the project to .env:
```bash
cp .env.example .env
```

### 2、Configure Wallet Mnemonic Phrase
Edit the .env file and configure the wallet mnemonic phrase.

**Note：** Keep the mnemonic phrase in a safe place and do not disclose it to anyone.

### 3、Port Configuration

The default port is 8080.
- Skip this step if you use the default port.
- Edit the config.json file for custom configuration if you need to modify the port.

### 4、Build the Docker Image
Build the Docker image locally:
```bash
docker build -t local/sign-server:latest sign-server .
```

### 5、Start service
Using Docker Compose to Start Services:
```bash
docker compose -f docker-compose.yml up -d
```

### Notes
Ensure all operations are performed in a secure environment.

It is recommended to check the firewall settings and ensure the corresponding port is accessible before deployment.
