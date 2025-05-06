import time
import uuid
from pyspark.sql import SparkSession

if __name__ == "__main__":
    spark = SparkSession.builder.getOrCreate()
    print("Start")
    time.sleep(10)
    with open(f"/Users/21370371/projects/goflint/check-{uuid.uuid4()}.txt", mode="w") as f:
        f.write("tut")

    print("Done")