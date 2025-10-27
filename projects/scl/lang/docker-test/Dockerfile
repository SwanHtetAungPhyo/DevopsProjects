FROM ubuntu:22.04

# Install SSH server and DevOps tools
RUN apt-get update && apt-get install -y \
    openssh-server \
    sudo \
    curl \
    wget \
    vim \
    htop \
    nginx \
    mysql-server \
    docker.io \
    systemctl \
    ufw \
    rsync \
    tar \
    gzip \
    net-tools \
    telnet \
    dnsutils \
    && rm -rf /var/lib/apt/lists/*

# Create SSH directory
RUN mkdir /var/run/sshd

# Create a test user with sudo privileges
RUN useradd -m -s /bin/bash testuser && \
    echo 'testuser:testpass' | chpasswd && \
    usermod -aG sudo testuser && \
    echo 'testuser ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

# Create devops user
RUN useradd -m -s /bin/bash devops && \
    echo 'devops:devopspass' | chpasswd && \
    usermod -aG sudo devops && \
    echo 'devops ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

# Configure SSH
RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config && \
    sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config && \
    sed -i 's/#PubkeyAuthentication yes/PubkeyAuthentication yes/' /etc/ssh/sshd_config && \
    sed -i 's/#Port 22/Port 22/' /etc/ssh/sshd_config

# Set root password
RUN echo 'root:rootpass' | chpasswd

# Create .ssh directories
RUN mkdir -p /home/testuser/.ssh && \
    chown testuser:testuser /home/testuser/.ssh && \
    chmod 700 /home/testuser/.ssh && \
    mkdir -p /home/devops/.ssh && \
    chown devops:devops /home/devops/.ssh && \
    chmod 700 /home/devops/.ssh && \
    mkdir -p /root/.ssh && \
    chmod 700 /root/.ssh

# Create test directories and files
RUN mkdir -p /var/log/myapp /etc/myapp /backup && \
    touch /var/log/myapp/app.log && \
    echo "Test config file" > /etc/myapp/config.txt && \
    chmod 755 /var/log/myapp && \
    chmod 644 /var/log/myapp/app.log

# Configure nginx (basic setup)
RUN systemctl enable nginx || true

# Create startup script
RUN echo '#!/bin/bash\n\
# Start SSH service\n\
service ssh start\n\
\n\
# Start nginx\n\
service nginx start || true\n\
\n\
# Keep container running\n\
tail -f /dev/null' > /start.sh && \
    chmod +x /start.sh

# Expose SSH port
EXPOSE 22 80 443

# Start services
CMD ["/start.sh"]