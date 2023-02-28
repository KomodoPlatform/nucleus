#!/usr/bin/env python3

import json
import os
import pathlib
import shutil
import subprocess

NUCLEUS_BIN = "nucleusd"
NUCLEUS_SYSROOT = str(pathlib.Path.home().joinpath(".nucleus"))

CHAIN_ID = "nucleus-1"
PAIR1 = "alice"
PAIR2 = "bob"
MONIKER = "nimda"

DEV_ADDRESS = "nuc1llp0f6qxemgh4g4m5ewk0ew0hxj76avukahzwu"

PAIR1_MNEMONIC = "explain chalk inch snake snack fade news bus horn grant stereo surface panic sister absurd lens speed never inhale element junk senior bubble return"
# PAIR1_ADDRESS="nuc15d4sf4z6y0vk9dnum8yzkvr9c3wq4q89hzvgjk"

PAIR2_MNEMONIC = "young twist cube soldier sorry false dry actor assault roast mosquito brain fix seat wrong sauce impact smoke expect donkey hard decline polar dinosaur"
# PAIR2_ADDRESS="nuc1yfygf0zr5s69ces9r0h72hqv23nkqz9nej732n"


def check_if_already_initialized():
    if os.path.isdir(NUCLEUS_SYSROOT):
        answer = input(NUCLEUS_SYSROOT + " already exists.\n" +
                       "Type 'yes' if you want to remove it and create fresh configurations: ")
        if answer.lower() == "yes":
            shutil.rmtree(NUCLEUS_SYSROOT)
        else:
            print("Aborting!")
            quit()


def init_genesis():
    subprocess.run([NUCLEUS_BIN, "init", MONIKER,
                   "--chain-id="+CHAIN_ID, "--home="+NUCLEUS_SYSROOT])
    subprocess.run([NUCLEUS_BIN, "config", "chain-id",
                   CHAIN_ID, "--home="+NUCLEUS_SYSROOT])
    subprocess.run([NUCLEUS_BIN, "config", "keyring-backend",
                   "test", "--home="+NUCLEUS_SYSROOT])

    process_1 = subprocess.Popen(
        ["echo", PAIR1_MNEMONIC], stdout=subprocess.PIPE)
    process_2 = subprocess.Popen([NUCLEUS_BIN, "keys", "add", PAIR1, "--keyring-backend=test",
                                 "--recover", "--home="+NUCLEUS_SYSROOT], stdin=process_1.stdout)
    process_2.communicate()

    process_1 = subprocess.Popen(
        ["echo", PAIR2_MNEMONIC], stdout=subprocess.PIPE)
    process_2 = subprocess.Popen([NUCLEUS_BIN, "keys", "add", PAIR2, "--keyring-backend=test",
                                 "--recover", "--home="+NUCLEUS_SYSROOT], stdin=process_1.stdout)
    process_2.communicate()

    process_1 = subprocess.run([NUCLEUS_BIN, "keys", "show", "-a", PAIR1, "--keyring-backend=test",
                               "--home="+NUCLEUS_SYSROOT], capture_output=True, text=True, encoding="utf-8")
    subprocess.run([NUCLEUS_BIN, "add-genesis-account", process_1.stdout.strip(),
                   "10000000000unucl", "--home="+NUCLEUS_SYSROOT])

    process_1 = subprocess.run([NUCLEUS_BIN, "keys", "show", "-a", PAIR2, "--keyring-backend=test",
                               "--home="+NUCLEUS_SYSROOT], capture_output=True, text=True, encoding="utf-8")
    subprocess.run([NUCLEUS_BIN, "add-genesis-account", process_1.stdout.strip(),
                   "10000000000unucl", "--home="+NUCLEUS_SYSROOT])

    subprocess.run([NUCLEUS_BIN, "add-genesis-account", DEV_ADDRESS,
                   "10000000000unucl", "--home="+NUCLEUS_SYSROOT])

    subprocess.run([NUCLEUS_BIN,
                    "gentx",
                    PAIR1,
                    "200000000unucl",
                    "--chain-id="+CHAIN_ID,
                    "--moniker="+MONIKER,
                    "--commission-max-change-rate=0.01",
                    "--commission-max-rate=0.20",
                    "--commission-rate=0.05",
                    "--fees=2500unucl",
                    "--from="+PAIR1,
                    "--keyring-backend=test",
                    "--home="+NUCLEUS_SYSROOT])

    subprocess.run([NUCLEUS_BIN, "collect-gentxs", "--home="+NUCLEUS_SYSROOT])

    with open(NUCLEUS_SYSROOT + "/config/genesis.json") as r:
        text = r.read().replace("stake", "unucl")
    with open(NUCLEUS_SYSROOT + "/config/genesis.json", "w") as w:
        w.write(text)


def update_app_toml():
    os.remove(NUCLEUS_SYSROOT + "/config/app.toml")
    shutil.copyfile("app.toml", NUCLEUS_SYSROOT + "/config/app.toml")


def update_rpc_host():
    with open(NUCLEUS_SYSROOT + "/config/config.toml") as r:
        text = r.read().replace("127.0.0.1:26657", "0.0.0.0:26657")
    with open(NUCLEUS_SYSROOT + "/config/config.toml", "w") as w:
        w.write(text)

def add_denom_metadata(data):
    denom_metadata = [{
        "description": "The native staking token of Nucleus.",
        "denom_units": [
            {
                "denom": "unucl",
                "exponent": 0,
                "aliases": [
                    "micronucl"
                ]
            },
            {
                "denom": "mnucl",
                "exponent": 3,
                "aliases": [
                    "millinucl"
                ]
            },
            {
                "denom": "nucl",
                "exponent": 6,
                "aliases": []
            }
        ],
        "base": "unucl",
        "display": "nucl",
        "name": "",
        "symbol": "",
        "uri": "",
        "uri_hash": ""
    }]

    data["app_state"]["bank"]["denom_metadata"] = denom_metadata

def add_htlc_genesis_config(data):
    htlc_genesis_config = {
        "params": {
            "asset_params": []
        },
        "htlcs": [],
        "supplies": [],
        "previous_block_time": "1970-01-01T00:00:01Z"
    }

    data["app_state"]["htlc"] = htlc_genesis_config

if __name__ == "__main__":
    check_if_already_initialized()
    init_genesis()
    update_app_toml()
    update_rpc_host()

    with open(NUCLEUS_SYSROOT + "/config/genesis.json", "r+") as genesis_file:
        data = json.load(genesis_file)

        add_denom_metadata(data)
        add_htlc_genesis_config(data)

        genesis_file.seek(0)
        json.dump(data, genesis_file, indent=2)

    print("\nConfigurations files generated under " + NUCLEUS_SYSROOT)
