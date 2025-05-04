import yaml
import json
import sys
import os
import subprocess

CONFIG_PATH = './config/backends.json'
COMPOSE_FILE = 'docker-compose.yml'

def create_compose(backends):
    services = {
        'load-balancer': {
            'build': '.',
            'ports': ['8080:8080'],
            'depends_on': list(backends.keys()),
            'volumes': ['./config:/root/config']
        }
    }

    for name, info in backends.items():
        services[name] = {
            'build': {
                'context': './backends',
                'dockerfile': 'Dockerfile.backend'
            },
            'environment': [f"PORT={info['port']}"]
        }

    compose = {
        'version': '3.8',
        'services': services
    }

    with open(COMPOSE_FILE, 'w') as f:
        yaml.dump(compose, f, default_flow_style=False)

    with open(CONFIG_PATH, 'w') as f:
        json.dump([{
            "url": f"http://{name}:{info['port']}",
            "weight": info['weight']
        } for name, info in backends.items()], f, indent=2)

def parse_args(args, command):
    backends = {}
    if command in ["init", "add"]:
        for arg in args:
            name, port, weight = arg.split(":")
            backends[name] = {
                'port': int(port),
                'weight': int(weight)
            }
    elif command == "remove":
        # No need to parse args for remove command; we will handle it directly
        return args
    return backends

def load_existing_backends():
    if os.path.exists(CONFIG_PATH):
        with open(CONFIG_PATH, 'r') as f:
            config = json.load(f)
        backends = {}
        for entry in config:
            host = entry['url'].split("//")[1].split(":")
            name = host[0]
            port = int(host[1])
            backends[name] = {
                'port': port,
                'weight': entry['weight']
            }
        return backends
    return {}

def remove_backend(backends, identifier):
    """Removes a backend by name or port."""
    if identifier.isdigit():  # If identifier is a port number
        identifier = int(identifier)
        backends = {name: data for name, data in backends.items() if data['port'] != identifier}
    else:  # If identifier is a name
        if identifier in backends:
            del backends[identifier]
        else:
            print(f"Backend with name {identifier} does not exist.")
            return backends
    return backends

def restart_docker():
    print("🔁 Restarting Docker containers...")
    subprocess.run(["docker-compose", "up", "--build"], check=True)

def main():
    if len(sys.argv) < 3:
        print("Usage:")
        print("  python cli.py init <backends>   - Initialize new backends")
        print("  python cli.py add <new_backends> - Add new backends")
        print("  python cli.py remove <backend>   - Remove a backend by name or port")
        print("Example:")
        print("  python cli.py init backend1:9001:2 backend2:9002:3")
        print("  python cli.py add backend3:9003:1")
        print("  python cli.py remove backend3")
        print("  python cli.py remove 9003")
        return

    command = sys.argv[1]
    if command == "remove":
        if len(sys.argv) < 3:
            print("Usage: python cli.py remove <backend_name_or_port>")
            return
        backend_identifier = sys.argv[2]
        existing_backends = load_existing_backends()
        print(f"🗑 Removing backend: {backend_identifier}")
        updated_backends = remove_backend(existing_backends, backend_identifier)
        create_compose(updated_backends)
        restart_docker()
    else:
        new_backends = parse_args(sys.argv[2:], command)
        existing_backends = load_existing_backends()

        if command == "init":
            print("🛠 Creating new configuration...")
            create_compose(new_backends)
            restart_docker()
        elif command == "add":
            print("Adding new backends to existing config...")
            merged = {**existing_backends, **new_backends}
            create_compose(merged)
            restart_docker()
        else:
            print("Unknown command. Use 'init', 'add', or 'remove'.")

if __name__ == "__main__":
    main()
