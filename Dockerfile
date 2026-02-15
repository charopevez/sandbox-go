FROM golang:latest

# dev tools
RUN apt-get update && apt-get install -y \
    git \
    curl \
    postgresql-client \
    && rm -rf /var/lib/apt/lists/*

# Go tools for VS Code
RUN go install golang.org/x/tools/gopls@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest && \
    go install golang.org/x/tools/cmd/goimports@latest

# non-root user
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID
RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
    && mkdir -p /home/$USERNAME/go \
    && chown -R $USERNAME:$USERNAME /home/$USERNAME

ENV GOPATH=/home/$USERNAME/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR /workspace
