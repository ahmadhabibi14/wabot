# WhatsApp Bot

## Development

- **Dependencies** :
  - Go Programming Language https://go.dev/
  - SQLite Database https://www.sqlite.org/
  - OpenAI API key, generate at https://platform.openai.com/account/api-keys
- **.env Variables**:
  - `OPENAI_API_KEY`
- **How to run**:
  - `go get` to install required modules
  - `start.sh` (Make sure to assign execute permission)
 
## Deployment

#### 1. SSH Configuration
Generate SSH key
```sh
ssh-keygen
```
At local computer, copy publickey to `~/.ssh/authorized_keys` at remote server

#### 2. Github Secrets Variables
- `SERVER_IP` remote server IP
- `SERVER_PASS` remote server Password
- `SERVER_PORT` opened port
- `SERVER_USER` remote server user to action
- `SSH_PRIVATE_KEY` generated ssh private key on local computer
  
#### 3. Remote Server Configuration
Install dependencies
- Go Programming Language
- SQlite

Set SystemD configuration
```sh
# Go to SystemD services directory
cd /lib/systemd/system/

# Create a SystemD service file
touch whatsappbot.service
```
At file `/lib/systemd/system/whatsappbot.service`, write this configuration
```sh
[Unit]
Description=Habi-Bot (WhatsApp Bot)

[Service]
User=root
Type=simple
Restart=always
RestartSec=5s
WorkingDirectory=/root/wabot
ExecStart=/root/wabot/main

[Install]
WantedBy=multi-user.target
```

Enable SystemD service:
```sh
sudo systemctl start whatsappbot.service
```

#### Deploy with existing script

Go to `/deploy` directory, and run script

```shell
# Run this script, note that your ssh private key is authorized by server
./deploy.sh
```