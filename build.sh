#!/bin/bash

# Check if fyne is installed
if ! command -v fyne &> /dev/null
then
                echo "fyne could not be found, installing..."
                go install fyne.io/fyne/v2/cmd/fyne@latest
fi

# Check if fyne-cross is installed
if ! command -v fyne-cross &> /dev/null
then
                echo "fyne-cross could not be found, installing..."
                go install github.com/fyne-io/fyne-cross@latest
fi

# Define the list of platforms
platforms=("darwin" "windows" "linux" "all")

# Define the app id
app_id="com.millt.xivrus_cleaner"

# Prompt the user for the platform
echo "Please enter the platform (darwin, windows, linux, all):"
read platform_input

# If "all" is specified, build for all platforms
if [ "$platform_input" = "all" ]; then
        for platform in "${platforms[@]}"
        do
                # Build the application for the current platform
                fyne-cross $platform -app-id $app_id
        done
else
        # Build for the specified platform
        fyne-cross $platform_input -app-id $app_id
fi