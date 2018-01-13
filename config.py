import os.path
import json

if os.path.exists('config.json'):
    with open('config.json') as config_json_file:
        config_json = json.load(config_json_file)
        paths = config_json.get('paths', [])
        database_path = config_json.get('database_path', 'sqlite:///db/rolodecks.sqlite')
