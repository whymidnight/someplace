import json
import sys

def main():
    file = sys.argv[1]
    json_content = {}
    with open(file, "r") as file_contents:
        json_content = json.load(file_contents)

    for item_k in json_content["items"].keys():
        json_content["items"][item_k]["onChain"] = False

    with open(file, "w") as file_contents:
        file_contents.write(json.dumps(json_content, indent="  "))


if __name__ == "__main__":
    main()
