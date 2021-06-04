# Solana

## Create token

```shell
solana config get
Config File: /Users/dmitrymomot/.config/solana/cli/config.yml
RPC URL: http://127.0.0.1:8899
WebSocket URL: ws://127.0.0.1:8900/ (computed)
Keypair Path: /Users/dmitrymomot/.config/solana/id.json
Commitment: confirmed
➜  ~ solana config set --url https://devnet.solana.com
Config File: /Users/dmitrymomot/.config/solana/cli/config.yml
RPC URL: https://devnet.solana.com
WebSocket URL: wss://devnet.solana.com/ (computed)
Keypair Path: /Users/dmitrymomot/.config/solana/id.json
Commitment: confirmed
➜  ~ solana-keygen new -o /Users/dmitrymomot/.config/solana/id.json
Refusing to overwrite /Users/dmitrymomot/.config/solana/id.json without --force flag
➜  ~ spl-token create-token
Creating token 3nkWEWXbDfgsGis8sEhP8Pb7pHvyNCdnye71qdQrYFhE
Fee payer, 6z5rXfBP1t8aycK5VYX1KFqYkz336TYvEDfadaV1DXrX, has insufficient balance: 0.0014716 required, 0 available
➜  ~ solana airdrop 10
Requesting airdrop of 10 SOL

Signature: 28b5EAPtu9rNzcopcAzU66wFtUUbYLy6FSvucjYDhERWXGT89wt8fpDXLgEupnX1ZtQKqqJ5T2bXu54yY3gJVkU8

10 SOL
➜  ~ spl-token create-token
Creating token 4iC8n6BB6mxozHKYKUSKcZFJrLQNPTfs8ZAFxob2kX7A
Signature: 4HpPiC1boW5AKUdQrxmmQqkuug6en1URi2s89h2CgFiiiZiEojFPa2krGFqo7SqQjPr7As9RaVjk9X8xx5GhipMc
➜  ~ spl-token supply 4iC8n6BB6mxozHKYKUSKcZFJrLQNPTfs8ZAFxob2kX7A
0
➜  ~ spl-token create-account 4iC8n6BB6mxozHKYKUSKcZFJrLQNPTfs8ZAFxob2kX7A
Creating account G5Kxg7iy8d5gzkDb2MtRTvHgbK7zTjE5A9RbUbnrFito
Signature: 5VRgdEXFS31hGfkEyEoxjn4S58LAv8Sveg13tV24KCSzhbqAgMSvu8B61HGPMrLrAQxZe39vGohrfnJu3g1TzuZs
➜  ~ spl-token balance 4iC8n6BB6mxozHKYKUSKcZFJrLQNPTfs8ZAFxob2kX7A
0
➜  ~ spl-token mint 4iC8n6BB6mxozHKYKUSKcZFJrLQNPTfs8ZAFxob2kX7A 1000
Minting 1000 tokens
  Token: 4iC8n6BB6mxozHKYKUSKcZFJrLQNPTfs8ZAFxob2kX7A
  Recipient: G5Kxg7iy8d5gzkDb2MtRTvHgbK7zTjE5A9RbUbnrFito
Signature: 43kfqCY4gtMxrGEknbkgRzeshvgDZNtp5KvxmQUjVs6ar5tRrXMnzgUDQNF5XjdKvw81CXzhMo27kUMrfXqfSYyq
➜  ~ spl-token supply 4iC8n6BB6mxozHKYKUSKcZFJrLQNPTfs8ZAFxob2kX7A
1000
➜  ~ spl-token balance 4iC8n6BB6mxozHKYKUSKcZFJrLQNPTfs8ZAFxob2kX7A
1000
➜  ~ spl-token accounts
Token                                         Balance
---------------------------------------------------------------
4iC8n6BB6mxozHKYKUSKcZFJrLQNPTfs8ZAFxob2kX7A  1000
```