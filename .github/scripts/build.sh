#!/bin/bash
set -e

# Create output directories
mkdir -p dist
mkdir -p temp_dist

echo "Starting build process..."

# Loop through directories in src
# Use null delimiter for safe handling of spaces
find src -mindepth 1 -maxdepth 1 -type d -print0 | while IFS= read -r -d '' dir; do
    dirname=$(basename "$dir")
    # Replace spaces with underscores for the package name
    tool_name="${dirname// /_}"
    
    echo "----------------------------------------"
    echo "Processing: $dirname"
    echo "Target Package Name: $tool_name"
    
    # Create temp dir for this tool
    dest_dir="temp_dist/$tool_name"
    mkdir -p "$dest_dir"
    
    if [[ "$dirname" == "install scripts" ]]; then
        echo "Type: Scripts"
        # Special handling for install scripts
        cp "$dir"/*.sh "$dest_dir/" 2>/dev/null || true
        cp "$dir"/README.md "$dest_dir/" 2>/dev/null || true
    else
        echo "Type: Codebase"
        # Build logic
        pushd "$dir" > /dev/null
        
        # Build
        if [ -f Makefile ]; then
            echo "Building with make..."
            make build
        elif [ -f go.mod ]; then
            echo "Building with go build..."
            # Fallback if no Makefile
            go build -o "$tool_name"
        fi

        # Debug: List files after build to verify binary existence
        echo "Files after build in $dirname:"
        ls -la
        
        # Collect Artifacts
        # Find executable files
        # Using -executable which is supported by GNU find (standard on Linux/Ubuntu)
        find . -maxdepth 1 -type f -executable | while read bin; do
            filename=$(basename "$bin")
            # Filter out scripts and source files if they happen to be executable
            if [[ "$filename" != *.sh && "$filename" != *.py && "$filename" != *.go && "$filename" != *.mod && "$filename" != *.sum ]]; then
                echo "Found binary: $filename"
                cp "$bin" "../../$dest_dir/"
            fi
        done
        
        # Copy READMEs
        cp README.md "../../$dest_dir/" 2>/dev/null || true
        cp readme.md "../../$dest_dir/" 2>/dev/null || true
        
        popd > /dev/null
    fi
    
    # Check if we found anything
    if [ -z "$(ls -A "$dest_dir")" ]; then
        echo "Warning: No artifacts found for $tool_name"
    else
        # Create Zip
        echo "Zipping $tool_name..."
        pushd temp_dist > /dev/null
        zip -r "../dist/$tool_name.zip" "$tool_name"
        popd > /dev/null
    fi
done

# Create combined zip
echo "----------------------------------------"
echo "Creating summary package: all_tools.zip"
pushd dist > /dev/null
zip -r all_tools.zip *.zip
popd > /dev/null

echo "Build complete. Artifacts in dist/:"
ls -lh dist/
