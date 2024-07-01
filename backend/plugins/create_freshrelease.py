import os
import re
import shutil

def replace_content(content):
    """Replace occurrences of 'JIra' with 'Freshrelease' and 'jira' with 'freshrelease'."""
    content = re.sub(r'Jira', 'Freshrelease', content)
    content = re.sub(r'jira', 'freshrelease', content)
    return content

def process_file(source, destination):
    """Copy and modify the file from source to destination."""
    with open(source, 'r') as f:
        content = f.read()

    # Replace content
    new_content = replace_content(content)

    # Write the new content to the destination file
    with open(destination, 'w') as f:
        f.write(new_content)

def process_directory(source_dir, dest_dir):
    """Recursively process the directory to copy and modify files."""
    if not os.path.exists(dest_dir):
        os.makedirs(dest_dir)

    for item in os.listdir(source_dir):
        source_path = os.path.join(source_dir, item)
        
        # Replace 'jira' with 'freshrelease' in the filename
        new_name = re.sub(r'Jira', 'Freshrelease', item)
        new_name = re.sub(r'jira', 'freshrelease', new_name)
        
        dest_path = os.path.join(dest_dir, new_name)
        
        if os.path.isdir(source_path):
            # Recursively process the subdirectory
            process_directory(source_path, dest_path)
        else:
            # Process the file
            process_file(source_path, dest_path)

def main():
    source_dir = "jira"
    dest_dir = "freshrelease"
    
    if not os.path.exists(source_dir):
        print(f"Source directory '{source_dir}' does not exist.")
        return
    
    if os.path.exists(dest_dir):
        print(f"Destination directory '{dest_dir}' already exists. Removing it first.")
        shutil.rmtree(dest_dir)
    
    process_directory(source_dir, dest_dir)
    print(f"Copied and modified files from '{source_dir}' to '{dest_dir}' successfully.")

if __name__ == "__main__":
    main()
