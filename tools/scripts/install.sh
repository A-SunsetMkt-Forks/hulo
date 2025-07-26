#!/bin/bash
# Copyright 2025 The Hulo Authors. All rights reserved.
# Use of this source code is governed by a MIT-style
# license that can be found in the LICENSE file.

set -e

VERSION=${1:-"latest"}
INSTALL_PATH=${2:-"$HOME/.local/bin"}

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

get_latest_version() {
    # 静默获取版本，不输出任何信息
    local api_url="https://api.github.com/repos/hulo-lang/hulo/releases/latest"
    local response=""

    # 尝试获取 API 响应
    if command -v curl >/dev/null 2>&1; then
        response=$(curl -s "$api_url" 2>/dev/null)
    elif command -v wget >/dev/null 2>&1; then
        response=$(wget -qO- "$api_url" 2>/dev/null)
    fi

    # 如果 API 调用成功，尝试解析
    if [ -n "$response" ]; then
        local version=""

        # 使用 jq 解析 (推荐)
        if command -v jq >/dev/null 2>&1; then
            version=$(echo "$response" | jq -r '.tag_name // empty' 2>/dev/null)
        else
            # 使用简单的 grep/sed 解析
            version=$(echo "$response" | grep -o '"tag_name"[[:space:]]*:[[:space:]]*"[^"]*"' | head -1 | sed 's/.*"tag_name"[[:space:]]*:[[:space:]]*"//;s/".*//')
        fi

        if [ -n "$version" ] && [ "$version" != "null" ] && [ "$version" != "" ]; then
            printf "%s" "$version"
            return 0
        fi
    fi

    # 如果 API 失败，使用默认版本
    printf "%s" "v0.2.0"
}

get_base_url() {
    local version=$1
    echo "https://github.com/hulo-lang/hulo/releases/download/$version"
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

    local actual_checksum
    if command -v sha256sum >/dev/null 2>&1; then
        actual_checksum=$(sha256sum "$file_path" | cut -d' ' -f1)
    elif command -v shasum >/dev/null 2>&1; then
        actual_checksum=$(shasum -a 256 "$file_path" | cut -d' ' -f1)
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

    write_install "Installing binary: $binary_path to: $install_dir"

    # Create installation directory
    mkdir -p "$install_dir"
    write_info "Created installation directory: $install_dir"

    # Copy binary file, keeping original filename
    local filename
    local install_path
    filename=$(basename "$binary_path")
    install_path="$install_dir/$filename"
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

    # Copy all files except bin directory and archives to HULO_MODULES directory
    write_step "Copying standard library files to: $hulo_modules_dir"
    local file_count=0
    for item in "$extract_dir"/*; do
        local basename
        basename=$(basename "$item")
        if [ "$basename" != "bin" ] && [[ "$basename" != *.zip ]] && [[ "$basename" != *.tar.gz ]]; then
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

    # Set HULO_PATH environment variable globally
    local global_profile="/etc/profile.d/hulo.sh"
    write_step "Configuring HULO_PATH environment variable"

    # Create global profile script
    sudo tee "$global_profile" > /dev/null << EOF
#!/bin/bash
export HULO_PATH="$hulo_modules_dir"
EOF

    # Make it executable
    sudo chmod +x "$global_profile"

    # Set for current session and display
    export HULO_PATH="$hulo_modules_dir"
    write_info "HULO_PATH set to: $hulo_modules_dir"
    write_info "Global profile created: $global_profile"
}

main() {
    # 获取最新版本
    if [ "$VERSION" = "latest" ]; then
        write_info "Fetching latest version from GitHub..."
        VERSION=$(get_latest_version)
        # 确保版本号正确获取
        if [ -z "$VERSION" ]; then
            VERSION="v0.2.0"
        fi
        write_success "Latest version found: $VERSION"
    fi

    write_header "Hulo Installer $VERSION"

    local os
    local arch
    local extension
    local base_url
    local filename
    local download_url
    local temp_dir
    local download_path
    local checksums_url
    local checksums_temp
    local expected_checksum

    os=$(get_operating_system)
    arch=$(get_architecture)
    extension=$(get_file_extension "$os")

    base_url=$(get_base_url "$VERSION")
    filename="hulo_${os}_${arch}.${extension}"
    download_url="$base_url/$filename"

    write_info "Target platform: $os $arch"
    write_info "Download URL: $download_url"
    temp_dir=$(mktemp -d)
    download_path="$temp_dir/$filename"

    # 动态获取 checksums.txt 文件
    checksums_url="$base_url/checksums.txt"
    write_download "Downloading checksums from: $checksums_url"

    # 下载 checksums.txt 文件
    checksums_temp=$(mktemp)
    if ! download_file "$checksums_url" "$checksums_temp"; then
        write_error_and_exit "Failed to download checksums.txt"
    fi

    # 解析 checksums.txt 文件
    declare -A checksums
    while IFS= read -r line; do
        line=$(echo "$line" | tr -d '\r' | xargs)
        if [ -n "$line" ] && [[ ! "$line" =~ ^# ]]; then
            local checksum_value
            local checksum_filename
            checksum_value=$(echo "$line" | awk '{print $1}')
            checksum_filename=$(echo "$line" | awk '{print $2}')
            if [ -n "$checksum_value" ] && [ -n "$checksum_filename" ]; then
                checksums["$checksum_filename"]="$checksum_value"
            fi
        fi
    done < "$checksums_temp"

    rm -f "$checksums_temp"
    write_info "Successfully loaded checksums for ${#checksums[@]} files"

    expected_checksum=${checksums[$filename]}
    write_info "Looking for checksum for: $filename"
    write_info "Available files in checksums: $(printf '%s ' "${!checksums[@]}")"
    if [ -z "$expected_checksum" ]; then
        write_error_and_exit "No checksum found for $filename in checksums.txt"
    fi
    write_info "Found checksum: $expected_checksum"

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
    # Find all executables in bin directory
    local bin_dir="$temp_dir/bin"
    if [ ! -d "$bin_dir" ]; then
        write_error_and_exit "bin directory not found in extracted archive"
    fi

    # 使用 mapfile 安全地创建数组
    local exe_files=()
    while IFS= read -r -d '' file; do
        exe_files+=("$file")
    done < <(find "$bin_dir" -maxdepth 1 -type f -perm -u=x -print0)

    if [ ${#exe_files[@]} -eq 0 ]; then
        write_error_and_exit "No executable files found in bin directory"
    fi

    for exe in "${exe_files[@]}"; do
        install_binary "$exe" "$INSTALL_PATH"
    done

    # Install standard library files
    install_stdlib "$temp_dir" "$INSTALL_PATH"

    write_header "Installation Summary"
    write_success "Installation completed successfully"
    write_info "You can now use 'hulo' command"
    write_info "Standard library and modules are available at: $HULO_PATH"
    write_info "HULO_PATH environment variable: $HULO_PATH"
}

main "$@"
