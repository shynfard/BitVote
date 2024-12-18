import pandas as pd
import re

# Load the data from CSV file
data = pd.read_csv("t2.csv")


def convert_to_seconds(value):
    # Match patterns for 'XmYs' or 'Ys' (where X is minutes and Y is seconds)
    match = re.match(r"(?:(\d+)m)?([\d.]+)s", str(value).strip())
    if match:
        minutes = float(match.group(1)) if match.group(1) else 0
        seconds = float(match.group(2))
        # Convert total time to seconds
        return minutes * 60 + seconds
    # Handle cases with only milliseconds or microseconds
    match = re.match(r"([\d.]+)(ms|µs)", str(value).strip())
    if match:
        number, unit = float(match.group(1)), match.group(2)
        if unit == "ms":
            return number / 1000  # milliseconds to seconds
        elif unit == "µs":
            return number / 1_000_000  # microseconds to seconds
    return None  # if format is unexpected


# Apply conversion to 'Poll Execution Time' and 'Mining Execution Time' columns
data["t"] = data["t"].apply(convert_to_seconds)
data["mt"] = data["mt"].apply(convert_to_seconds)

data["MEAN time"] = data["t"] / data["polls"]
sumTimes = data["t"].sum()
sumPolls = data["polls"].sum()
meanTime = sumTimes / sumPolls
print("Mean time: ", meanTime)
print("Polls: ", sumPolls)

# data = data[["polls", "Poll Execution Time", "Mining Execution Time", "MEAN time"]]
data.to_csv("tt1.csv", index=False)

# print("Data has been standardized to seconds and saved to 'standardized_time_data.csv'")
