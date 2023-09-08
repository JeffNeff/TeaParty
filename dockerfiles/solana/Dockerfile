FROM rust

RUN apt-get update
RUN apt-get -y install pkg-config build-essential libudev-dev nodejs npm
RUN npm i -g yarn
RUN rustup component add rustfmt
RUN rustup default nightly && rustup update
RUN useradd -ms /bin/bash user
USER 1000
RUN sh -c "$(curl -sSfL https://release.solana.com/v1.9.2/install)"
ENV PATH /home/user/.local/share/solana/install/active_release/bin:$PATH
RUN /home/user/.local/share/solana/install/active_release/bin/sdk/bpf/scripts/install.sh
RUN cargo install --git https://github.com/project-serum/anchor --tag v0.19.0 anchor-cli --locked
WORKDIR /usr/src/