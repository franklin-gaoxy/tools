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
        
        # Build Cross-Platform
        echo "Building Cross-Platform binaries..."
        
        # Define platforms
        platforms=("linux/amd64" "linux/arm64" "darwin/arm64")
        
        for platform in "${platforms[@]}"; do
            platform_split=(${platform//\// })
            GOOS=${platform_split[0]}
            GOARCH=${platform_split[1]}
            
            # Map output name suffix
            suffix="${GOOS}_${GOARCH}"
            # Handle darwin as macos for better naming
            if [ "$GOOS" == "darwin" ]; then
                suffix="macos_${GOARCH}"
            fi
            # Handle amd64 as x86 for better naming
            if [ "$GOARCH" == "amd64" ]; then
                suffix="${GOOS}_x86"
            fi
            
            output_name="${tool_name}_${suffix}"
            
            echo "Building for $GOOS/$GOARCH -> $output_name"
            
            # Build command
            env GOOS=$GOOS GOARCH=$GOARCH go build -o "$output_name" 2>/dev/null || echo "Build failed for $platform (might not be Go project)"
            
            # Move if exists
            if [ -f "$output_name" ]; then
                 cp "$output_name" "../../$dest_dir/"
            fi
        done
        
        # Original logic for backward compatibility (optional, or just rely on cross-builds if they cover everything)
        # But since we want specific versions, we might skip the default "native" build or rename it.
        # Let's keep the cross-build logic as the primary artifact generator for Go projects.
        
        # Collect Artifacts (Non-Go binaries or scripts)
        # ... (rest of logic)
        
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
# Package everything in temp_dist (which contains directories like tool_a/, tool_b/)
pushd temp_dist > /dev/null
zip -r "../dist/all_tools.zip" .
popd > /dev/null

echo "Build complete. Artifacts in dist/:"
ls -lh dist/
