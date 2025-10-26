FROM ubuntu:22.04

# Prevent interactive prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Install SSH server and essential tools
RUN apt-get update && apt-get install -y \
    openssh-server \
    sudo \
    curl \
    wget \
    vim \
    htop \
    nginx \
    postgresql \
    docker.io \
    systemctl \
    ufw \
    fail2ban \
    rsync \
    tar \
    gzip \
    net-tools \
    telnet \
    dnsutils \
    tree \
    git \
    nodejs \
    npm \
    python3 \
    python3-pip \
    && rm -rf /var/lib/apt/lists/*

# Create SSH directory
RUN mkdir /var/run/sshd

# Create users with different privileges
RUN useradd -m -s /bin/bash ubuntu && \
    echo 'ubuntu:ubuntu123' | chpasswd && \
    usermod -aG sudo ubuntu && \
    echo 'ubuntu ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

RUN useradd -m -s /bin/bash devops && \
    echo 'devops:devops123' | chpasswd && \
    usermod -aG sudo devops && \
    echo 'devops ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

RUN useradd -m -s /bin/bash admin && \
    echo 'admin:admin123' | chpasswd && \
    usermod -aG sudo admin && \
    echo 'admin ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

# Configure SSH
RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config && \
    sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config && \
    sed -i 's/#PubkeyAuthentication yes/PubkeyAuthentication yes/' /etc/ssh/sshd_config && \
    sed -i 's/#Port 22/Port 22/' /etc/ssh/sshd_config

# Set root password
RUN echo 'root:rootpass' | chpasswd

# Create .ssh directories for all users
RUN mkdir -p /home/ubuntu/.ssh && \
    chown ubuntu:ubuntu /home/ubuntu/.ssh && \
    chmod 700 /home/ubuntu/.ssh && \
    mkdir -p /home/devops/.ssh && \
    chown devops:devops /home/devops/.ssh && \
    chmod 700 /home/devops/.ssh && \
    mkdir -p /home/admin/.ssh && \
    chown admin:admin /home/admin/.ssh && \
    chmod 700 /home/admin/.ssh && \
    mkdir -p /root/.ssh && \
    chmod 700 /root/.ssh

# Create realistic directory structure for testing
RUN mkdir -p /var/www/html /var/log/myapp /etc/myapp /backup /opt/applications && \
    touch /var/log/myapp/app.log && \
    echo "Welcome to SCL Test Server" > /var/www/html/index.html && \
    echo '{"app": "test", "version": "1.0"}' > /etc/myapp/config.json && \
    chmod 755 /var/www/html && \
    chmod 644 /var/www/html/index.html && \
    chmod 755 /var/log/myapp && \
    chmod 644 /var/log/myapp/app.log && \
    chmod 755 /etc/myapp && \
    chmod 600 /etc/myapp/config.json

# Create sample projects and files
RUN mkdir -p /home/ubuntu/projects /home/devops/scripts /opt/data && \
    echo "# Sample Project" > /home/ubuntu/projects/README.md && \
    echo "#!/bin/bash\necho 'Hello from script'" > /home/devops/scripts/hello.sh && \
    chmod +x /home/devops/scripts/hello.sh && \
    chown -R ubuntu:ubuntu /home/ubuntu/projects && \
    chown -R devops:devops /home/devops/scripts

# Configure services
RUN systemctl enable ssh || true && \
    systemctl enable nginx || true

# Create startup script that starts all services
RUN echo '#!/bin/bash\n\
# Start SSH service\n\
service ssh start\n\
\n\
# Start Nginx\n\
service nginx start || true\n\
\n\
# Start PostgreSQL\n\
service postgresql start || true\n\
\n\
# Create some sample log entries\n\
echo "$(date): SSH server started" >> /var/log/myapp/app.log\n\
echo "$(date): Nginx started" >> /var/log/myapp/app.log\n\
echo "$(date): System ready for SCL testing" >> /var/log/myapp/app.log\n\
\n\
# Keep container running\n\
tail -f /dev/null' > /start-services.sh && \
    chmod +x /start-services.sh

# Expose SSH and web ports
EXPOSE 22 80 443

# Start all services
CMD ["/start-services.sh"]