FROM ubuntu:latest
RUN apt-get update && apt-get install -y curl
RUN sh -c "$(curl -sSfL https://release.solana.com/v1.8.3/install)"
ENV PATH="/root/.local/share/solana/install/active_release/bin:$PATH"
ENTRYPOINT [ "solana-test-validator" ]