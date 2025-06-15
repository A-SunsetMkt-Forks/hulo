import os
import subprocess
import sys
import platform

def get_executable_path():
    system = platform.system().lower()
    arch = platform.machine().lower()

    bin_dir = os.path.join(os.path.dirname(__file__), "..", "bin")

    # Normalize architecture naming
    if arch in ("x86_64", "amd64"):
        arch = "x86_64"
    elif arch in ("aarch64", "arm64"):
        arch = "arm64"
    else:
        raise RuntimeError(f"Unsupported architecture: {arch}")

    # File name varies by OS
    ext = ".exe" if system == "windows" else ""
    filename = f"hulo{ext}"

    exe_path = os.path.join(bin_dir, filename)

    if not os.path.exists(exe_path):
        raise FileNotFoundError(f"Executable not found: {exe_path}")

    return exe_path

def main():
    exe_path = get_executable_path()
    result = subprocess.run([exe_path] + sys.argv[1:])
    sys.exit(result.returncode)
