import json
import os
import shutil

from datetime import datetime

# Untouched version of configurations
RAW_NUCLEUS_CONFIGS_PATH = "./.assets/.nucleus"
# Updated version of configurations
INITIALIZED_NUCLEUS_CONFIGS_PATH = "./.nucleus"
# genesis.json path for initialization process
GENESIS_JSON_PATH = "./.nucleus/config/genesis.json"


def check_if_already_initialized():
    if os.path.isdir(INITIALIZED_NUCLEUS_CONFIGS_PATH):
        answer = input(os.getcwdb().decode() + "/.nucleus already exists.\n"
                       + "Type 'yes' if you want to remove it and create fresh configurations: ")

        if answer.lower() == "yes":
            shutil.rmtree(INITIALIZED_NUCLEUS_CONFIGS_PATH)
        else:
            print("Aborting!")
            quit()


def cp_nucleus_configs():
    shutil.copytree(RAW_NUCLEUS_CONFIGS_PATH, INITIALIZED_NUCLEUS_CONFIGS_PATH)


def fill_initial_values():
    with open(GENESIS_JSON_PATH, "r+") as genesis_file:
        content = json.load(genesis_file)
        content["genesis_time"] = datetime.utcnow().strftime(
            '%Y-%m-%dT%H:%M:%S.%fZ')
        genesis_file.seek(0)
        json.dump(content, genesis_file, indent=2)


if __name__ == "__main__":
    check_if_already_initialized()
    cp_nucleus_configs()
    fill_initial_values()
    print("\nConfigurations files generated under " + os.getcwdb().decode() + "/.nucleus")
