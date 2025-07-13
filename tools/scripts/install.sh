#!/bin/bash
# Copyright 2025 The Hulo Authors. All rights reserved.
# Use of this source code is governed by a MIT-style
# license that can be found in the LICENSE file.

set -e

VERSION=${1:-"v0.1.0"}
INSTALL_PATH=${2:-"$HOME/.local/bin"}

BASE_URL="https://github.com/hulo-lang/hulo/releases/download/$VERSION"

write_info() {
    echo -e "[INFO] $1"
}

write_success() {
    echo -e "[SUCCESS] $1"
}

write_error_and_exit() {
    echo -e "[ERROR] $1" >&2
    exit 1
}

write_warning() {
    echo -e "[WARNING] $1"
}

write_step() {
    echo -e "[STEP] $1"
}

write_download() {
    echo -e "[DOWNLOAD] $1"
}

write_install() {
    echo -e "[INSTALL] $1"
}

write_verify() {
    echo -e "[VERIFY] $1"
}

write_header() {
    echo ""
    echo "=== $1 ==="
    echo ""
}

get_operating_system() {
    case "$(uname -s)" in
        Linux*)     echo "Linux" ;;
        Darwin*)    echo "Darwin" ;;
        CYGWIN*|MINGW*|MSYS*) echo "Windows" ;;
        *)          write_error_and_exit "Unsupported operating system" ;;
    esac
}

get_architecture() {
    case "$(uname -m)" in
        x86_64)     echo "x86_64" ;;
        aarch64)    echo "arm64" ;;
        arm64)      echo "arm64" ;;
        i386)       echo "i386" ;;
        *)          write_error_and_exit "Unsupported architecture: $(uname -m)" ;;
    esac
}

get_file_extension() {
    local os=$1
    if [ "$os" = "Windows" ]; then
        echo "zip"
    else
        echo "tar.gz"
    fi
}

download_file() {
    local url=$1
    local output_path=$2

    write_download "Downloading from: $url"

    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$output_path" "$url"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "$output_path" "$url"
    else
        write_error_and_exit "Neither curl nor wget found. Please install one of them."
    fi
    write_success "Download completed successfully"
}

verify_checksum() {
    local file_path=$1
    local expected_checksum=$2

    write_verify "Verifying SHA256 checksum..."

    if command -v sha256sum >/dev/null 2>&1; then
        local actual_checksum=$(sha256sum "$file_path" | cut -d' ' -f1)
    elif command -v shasum >/dev/null 2>&1; then
        local actual_checksum=$(shasum -a 256 "$file_path" | cut -d' ' -f1)
    else
        write_warning "Checksum verification skipped (sha256sum/shasum not available)"
        return 0
    fi

    if [ "$actual_checksum" = "$expected_checksum" ]; then
        write_success "Checksum verification passed"
    else
        write_error_and_exit "Checksum verification failed. Expected: $expected_checksum, Got: $actual_checksum"
    fi
}

extract_file() {
    local file_path=$1
    local extract_dir=$2

    write_step "Extracting archive: $file_path"

    if [[ "$file_path" == *.zip ]]; then
        if command -v unzip >/dev/null 2>&1; then
            unzip -q "$file_path" -d "$extract_dir"
        else
            write_error_and_exit "unzip command not found. Please install unzip."
        fi
    elif [[ "$file_path" == *.tar.gz ]]; then
        if command -v tar >/dev/null 2>&1; then
            tar -xzf "$file_path" -C "$extract_dir"
        else
            write_error_and_exit "tar command not found. Please install tar."
        fi
    else
        write_error_and_exit "Unknown file format: $file_path"
    fi
    write_success "Archive extracted successfully"
}

install_binary() {
    local binary_path=$1
    local install_dir=$2

    write_install "Installing hulo binary to: $install_dir"

    # Create installation directory
    mkdir -p "$install_dir"
    write_info "Created installation directory: $install_dir"

    # Copy binary file
    local install_path="$install_dir/hulo"
    cp "$binary_path" "$install_path"
    chmod +x "$install_path"

    write_success "Binary installed successfully: $install_path"

    # Add to PATH environment variable
    local shell_rc=""
    if [ -n "$ZSH_VERSION" ]; then
        shell_rc="$HOME/.zshrc"
    elif [ -n "$BASH_VERSION" ]; then
        shell_rc="$HOME/.bashrc"
    fi

    if [ -n "$shell_rc" ] && ! grep -q "$install_dir" "$shell_rc" 2>/dev/null; then
        echo "export PATH=\"$install_dir:\$PATH\"" >> "$shell_rc"
        write_info "Added to PATH: $install_dir"
        write_info "Please restart your terminal or run 'source $shell_rc' to apply changes"
    else
        write_info "PATH already contains: $install_dir"
    fi
}

install_stdlib() {
    local extract_dir=$1
    local install_dir=$2

    write_install "Installing standard library and modules"

    # Create HULO_MODULES directory in user's home directory
    local hulo_modules_dir="$HOME/HULO_MODULES"
    mkdir -p "$hulo_modules_dir"
    write_info "Created modules directory: $hulo_modules_dir"

    # Copy all files except executable and archives to HULO_MODULES directory
    write_step "Copying standard library files to: $hulo_modules_dir"
    local file_count=0
    for item in "$extract_dir"/*; do
        local basename=$(basename "$item")
        if [ "$basename" != "hulo" ] && [ "$basename" != "hulo.exe" ] && [[ "$basename" != *.zip ]] && [[ "$basename" != *.tar.gz ]]; then
            if [ -d "$item" ]; then
                cp -r "$item" "$hulo_modules_dir/"
                ((file_count++))
            else
                cp "$item" "$hulo_modules_dir/"
                ((file_count++))
            fi
        fi
    done
    write_success "Standard library installed successfully ($file_count items)"

    # Set HULOPATH environment variable globally
    local global_profile="/etc/profile.d/hulo.sh"
    write_step "Configuring HULOPATH environment variable"

    # Create global profile script
    sudo tee "$global_profile" > /dev/null << EOF
#!/bin/bash
export HULOPATH="$hulo_modules_dir"
EOF

    # Make it executable
    sudo chmod +x "$global_profile"

    # Set for current session and display
    export HULOPATH="$hulo_modules_dir"
    write_info "HULOPATH set to: $hulo_modules_dir"
    write_info "Global profile created: $global_profile"
}

main() {
    write_header "Hulo Installer $VERSION"

    local os=$(get_operating_system)
    local arch=$(get_architecture)
    local extension=$(get_file_extension "$os")

    local filename="hulo_${os}_${arch}.${extension}"
    local download_url="$BASE_URL/$filename"

    write_info "Target platform: $os $arch"
    write_info "Download URL: $download_url"
    local temp_dir=$(mktemp -d)
    local download_path="$temp_dir/$filename"

    # 动态获取 checksums.txt 文件
    local checksums_url="$BASE_URL/checksums.txt"
    write_download "Downloading checksums from: $checksums_url"

    # 下载 checksums.txt 文件
    local checksums_temp=$(mktemp)
    if ! download_file "$checksums_url" "$checksums_temp"; then
        write_error_and_exit "Failed to download checksums.txt"
    fi

    # 解析 checksums.txt 文件
    declare -A checksums
    while IFS= read -r line; do
        line=$(echo "$line" | tr -d '\r' | xargs)
        if [ -n "$line" ] && [[ ! "$line" =~ ^# ]]; then
            checksum=$(echo "$line" | awk '{print $1}')
            filename=$(echo "$line" | awk '{print $2}')
            if [ -n "$checksum" ] && [ -n "$filename" ]; then
                checksums["$filename"]="$checksum"
            fi
        fi
    done < "$checksums_temp"

    rm -f "$checksums_temp"
    write_info "Successfully loaded checksums for ${#checksums[@]} files"

    local expected_checksum=${checksums[$filename]}
    if [ -z "$expected_checksum" ]; then
        write_error_and_exit "No checksum found for $filename in checksums.txt"
    fi

    # Cleanup function
    cleanup() {
        if [ -d "$temp_dir" ]; then
            rm -rf "$temp_dir"
        fi
    }

    # Set trap to cleanup on exit
    trap cleanup EXIT

    download_file "$download_url" "$download_path"
    verify_checksum "$download_path" "$expected_checksum"
    extract_file "$download_path" "$temp_dir"

    # Find binary file
    local binary_path=""
    if [ "$os" = "Windows" ]; then
        binary_path=$(find "$temp_dir" -name "hulo.exe" -type f | head -n 1)
    else
        binary_path=$(find "$temp_dir" -name "hulo" -type f | head -n 1)
    fi

    if [ -z "$binary_path" ] || [ ! -f "$binary_path" ]; then
        write_error_and_exit "Binary file not found in extracted archive"
    fi

    install_binary "$binary_path" "$INSTALL_PATH"

    # Install standard library files
    install_stdlib "$temp_dir" "$INSTALL_PATH"

    write_header "Installation Summary"
    write_success "Installation completed successfully"
    write_info "You can now use 'hulo' command"
    write_info "Standard library and modules are available at: $HULOPATH"
    write_info "HULOPATH environment variable: $HULOPATH"
}

main "$@"
