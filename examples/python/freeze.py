import time
from pyspark.sql import SparkSession

if __name__ == "__main__":
    spark = SparkSession.builder.getOrCreate()
    print("Start")
    time.sleep(10)
    print("Done")