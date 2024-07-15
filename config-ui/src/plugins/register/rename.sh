#!/bin/bash

# Set the source and destination directories
SOURCE_DIR="jira"
DEST_DIR="freshrelease"

# Create the destination directory
mkdir -p $DEST_DIR

# Copy the contents of the source directory to the destination directory
cp -r $SOURCE_DIR/* $DEST_DIR

# Python script for renaming contents
python3 << 'EOF'
import os

source_dir = 'freshrelease'

def replace_in_file(file_path, search_text, replace_text):
    with open(file_path, 'r') as file:
        filedata = file.read()
    
    newdata = filedata.replace(search_text, replace_text)
    
    with open(file_path, 'w') as file:
        file.write(newdata)

def rename_items(root_dir, search_text, replace_text):
    for dirpath, dirnames, filenames in os.walk(root_dir, topdown=False):
        # Replace content in files
        for filename in filenames:
            file_path = os.path.join(dirpath, filename)
            replace_in_file(file_path, search_text, replace_text)
        
        # Rename files
        for filename in filenames:
            if search_text in filename:
                old_file_path = os.path.join(dirpath, filename)
                new_file_path = os.path.join(dirpath, filename.replace(search_text, replace_text))
                os.rename(old_file_path, new_file_path)
        
        # Rename directories
        for dirname in dirnames:
            if search_text in dirname:
                old_dir_path = os.path.join(dirpath, dirname)
                new_dir_path = os.path.join(dirpath, dirname.replace(search_text, replace_text))
                os.rename(old_dir_path, new_dir_path)

# Perform the replacements and renaming
rename_items(source_dir, 'jira', 'freshrelease')
rename_items(source_dir, 'Jira', 'Freshrelease')

EOF

echo "Done!"

