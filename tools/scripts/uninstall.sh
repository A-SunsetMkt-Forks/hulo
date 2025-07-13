#!/bin/bash
# Copyright 2025 The Hulo Authors. All rights reserved.
# Use of this source code is governed by a MIT-style
# license that can be found in the LICENSE file.

set -e

INSTALL_PATH=${1:-"$HOME/.local/bin"}

# Colors for output
INFO_COLOR='\033[1;36m'
SUCCESS_COLOR='\033[1;32m'
ERROR_COLOR='\033[1;31m'
WARNING_COLOR='\033[1;33m'
NC='\033[0m' # No Color

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

write_remove() {
    echo -e "[REMOVE] $1"
}

write_header() {
    echo ""
    echo "=== $1 ==="
    echo ""
}

remove_binary() {
    local install_dir=$1

    write_remove "Removing hulo binary from: $install_dir"

    local binary_path="$install_dir/hulo"
    if [ -f "$binary_path" ]; then
        rm -f "$binary_path"
        write_success "Binary removed successfully: $binary_path"
    else
        write_warning "Binary not found: $binary_path"
    fi

    # Remove from PATH in shell config files
    local shell_rc=""
    if [ -n "$ZSH_VERSION" ]; then
        shell_rc="$HOME/.zshrc"
    elif [ -n "$BASH_VERSION" ]; then
        shell_rc="$HOME/.bashrc"
    fi

    if [ -n "$shell_rc" ] && [ -f "$shell_rc" ]; then
        if grep -q "export PATH.*$install_dir" "$shell_rc" 2>/dev/null; then
            # Remove the PATH line that contains the install directory
            sed -i.bak "/export PATH.*$install_dir/d" "$shell_rc"
            write_info "Removed from PATH: $install_dir"
        else
            write_info "PATH does not contain: $install_dir"
        fi
    fi
}

remove_hulo_modules() {
    write_remove "Removing HULO_MODULES directory"

    # Get HULOPATH from environment or shell config
    local hulo_path="$HULOPATH"
    if [ -z "$hulo_path" ]; then
        # Try to get from shell config
        local shell_rc=""
        if [ -n "$ZSH_VERSION" ]; then
            shell_rc="$HOME/.zshrc"
        elif [ -n "$BASH_VERSION" ]; then
            shell_rc="$HOME/.bashrc"
        fi

        if [ -n "$shell_rc" ] && [ -f "$shell_rc" ]; then
            hulo_path=$(grep "export HULOPATH=" "$shell_rc" | sed 's/export HULOPATH="\(.*\)"/\1/' 2>/dev/null || echo "")
        fi
    fi

    if [ -n "$hulo_path" ] && [ -d "$hulo_path" ]; then
        write_info "Found modules directory: $hulo_path"

        # Confirm deletion
        read -p "Are you sure you want to delete the HULO_MODULES directory? (y/N): " confirm
        if [[ $confirm =~ ^[Yy]$ ]]; then
            rm -rf "$hulo_path"
            write_success "Modules directory removed successfully: $hulo_path"
        else
            write_info "Skipped deletion of modules directory"
        fi
    else
        write_warning "HULO_MODULES directory not found"
    fi
}

remove_environment_variable() {
    write_remove "Removing HULOPATH environment variable"

    # Remove from shell config files
    local shell_rc=""
    if [ -n "$ZSH_VERSION" ]; then
        shell_rc="$HOME/.zshrc"
    elif [ -n "$BASH_VERSION" ]; then
        shell_rc="$HOME/.bashrc"
    fi

    if [ -n "$shell_rc" ] && [ -f "$shell_rc" ]; then
        if grep -q "export HULOPATH=" "$shell_rc" 2>/dev/null; then
            sed -i.bak '/export HULOPATH=/d' "$shell_rc"
            write_info "Removed HULOPATH from user environment variables"
        else
            write_info "HULOPATH not found in user environment variables"
        fi
    fi

    # Remove from global profile if exists
    local global_profile="/etc/profile.d/hulo.sh"
    if [ -f "$global_profile" ]; then
        write_warning "Found global HULOPATH configuration. Administrator privileges may be required."
        read -p "Do you want to remove the global HULOPATH configuration? (y/N): " confirm
        if [[ $confirm =~ ^[Yy]$ ]]; then
            if sudo rm -f "$global_profile"; then
                write_success "Removed HULOPATH from global environment variables"
            else
                write_warning "Failed to remove global HULOPATH (requires administrator privileges)"
            fi
        else
            write_info "Skipped removal of global HULOPATH"
        fi
    else
        write_info "HULOPATH not found in global environment variables"
    fi

    # Remove from current session
    if [ -n "$HULOPATH" ]; then
        unset HULOPATH
        write_info "Removed HULOPATH from current session"
    fi
}

main() {
    write_header "Hulo Uninstaller"

    write_info "Starting uninstallation process"
    write_info "Install path: $INSTALL_PATH"

    # Remove binary
    write_step "Step 1: Removing binary file"
    remove_binary "$INSTALL_PATH"

    # Remove HULO_MODULES directory
    write_step "Step 2: Removing modules directory"
    remove_hulo_modules

    # Remove environment variable
    write_step "Step 3: Removing environment variables"
    remove_environment_variable

    write_header "Uninstallation Summary"
    write_success "Uninstallation completed successfully"
    write_info "Hulo has been completely removed from your system"
}

main "$@"
