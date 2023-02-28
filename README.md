# Nucleus

Nucleus Node Implementation & Protocol


## What is Nucleus?
The Nucleus interoperability layer is a cutting-edge Inter-Blockchain Communication (IBC) solution developed by Komodo. The [AtomicDEX](https://github.com/KomodoPlatform/AtomicDEX-API) integration facilitates seamless cross-protocol atomic swaps and other XCAT (cross-chain atomic transaction) operations between Cosmos (tendermint) assets and non-tendermint assets, such as UTXO blockchains like Bitcoin, Litecoin, DGB, Doge, as well as EVM-compatible ledgers including Ethereum, Polygon, BNB Smart Chain, and numerous others. With Nucleus, DeFi-developers as well as endusers can enjoy a remarkable level of interoperability between diverse blockchain protocols/ecosystems, while benefiting from the security and reliability of Komodo's advanced technology stack.

## Getting Started

### Build Dependencies

In order to build nucleus successfuly, you will need `golang(+1.16), gnu make and gcc` to be installed on your system.

### Installation

- From source code:
  ```sh
  make clean install
  ```

### Configuration

Check `init.py` bootstrapping script and update/add any configuration you want. You can also manually utilize everything under `~/.nucleusd` directory.

### Usage
```sh
# First, initialize `~/.nucleus` (make sure nucleusd is installed befoer executing this)
python3 init.py
# or for linux and macos users
./init.py

nucleusd start

```


## Contribute
## References