import logging
import time

logging.basicConfig(level=logging.INFO)
for i in range(3):
    print("Hello")
    time.sleep(0.01)  # small delay is better to flush out to terminal
