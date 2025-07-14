import subprocess
import os
import shutil
try:
    from rich.console import Console
    HAS_RICH = True
except ImportError:
    HAS_RICH = False

print_pretty = Console().print if HAS_RICH else print

def get_go_mod_deps():
    try:
        output = subprocess.check_output(['go', 'list', '-m', 'all'], text=True)
        deps = []
        for line in output.strip().split('\n')[1:]:  # Skip the first line (the main module)
            parts = line.split()
            if len(parts) == 2:
                module_path, version = parts
                deps.append((module_path, version))
        return deps
    except subprocess.CalledProcessError as e:
        print_pretty("Error getting module list:", e)
        return []

def get_gomodcache():
    try:
        output = subprocess.check_output(['go', 'env', 'GOMODCACHE'], text=True)
        return output.strip()
    except subprocess.CalledProcessError as e:
        print_pretty("Error getting GOMODCACHE:", e)
        return None

def escape_module_path(module_path):
    # Convert the module path to directory-safe format (Go escapes some characters)
    return module_path.replace('/', os.sep).replace('\\', os.sep)

def remove_cached_modules(deps, gomodcache):
    removed_count = 0
    for module_path, version in deps:
        escaped_path = escape_module_path(module_path)
        full_path = os.path.join(gomodcache, f"{escaped_path}@{version}")
        if os.path.exists(full_path):
            shutil.rmtree(full_path)
            print_pretty(f"Removed {full_path}")
            removed_count += 1
        # Silently skip modules that are not found
    
    if removed_count == 0:
        print_pretty("No cached modules found to remove")

def main():
    deps = get_go_mod_deps()
    if not deps:
        print_pretty("No dependencies found.")
        return

    gomodcache = get_gomodcache()
    if not gomodcache:
        print_pretty("Could not find GOMODCACHE.")
        return

    remove_cached_modules(deps, gomodcache)

if __name__ == "__main__":
    try:
        main()
    except PermissionError:
        print_pretty("🚫🚫🚫 Failed to clean dependecies! Please execute as sudo! 🚫🚫🚫", 
                     style="red")
