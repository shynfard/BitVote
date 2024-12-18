import pandas as pd
import matplotlib.pyplot as plt
import re


plt.rcParams.update({"pdf.fonttype": 42})
plt.rcParams.update({"ps.fonttype": 42})

# Load the CSV data
data = pd.read_csv("data.csv")


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
# data["t"] = data["t"].apply(convert_to_seconds)
# data["t"] = data["t"].apply(convert_to_seconds)
# data["mt"] = data["mt"].apply(convert_to_seconds)
# data["MEAN time"] = data["t"] / data["pp"]
# data["mt"] = data["t"] / data["i"]


# # Plotting Poll Execution Time vs Polls
# plt.figure(figsize=(10, 6))
# plt.plot(data["polls"], data["Poll Execution Time"], marker="o")
# plt.title("Poll Execution Time vs Number of Polls")
# plt.xlabel("Number of Polls")
# plt.ylabel("Poll Execution Time")
# plt.grid()
# plt.show()

# # Plotting Poll Size of poll vs Polls
# plt.figure(figsize=(10, 6))
# plt.plot(data["pp"], data["MEAN time"], marker="o")
# plt.title("Average Vote Execution Time for Number of Participants")
# plt.xlabel("Number of Participants")
# plt.ylabel("Vote Execution Time")
# plt.grid()
# plt.show()

# # Plotting Mining Execution Time vs Polls
plt.figure(figsize=(10, 6))
plt.plot(data["p"], data["s"], marker="o")
plt.title("Average single vote counting Time vs Number of votes")
plt.xlabel("Number of votes")
plt.ylabel("Time")
plt.grid()
plt.show()

# # Plotting Mining Memory Usage vs Polls
# plt.figure(figsize=(10, 6))
# plt.plot(data["polls"], data["Mining Memory Usage"], marker="o")
# plt.title("Mining Memory Usage vs Number of Polls")
# plt.xlabel("Number of Polls")
# plt.ylabel("Mining Memory Usage")
# plt.grid()
# plt.show()
