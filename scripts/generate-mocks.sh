#!/bin/bash

# Function to generate mocks for a given directory
generate_mocks() {
    local BASE_DIR=$1
    local MOCK_PACKAGE=$2

    # Find all files that end with _repo.go or _service.go inside the specified directory
    find "$BASE_DIR" -name "*_repo.go" -o -name "*_service.go" | while read -r FILE
    do
        # Extract the filename without extension
        FILE_NAME=$(basename "$FILE" .go)

        # Set the mock destination path
        MOCK_DESTINATION="$BASE_DIR/mocks/${FILE_NAME}_mock.go"

        # Run mockgen command
        mockgen -source="$FILE" -destination="$MOCK_DESTINATION" -package="$MOCK_PACKAGE"

        # Output a success message
        echo "Generated mock for $FILE at $MOCK_DESTINATION"
    done
}

# Generate mocks for both repository and service directories
generate_mocks "internal/repository" "repository_mock"
generate_mocks "internal/service" "service_mock"
