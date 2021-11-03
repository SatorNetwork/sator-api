solc --optimize --overwrite --abi ./erc20.sol -o build
solc --optimize --overwrite --bin ./erc20.sol -o build
abigen --abi=./build/erc20.abi --bin=./build/erc20.bin --pkg=erc20 --out=./erc20.go