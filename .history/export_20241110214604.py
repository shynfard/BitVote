import pandas as pd
import re

# Load the data from CSV file
data = pd.read_csv("outPut2.csv")


# Function to convert time to seconds
def convert_to_seconds(value):
    match = re.match(r"(\d+\.?\d*)\s*(ms|s|m|µs)", str(value).strip())
    if match:
        number, unit = float(match.group(1)), match.group(2)
        if unit == "ms":
            return number / 1000  # milliseconds to seconds
        elif unit == "s":
            return number  # already in seconds
        elif unit == "m":
            return number * 60  # minutes to seconds
        elif unit == "µs":
            return number / 1_000_000  # microseconds to seconds
    return None  # if format is unexpected


# Apply conversion to 'Poll Execution Time' and 'Mining Execution Time' columns
data["Poll Execution Time"] = data["Poll Execution Time"].apply(convert_to_seconds)
data["Mining Execution Time"] = data["Mining Execution Time"].apply(convert_to_seconds)

# Save the updated DataFrame to a new CSV file
data.to_csv("standardized_time_data.csv", index=False)

print("Data has been standardized and saved to 'standardized_time_data.csv'")
