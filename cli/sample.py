import logging

logging.basicConfig(level=logging.INFO)
for i in range(10000):
    # Use flush=True to ensure immediate output
    print(i)
