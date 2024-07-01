import os
import re
import sys
from datetime import datetime, timedelta

def get_files(directory, pattern):
    """Get all files matching the pattern in the directory."""
    return sorted([f for f in os.listdir(directory) if re.match(pattern, f)])

def get_incremented_dates(num_dates, reference_date):
    """Generate a list of incremented dates starting from the reference date."""
    dates = []
    current_date = reference_date
    
    while len(dates) < num_dates:
        # Increment day
        current_date += timedelta(days=1)
        dates.append(current_date)
    
    return dates

def rename_and_update_files(directory):
    """Rename files, update their contents, and replace old names in all .go files."""
    pattern = r'(\d{8})_(.*)\.go'
    files = get_files(directory, pattern)
    num_files = len(files)
    
    # Start from the first day of the last month
    today = datetime.now()
    last_month = today.replace(day=1) - timedelta(days=1)
    first_day_last_month = last_month.replace(day=1)
    
    new_dates = get_incremented_dates(num_files, first_day_last_month)
    
    # Mapping dictionary
    old_to_new_map = {}
    
    for old_name, new_date in zip(files, new_dates):
        old_date, rest = re.match(pattern, old_name).groups()
        new_date_str = new_date.strftime("%Y%m%d")
        new_name = f"{new_date_str}_{rest}.go"
        
        old_path = os.path.join(directory, old_name)
        new_path = os.path.join(directory, new_name)
        
        # Rename the file
        os.rename(old_path, new_path)
        
        # Update the contents of the file
        with open(new_path, 'r') as file:
            content = file.read()
        
        updated_content = content.replace(old_date, new_date_str)
        
        with open(new_path, 'w') as file:
            file.write(updated_content)

        # Add to mapping
        old_to_new_map[old_date] = new_date_str

        print(f"Renamed {old_name} to {new_name} and updated contents.")
    
    # Replace old dates with new ones in all .go files
    replace_old_with_new_in_files(directory, old_to_new_map)

def replace_old_with_new_in_files(directory, old_to_new_map):
    """Replace old dates with new ones in all .go files in the directory."""
    for root, _, files in os.walk(directory):
        for file_name in files:
            if file_name.endswith(".go"):
                file_path = os.path.join(root, file_name)
                
                with open(file_path, 'r') as file:
                    content = file.read()
                
                updated_content = content
                
                for old_date, new_date in old_to_new_map.items():
                    updated_content = updated_content.replace(old_date, new_date)
                
                with open(file_path, 'w') as file:
                    file.write(updated_content)
                
                print(f"Updated references in {file_name}.")

# Path to the directory containing the files
if len(sys.argv) != 2:
    print("Usage: python script.py <directory>")
    sys.exit(1)

directory = sys.argv[1]

# Call the function
rename_and_update_files(directory)
